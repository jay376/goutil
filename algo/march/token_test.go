package march

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	qps := 3000
	l := NewTokenLimiter(qps)

	time.Sleep(time.Second)
	assert.Equal(t, l.Take(qps*2), qps)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, l.Take(qps), qps/10)
}
