package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenLimit struct {
	rate   float64 // 速率，往桶里放入的数量
	bucket int     // 桶的大小
	mu     sync.Mutex
	total  float64   // 桶里面现在的数量
	last   time.Time // 上一次消耗的token的时间
}

func NewTokenLimit(rate float64, bucket int) *TokenLimit {
	return &TokenLimit{rate: rate, bucket: bucket}
}

func (this *TokenLimit) Allow() bool {
	return this.AllowN(time.Now(), 1)
}

func (this *TokenLimit) AllowN(now time.Time, n int) bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	// 计算上一次请求的补充了多少token
	delta := now.Sub(this.last).Seconds() * this.rate
	this.total += delta

	if this.total > float64(this.bucket) {
		this.total = float64(this.bucket)
	}

	if this.total < float64(n) {
		return false
	}

	this.total -= float64(n)
	this.last = now

	return true
}

func main() {
	limit := NewTokenLimit(1, 5)
	for {
		n := 3
		for i := 0; i < n; i++ {
			go func(i int) {
				if !limit.Allow() {
					fmt.Printf("forbid [%d] \n", i)
				} else {
					fmt.Printf("allow [%d] \n", i)
				}
			}(i)
		}
		time.Sleep(time.Second)
		fmt.Println("=============================")
	}
}
