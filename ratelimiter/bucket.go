package ratelimiter

import (
	"errors"
	"time"
)

type Bucket struct {
	Capacity    int
	Count       int
	RefillRate  int
	LastRefill  time.Time
}

func NewBucket(capacity int, count int, rValue int, rRate int) *Bucket {
	return &Bucket{
		Capacity:    capacity,
		Count:       count,
		RefillRate:  rRate,
		LastRefill:  time.Now(),
	}
}

func (b *Bucket) Refill() {
	now := time.Now()
	elapsedTime := now.Sub(b.LastRefill)
	if elapsedTime.Seconds() >= 1 {
		b.Count += b.RefillRate
		b.LastRefill = now
	}
}

func (b *Bucket) UseToken() error {
	if !b.HasToken() {
		return errors.New("bucket has no tokens.")
	}

	b.Count--
	return nil
}

func (b *Bucket) HasToken() bool {
	return b.Count > 0
}
