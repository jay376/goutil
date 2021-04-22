package elect

import (
	"database/sql"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"git.garena.com/shopee/sz-devops/observability/monitoring-platform/pkg/logutil"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var g *gorm.DB

func initDB() {
	EndPoint := "root:root@tcp(localhost:3456)/shopee_infra_monitor_datasource_sg_db"
	if db, err := sql.Open("mysql", EndPoint); err != nil {
		panic(err)
	} else {
		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(100)
		g, _ = gorm.Open("mysql", db)
		if err = g.DB().Ping(); err != nil {
			panic(err)
		}
	}
	go func() {
		logconf := &logutil.LoggingConf{
			Debug:       true,
			OutputPaths: []string{"./info.log"},
			ErrPaths:    []string{"./error.log"},
		}
		l, _ := logutil.SetupZapLogging(logconf, "elector-test")
		defer l.Sync()
		time.Sleep(10 * time.Second)
	}()
}

func BenchmarkKeepalive(b *testing.B) {
	initDB()
	num := 10
	electors := make([]*ElectorDB, num)
	ttl := 3 * time.Second
	for idx := 0; idx < num; idx++ {
		electors[idx] = NewElectorDB(g, Name("bench-test"), ttl, "")
		electors[idx].id = 2
		b.Log(electors[idx].node)
	}

	var once sync.Once
	hasElected := int32(0)
	duration := 1 * time.Minute
	var rwMu sync.RWMutex
	runFunc := func(idx int32) {
		begin := time.Now()
		i := int(idx) % num
		for time.Since(begin) < (duration - ttl) {
			rwMu.RLock()
			electors[i].keepAlive()
			rwMu.RUnlock()
			once.Do(func() {
				atomic.StoreInt32(&hasElected, 1)
			})
			time.Sleep(10 * time.Millisecond)
		}
	}

	go func() {
		begin := time.Now()
		for time.Since(begin) < duration {
			masterNum := 0
			rwMu.Lock()
			for _, elector := range electors {
				if elector.IsMaster() {
					masterNum++
				}
			}
			rwMu.Unlock()
			if atomic.LoadInt32(&hasElected) == 1 {
				assert.Equal(b, 1, masterNum)
			} else {
				assert.Equal(b, 0, masterNum)
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()

	var idx int32
	b.SetParallelism(5)
	b.RunParallel(func(pb *testing.PB) {
		n := atomic.AddInt32(&idx, 1)
		runFunc(n - 1)
		for pb.Next() {
		}
	})
}
