package vone

import (
	"errors"
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
	//log.Println("XXX *AdaptiveRateLimiter) ShouldAbort() sleep = ", l.sleep)
	if l.stop == nil {
		time.Sleep(l.sleep)
		return false
	}
	select {
	case <-l.stop:
		return true
	case <-time.After(l.sleep):
		return false
	}
}

func (l *AdaptiveRateLimiter) CheckError(err error) error {
	//log.Println("AdaptiveRateLimiter) CheckError: ", err)
	if err == nil || !l.rateLimitSurpassed(err) {
		//	log.Println("nil or not )rateLimitSurpassed")
		l.sleep /= 2
		if l.sleep < minSleepTime {
			l.sleep = minSleepTime
		}
		//	log.Println("now sleep:", l.sleep, "return nil")
		return err
	}
	//log.Println("not nil")
	l.sleep *= 2
	if l.sleep > maxSleepTime {
		l.sleep = maxSleepTime
		//	log.Println("return ", err)
		return err
	}
	//log.Println("return ", ErrOnceMore)
	return ErrOnceMore
}
