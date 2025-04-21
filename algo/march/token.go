package march

import "time"

type TokenLimiter struct {
	cap          int
	tokens       int
	lastFillTime time.Time
}

func NewTokenLimiter(qps int) *TokenLimiter {
	return &TokenLimiter{
		cap:          qps,
		tokens:       0,
		lastFillTime: time.Now(),
	}
}

func (l *TokenLimiter) Take(tokens int) int {
	l.fill()
	ret := min(l.tokens, tokens)
	l.tokens = max(l.tokens-tokens, 0)
	return ret
}

func (l *TokenLimiter) fill() {
	now := time.Now()
	add := int(now.Sub(l.lastFillTime) * time.Duration(l.cap) / time.Second)
	if add > 0 {
		l.tokens = min(add+l.tokens, l.cap)
		l.lastFillTime = now
	}
}
