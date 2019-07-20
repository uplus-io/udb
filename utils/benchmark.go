/*
 * Copyright (c) 2019 uplus.io
 */

package utils

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Benchmark struct {
	Name       string
	Concurrent int
	Max        int

	startTime      time.Time
	finishDuration time.Duration

	DurationChannel chan time.Duration
	CountChannel    chan int64

	total    int64
	duration time.Duration
}

func NewBenchmark(name string, concurrent int, max int) *Benchmark {
	benchmark := &Benchmark{Name: name, Concurrent: concurrent, Max: max}
	benchmark.CountChannel = make(chan int64, concurrent)
	benchmark.DurationChannel = make(chan time.Duration, concurrent)
	return benchmark
}

func (p *Benchmark) Start() {
	p.startTime = time.Now()
}

func (p *Benchmark) Finish() {
	p.finishDuration = time.Since(p.startTime)
}

func (p Benchmark) RandomInt() []byte {
	return []byte(strconv.Itoa(rand.Intn(p.Max)))
}

func (p Benchmark) RandomInt16() []byte {
	key := fmt.Sprintf("%16d", rand.Intn(p.Max))
	return []byte(key)
}
func (p Benchmark) GenerateKey(i int) []byte {
	key := fmt.Sprintf("%16d", i)
	return []byte(key)
}

func (p *Benchmark) Put(d time.Duration, times int64) {
	//p.DurationChannel <- d
	//p.CountChannel <- 1
	p.total += times
	p.duration += d
}

func (p *Benchmark) Watch() {
	for i := 0; i < p.Concurrent; i++ {
		p.total += <-p.CountChannel
		p.duration += <-p.DurationChannel
	}
}

func (p Benchmark) Print() {
	log.Printf("%s benchmark stat-------------------------\n", p.Name)
	log.Printf("start:%s duration:%s", p.startTime, p.finishDuration)
	log.Printf("total:%d", p.total)
	log.Printf("average time cost: %.6fms", float64((p.duration / time.Millisecond))/float64(p.total))
	qps := float64(p.total*1000) / float64(p.duration/time.Millisecond)
	log.Printf("qps: %.2f\n\n", qps)
}
