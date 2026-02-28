package ws

import "testing"

func TestKDSConnectionLimiterAcquireRelease(t *testing.T) {
	t.Parallel()

	limiter := newKDSConnectionLimiter(2)

	if !limiter.Acquire(1, 10) {
		t.Fatal("expected first acquire to succeed")
	}
	if !limiter.Acquire(1, 10) {
		t.Fatal("expected second acquire to succeed")
	}
	if limiter.Acquire(1, 10) {
		t.Fatal("expected third acquire to be blocked by max limit")
	}

	limiter.Release(1, 10)
	if !limiter.Acquire(1, 10) {
		t.Fatal("expected acquire to succeed after release")
	}
}

func TestKDSConnectionLimiterIsolatedPerRoom(t *testing.T) {
	t.Parallel()

	limiter := newKDSConnectionLimiter(1)

	if !limiter.Acquire(1, 10) {
		t.Fatal("expected first room acquire to succeed")
	}
	if !limiter.Acquire(1, 11) {
		t.Fatal("expected different room acquire to succeed")
	}

	if limiter.Acquire(1, 10) {
		t.Fatal("expected same room acquire to be blocked")
	}
	if limiter.Acquire(1, 11) {
		t.Fatal("expected second room acquire to be blocked at max")
	}
}
