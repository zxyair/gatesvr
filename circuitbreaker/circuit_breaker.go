package circuitbreaker

import (
	"sync/atomic"
	"time"
)

type State int32

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state            State
	failCount        int32
	successCount     int32
	requestCount     int32
	failureThreshold int32
	halfOpenRate     float64
	retryTimePeriod  time.Duration
	lastFailTime     time.Time
}

func NewCircuitBreaker(failureThreshold int, halfOpenRate float64, retryTimePeriod time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            Closed,
		failureThreshold: int32(failureThreshold),
		halfOpenRate:     halfOpenRate,
		retryTimePeriod:  retryTimePeriod,
	}
}

// 是否允许请求
func (cb *CircuitBreaker) AllowRequest() bool {
	now := time.Now()
	switch State(atomic.LoadInt32((*int32)(&cb.state))) {
	case Open:
		if now.Sub(cb.lastFailTime) >= cb.retryTimePeriod {
			atomic.StoreInt32((*int32)(&cb.state), int32(HalfOpen))
			return true
		}
		return false
	case HalfOpen:
		atomic.AddInt32(&cb.requestCount, 1)
		return true
	case Closed:
		return true
	default:
		return true
	}
}

// 记录成功
func (cb *CircuitBreaker) RecordSuccess() {
	if State(atomic.LoadInt32((*int32)(&cb.state))) == HalfOpen {
		atomic.AddInt32(&cb.successCount, 1)
		requestCount := atomic.LoadInt32(&cb.requestCount)
		successCount := atomic.LoadInt32(&cb.successCount)
		if successCount > int32(float64(requestCount)*cb.halfOpenRate) {
			atomic.StoreInt32((*int32)(&cb.state), int32(Closed))
			cb.resetCount()
		}
	} else {
		cb.resetCount()
	}
}

// 记录失败
func (cb *CircuitBreaker) RecordFail() {
	atomic.AddInt32(&cb.failCount, 1)
	cb.lastFailTime = time.Now()
	if State(atomic.LoadInt32((*int32)(&cb.state))) == HalfOpen {
		atomic.StoreInt32((*int32)(&cb.state), int32(Open))
		cb.lastFailTime = time.Now()
	} else if atomic.LoadInt32(&cb.failCount) >= cb.failureThreshold {
		atomic.StoreInt32((*int32)(&cb.state), int32(Open))
	}
}

func (cb *CircuitBreaker) resetCount() {
	atomic.StoreInt32(&cb.requestCount, 0)
	atomic.StoreInt32(&cb.successCount, 0)
	atomic.StoreInt32(&cb.failCount, 0)
}

func (cb *CircuitBreaker) State() State {
	return State(atomic.LoadInt32((*int32)(&cb.state)))
}
