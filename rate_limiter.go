package vone

import (
	"errors"
	"sync"
	"time"
)

type RateLimiter interface {
	ShouldAbort() bool
	CheckError(err error) error
}

const (
	minSleepTime = 1 * time.Microsecond
	maxSleepTime = 120 * time.Second
)

type RateLimitSurpassed = func(error) bool

type AdaptiveRateLimiter struct {
	stop               chan struct{}
	mu                 sync.Mutex
	sleep              time.Duration
	rateLimitSurpassed RateLimitSurpassed
}

var _ RateLimiter = &AdaptiveRateLimiter{}

func NewAdaptiveRateLimiter(rateLimitSurpassed RateLimitSurpassed, stop chan struct{}) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		stop:               stop,
		sleep:              minSleepTime,
		rateLimitSurpassed: rateLimitSurpassed,
	}
}

var (
	ErrOnceMore = errors.New("once more")
	ErrStop     = errors.New("stop")
)

func (l *AdaptiveRateLimiter) ShouldAbort() bool {
	l.mu.Lock()
	sleep := l.sleep
	l.mu.Unlock()

	if l.stop == nil {
		time.Sleep(sleep)
		return false
	}
	select {
	case <-l.stop:
		return true
	case <-time.After(sleep):
		return false
	}
}

func (l *AdaptiveRateLimiter) CheckError(err error) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err == nil || !l.rateLimitSurpassed(err) {
		l.sleep /= 2
		if l.sleep < minSleepTime {
			l.sleep = minSleepTime
		}
		return err
	}
	l.sleep *= 2
	if l.sleep > maxSleepTime {
		l.sleep = maxSleepTime
		return err
	}
	return ErrOnceMore
}
