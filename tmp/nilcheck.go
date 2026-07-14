package main

import (
	"fmt"
	"time"
)

type Bucket struct {
	Count      int
	LastRefill time.Time
}

func (b *Bucket) Refill() {
	now := time.Now()
	_ = now.Sub(b.LastRefill) // dereference nil?
	fmt.Println("Refill ran without panic")
}

func main() {
	m := make(map[string]*Bucket)
	bucket, exists := m["x"]
	if !exists {
		m["x"] = &Bucket{Count: 5, LastRefill: time.Now()}
	}
	fmt.Printf("bucket is nil? %v\n", bucket == nil)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC:", r)
		}
	}()
	bucket.Refill()
}
