package circuitbreaker

import (
	"testing"
	"time"
)

func TestCircuitBreaker_BasicFlow(t *testing.T) {
	cb := NewCircuitBreaker(3, 0.5, 100*time.Millisecond)

	// 初始状态应为 Closed
	t.Logf("Initial state: %v", cb.State())
	if cb.State() != Closed {
		t.Errorf("expected state Closed, got %v", cb.State())
	}

	// 连续失败3次，应该进入 Open
	for i := 0; i < 3; i++ {
		if !cb.AllowRequest() {
			t.Errorf("expected allow request in Closed, got false")
		}
		t.Logf("Request allowed in Closed state")
		cb.RecordFail()
		t.Logf("Recorded failure %d, state: %v", i+1, cb.State())
	}
	if cb.State() != Open {
		t.Errorf("expected state Open after failures, got %v", cb.State())
	}
	t.Logf("State switched to Open after 3 failures")

	// Open 状态下不允许请求
	t.Logf("State after failures: %v", cb.State())
	if cb.AllowRequest() {
		t.Errorf("expected not allow request in Open, got true")
	}
	t.Logf("Request denied in Open state")

	// 等待 retryTimePeriod 后，应该进入 HalfOpen 并允许请求
	time.Sleep(110 * time.Millisecond)
	t.Logf("State after sleep: %v", cb.State())
	if !cb.AllowRequest() {
		t.Errorf("expected allow request in HalfOpen, got false")
	}
	t.Logf("Request allowed in HalfOpen state")
	if cb.State() != HalfOpen {
		t.Errorf("expected state HalfOpen, got %v", cb.State())
	}
	t.Logf("State switched to HalfOpen after sleep")

	// 在 HalfOpen 状态下，成功率达到要求，应该回到 Closed
	cb.RecordSuccess() // 1 success, 1 request
	t.Logf("Recorded success in HalfOpen state, state: %v", cb.State())
	if cb.State() != Closed {
		t.Errorf("expected state Closed after enough success, got %v", cb.State())
	}
	t.Logf("State switched back to Closed after success")
}

func TestCircuitBreaker_HalfOpenFail(t *testing.T) {
	cb := NewCircuitBreaker(2, 0.5, 50*time.Millisecond)

	// 进入 Open
	cb.RecordFail()
	cb.RecordFail()
	t.Logf("State after 2 failures: %v", cb.State())
	if cb.State() != Open {
		t.Errorf("expected state Open, got %v", cb.State())
	}
	t.Logf("State switched to Open after 2 failures")

	time.Sleep(60 * time.Millisecond)
	t.Logf("State after sleep: %v", cb.State())
	if !cb.AllowRequest() {
		t.Errorf("expected allow request in HalfOpen, got false")
	}
	t.Logf("Request allowed in HalfOpen state")
	if cb.State() != HalfOpen {
		t.Errorf("expected state HalfOpen, got %v", cb.State())
	}
	t.Logf("State switched to HalfOpen after sleep")

	// 在 HalfOpen 状态下失败，不会立即回到 Open，但会累计失败
	cb.RecordFail()
	t.Logf("Recorded failure in HalfOpen state, state: %v", cb.State())
	// 这里可以根据你的业务逻辑扩展，比如失败后直接回到 Open等
	// 当前实现不会自动回到 Open，除非你在 HalfOpen 也加上失败阈值判断
}
