package core

import "time"

type retryPolicy struct {
	count            int
	shouldRetryCount int
	timeout          time.Duration
}

func NewRetryPolicy(shouldRetryCount int, timeout time.Duration) *retryPolicy {
	return &retryPolicy{
		count:            0,
		shouldRetryCount: shouldRetryCount,
		timeout:          timeout,
	}
}

func (rp *retryPolicy) Retry() bool {
	if rp.count > rp.shouldRetryCount {
		rp.count = 0
		return false
	}

	rp.count++
	return true
}

func (rp *retryPolicy) Clean() {
	rp.count = 0
}

func (rp *retryPolicy) Timeout() {
	time.Sleep(rp.timeout)
}
