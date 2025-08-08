package tokenbucket_test

import (
	"gatesvr/limite/tokenbucket"
	"testing"
	"time"
)

func TestTokenBucketRateLimtImpl_GetToken(t *testing.T) {
	limiter := tokenbucket.NewTokenBucketRateLimtImpl(3, 1000) // 3 tokens, refill rate of 1 token per 1000ms

	// Test initial token consumption (should succeed for first 3 attempts)
	for i := 0; i < 3; i++ {
		if !limiter.GetToken() {
			t.Errorf("Expected token to be available (attempt %d), but got false", i)
		} else {
			t.Logf("Successfully got token (attempt %d)", i)
		}
	}

	// Test token exhaustion (should fail on 4th attempt)
	if limiter.GetToken() {
		t.Errorf("Expected no token to be available (attempt 4), but got true")
	} else {
		t.Logf("Failed to get token (attempt 4), as expected")
	}

	// Test token refill after rate interval (should succeed after waiting)
	time.Sleep(1100 * time.Millisecond) // Wait slightly longer than the rate interval
	if !limiter.GetToken() {
		t.Errorf("Expected token to be available after refill, but got false")
	} else {
		t.Logf("Successfully got token after refill")
	}
}

func BenchmarkTokenBucketRateLimtImpl_GetToken(b *testing.B) {
	limiter := tokenbucket.NewTokenBucketRateLimtImpl(1000, 10) // High capacity and fast rate for benchmarking

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.GetToken()
	}
}
