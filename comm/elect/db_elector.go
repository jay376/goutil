package elect

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/xjbdjay/goutil/comm"
	"go.uber.org/zap"
)

// Item ...
type Item struct {
	ID          int64 `gorm:"primary_key"`
	Name        string
	Master      string
	UpdateTime  int64
	ExpiredTime int64
}

// TableName ...
func (Item) TableName() string {
	return "alert_elect_tab"
}

// ElectorDB ...
type ElectorDB struct {
	gorm    *gorm.DB
	id      int64
	name    string
	node    string
	master  string
	ttl     time.Duration
	expired time.Time
	sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewElectorDB name is unique for each elect group. Such as event-server, monitor-api
// addr is ServeAddr, if empty will generate random addr
func NewElectorDB(gorm *gorm.DB, name Name, ttl time.Duration, addr string) *ElectorDB {
	strs := strings.Split(addr, ":")
	var node string
	if addr == "" {
		node = comm.GetIntranetIP() + ":" + comm.RandStringBytes(8)
	} else {
		if len(strs) == 2 {
			if strs[0] == "" {
				node = comm.GetIntranetIP() + ":" + strs[1]
			}
		} else {
			node = addr
		}
	}

	e := &ElectorDB{
		gorm: gorm,
		name: string(name),
		node: node,
		ttl:  ttl,
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())

	return e
}

// Start ...
func (e *ElectorDB) Start() {
	e.compagin()
	go e.loop()
}

// Stop ...
func (e *ElectorDB) Stop() {
	e.cancel()
}

func (e *ElectorDB) init() {
	var item Item
	query := e.gorm.Where("name = ?", e.name).Find(&item)
	if query.Error != nil && query.Error != gorm.ErrRecordNotFound {
		zap.L().Error("query error", zap.Error(query.Error))
		return
	}

	if item.ID == 0 {
		item.Name = e.name
		item.UpdateTime = time.Now().Unix()
		item.ExpiredTime = time.Now().Unix()
		if e.gorm.Create(&item).Error != nil {
			zap.L().Error("create error", zap.Error(query.Error))
			return
		}
	}
	query = e.gorm.Where("name = ?", e.name).First(&item)
	if query.Error == nil {
		e.id = item.ID
	}
}

func (e *ElectorDB) compagin() {
	if e.id == 0 {
		e.init()
	}

	if err := e.keepAlive(); err != nil {
		// if db error occur, try again
		time.Sleep(time.Second)
		e.keepAlive() //nolint
	}
}

func (e *ElectorDB) loop() {
	e.init()
	ticker := time.NewTicker(e.ttl / 3)
	zap.L().Info("elector begin",
		zap.String("node", e.node), zap.Int("ttl seconds", int(e.ttl/time.Second)))
	for {
		select {
		case <-e.ctx.Done():
			return
		case <-ticker.C:
			e.compagin()
		}
	}
}

// IsMaster ...
func (e *ElectorDB) IsMaster() bool {
	e.Lock()
	defer e.Unlock()
	return e.node == e.master && time.Now().Before(e.expired)
}

// Master ...
func (e *ElectorDB) Master() string {
	e.Lock()
	defer e.Unlock()
	return e.master
}

// Node ...
func (e *ElectorDB) Node() string {
	e.Lock()
	defer e.Unlock()
	return e.node
}

func (e *ElectorDB) keepAlive() (err error) {
	logger := zap.L().With(zap.String("func", "keepAlive"))
	var item Item
	tx := e.gorm.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	// Lock row
	// 注：此时的update_time是数据库的当时时间
	if err = tx.Raw(`select id, name, master, UNIX_TIMESTAMP(NOW()) as update_time,
                     expired_time  from alert_elect_tab where id = ? for update`,
		e.id).Scan(&item).Error; err != nil {
		logger.Error("query master error", zap.Error(err))
		if err == gorm.ErrRecordNotFound {
			e.id = 0
		}
		return
	}
	logger.Debug("lock", zap.Any("item", item))

	// master renew
	now := item.UpdateTime
	ttl := int(e.ttl / time.Second)
	if now < item.ExpiredTime {
		if item.Master != e.node {
			logger.Debug("I am not master", zap.String("master", item.Master), zap.String("node", e.node))
			e.Lock()
			e.master = item.Master
			e.Unlock()
		} else {
			up := tx.Exec(`update alert_elect_tab set master = ? , update_time=UNIX_TIMESTAMP(NOW()),
                          expired_time = UNIX_TIMESTAMP(NOW()) + ? where id = ?`, e.node, ttl, e.id)
			if up.Error == nil && up.RowsAffected == 1 {
				e.Lock()
				e.master = item.Master
				e.expired = time.Now().Add(e.ttl)
				logger.Debug("renew lease ", zap.String("node", e.node), zap.Time("lease expiredAt", e.expired))
				e.Unlock()
			}
		}
		return
	}

	// to be master
	up := tx.Exec(`update alert_elect_tab set master = ? , update_time=UNIX_TIMESTAMP(NOW()), expired_time =
         UNIX_TIMESTAMP(NOW()) + ? where id = ? and expired_time <= UNIX_TIMESTAMP(NOW())`, e.node, ttl, e.id)
	if up.Error != nil || up.RowsAffected != 1 {
		err = up.Error
		logger.Error("compagin failed", zap.Error(up.Error), zap.Int64("rowsAffected", up.RowsAffected))
	} else {
		e.Lock()
		e.master = e.node
		e.expired = time.Now().Add(e.ttl)
		logger.Debug("compagin suc", zap.String("node", e.node), zap.Time("lease expiredAt", e.expired))
		e.Unlock()
	}

	return
}
