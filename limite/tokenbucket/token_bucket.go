package tokenbucket

import (
	"sync"
	"time"
)

type RateLimit interface {
	GetToken() bool
}

//type TokenBucketRateLimtImpl struct {
//	capacity    int
//	rate        int
//	curCapacity int64
//	timestamp   int64
//}

func NewTokenBucketRateLimtImpl(capacity, rate int64) *TokenBucketRateLimtImpl {
	return &TokenBucketRateLimtImpl{
		capacity:    capacity,
		rate:        rate,
		curCapacity: int64(capacity),
		timestamp:   time.Now().UnixMilli(),
	}
}

//	func (t *TokenBucketRateLimtImpl) GetToken() bool {
//		if atomic.LoadInt64(&t.curCapacity) > 0 {
//			if atomic.AddInt64(&t.curCapacity, -1) >= 0 {
//				return true
//			}
//			atomic.AddInt64(&t.curCapacity, 1) // Rollback if decrement went below 0
//		}
//
//		current := time.Now().UnixMilli()
//		interval := current - atomic.LoadInt64(&t.timestamp)
//
//		if interval >= int64(t.rate) {
//			if interval >= int64(t.rate)*2 {
//				added := int64(interval/int64(t.rate) - 1)
//				newCapacity := atomic.AddInt64(&t.curCapacity, added)
//				if newCapacity > int64(t.capacity) {
//					atomic.StoreInt64(&t.curCapacity, int64(t.capacity))
//				}
//			}
//			atomic.StoreInt64(&t.timestamp, current)
//			log.Debugf("token bucket rate limit, current: %d, rate: %d, interval: %d, capacity: %d", atomic.LoadInt64(&t.curCapacity), t.rate, interval, t.capacity)
//			return true
//		}
//
//		return false
//	}
type TokenBucketRateLimtImpl struct {
	curCapacity int64
	capacity    int64
	rate        int64 // 毫秒
	timestamp   int64 // 毫秒
	mu          sync.Mutex
}

func (t *TokenBucketRateLimtImpl) GetToken() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().UnixMilli()
	elapsed := now - t.timestamp

	// 补充令牌
	if elapsed > 0 {
		tokensToAdd := elapsed / t.rate
		if tokensToAdd > 0 {
			t.curCapacity += tokensToAdd
			if t.curCapacity > t.capacity {
				t.curCapacity = t.capacity
			}
			t.timestamp = now
		}
	}

	if t.curCapacity > 0 {
		t.curCapacity--
		return true
	}
	return false
}
