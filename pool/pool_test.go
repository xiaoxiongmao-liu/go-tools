package pool

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormal(t *testing.T) {
	a := atomic.Int64{}
	f := func(i int) {
		time.Sleep(1 * time.Millisecond)
		a.Add(int64(i))
	}
	option := WithWorkerFunc[int](f)
	p := NewPool[int](option)
	defer p.Close()
	for i := 0; i < 10; i++ {
		p.Put(i)
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(2 * time.Second)
	assert.Equal(t, 45, int(a.Load()))
}

func TestBlock(t *testing.T) {
	a := atomic.Int64{}
	f := func(i int) {
		time.Sleep(1 * time.Second)
		a.Add(int64(i))
	}
	options := []Option[int]{WithWorkerFunc[int](f), WithQueueSize[int](1), WithWorkerPoolSize[int](3), WithQueueFullStrategy[int](QueueFullStrategyBlock)}
	p := NewPool[int](options...)
	defer p.Close()
	for i := 0; i < 10; i++ {
		p.Put(i)
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(11 * time.Second)
	assert.Equal(t, 45, int(a.Load()))
}

func TestDropLatest(t *testing.T) {
	a := atomic.Int64{}
	f := func(i int) {
		time.Sleep(1 * time.Second)
		a.Add(int64(i))
	}
	options := []Option[int]{WithWorkerFunc[int](f), WithQueueSize[int](1), WithWorkerPoolSize[int](3), WithQueueFullStrategy[int](QueueFullStrategyDropLatest)}
	p := NewPool[int](options...)
	defer p.Close()
	for i := 0; i < 10; i++ {
		p.Put(i)
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(2 * time.Second)

	assert.Equal(t, 6, int(a.Load()))
}

func TestDropOldest(t *testing.T) {
	a := atomic.Int64{}
	f := func(i int) {
		time.Sleep(1 * time.Second)
		a.Add(int64(i))
	}
	options := []Option[int]{WithWorkerFunc[int](f), WithQueueSize[int](1), WithWorkerPoolSize[int](3), WithQueueFullStrategy[int](QueueFullStrategyDropOldest)}
	p := NewPool[int](options...)
	defer p.Close()
	for i := 0; i < 10; i++ {
		p.Put(i)
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(2 * time.Second)

	assert.Equal(t, 12, int(a.Load()))
}
