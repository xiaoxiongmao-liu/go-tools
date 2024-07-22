package pool

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultIdleDuration = 10 * time.Second // 默认空闲协程超时时间

	defaultQueueSize         = 100 // 默认任务队列大小
	defaultMinIdleWorkerSize = 1   // 默认最小空闲协程数
	defaultWorkerPoolSize    = 16  // 默认协程池大小
)

type QueueFullStrategy string

const (
	QueueFullStrategyBlock      QueueFullStrategy = "block"       // 阻塞
	QueueFullStrategyDropLatest QueueFullStrategy = "drop_latest" // 丢弃最新，默认
	QueueFullStrategyDropOldest QueueFullStrategy = "drop_oldest" // 丢弃最早
)

type Pool[T any] struct {
	jobQueue  chan T            // 任务队列
	queueSize int               // 队列大小
	strategy  QueueFullStrategy // 任务队列满策略

	minIdleWorkerSize int64         // 最小空闲协程数，最小是1
	workerPoolSize    int64         // 协程池大小
	workerFunc        func(T)       // 协程处理函数
	idleDuration      time.Duration // 空闲协程超时时间

	workerCount *atomic.Int64 // 当前工作协程数量

	lock sync.Mutex     // 锁
	done chan struct{}  // 关闭信号
	wg   sync.WaitGroup // 关闭线程池等待
}

type Option[T any] func(*Pool[T])

func WithQueueSize[T any](queueSize int) Option[T] {
	return func(p *Pool[T]) {
		if queueSize > 0 {
			p.queueSize = queueSize
		}
	}
}

func WithMinIdleWorkerSize[T any](minIdleWorkerSize int64) Option[T] {
	return func(p *Pool[T]) {
		if minIdleWorkerSize > 0 {
			p.minIdleWorkerSize = minIdleWorkerSize
		}
	}
}

func WithWorkerPoolSize[T any](workerPoolSize int64) Option[T] {
	return func(p *Pool[T]) {
		if workerPoolSize > 0 {
			p.workerPoolSize = workerPoolSize
		}
	}
}

func WithWorkerFunc[T any](workerFunc func(T)) Option[T] {
	return func(p *Pool[T]) {
		if workerFunc != nil {
			p.workerFunc = workerFunc
		}
	}
}

func WithIdleDuration[T any](idleDuration time.Duration) Option[T] {
	return func(p *Pool[T]) {
		if idleDuration >= defaultIdleDuration {
			p.idleDuration = idleDuration
		}
	}
}

func WithQueueFullStrategy[T any](strategy QueueFullStrategy) Option[T] {
	return func(p *Pool[T]) {
		if strategy != "" {
			p.strategy = strategy
		}
	}
}

func defaultPool[T any]() *Pool[T] {
	return &Pool[T]{
		queueSize:         defaultQueueSize,
		minIdleWorkerSize: defaultMinIdleWorkerSize,
		workerPoolSize:    defaultWorkerPoolSize,
		idleDuration:      defaultIdleDuration,
		done:              make(chan struct{}),
		workerCount:       &atomic.Int64{},
	}
}

// NewPool 创建协程池
func NewPool[T any](options ...Option[T]) *Pool[T] {
	p := defaultPool[T]()
	for _, option := range options {
		option(p)
	}
	// 初始化任务队列
	p.jobQueue = make(chan T, p.queueSize)

	// 协程池大小不能小于最小空闲协程数
	if p.minIdleWorkerSize > p.workerPoolSize {
		p.workerPoolSize = p.minIdleWorkerSize
	}
	// 启动协程池
	p.start()
	return p
}

// start 启动协程池
func (p *Pool[T]) start() {
	for i := int64(0); i < p.minIdleWorkerSize; i++ {
		w := Worker[T]{
			ticker: time.NewTicker(p.idleDuration),
			pool:   p,
		}
		p.wg.Add(1)
		p.workerCount.Add(1)
		go w.run()
	}
}

func (p *Pool[T]) Put(t T) {
	select {
	case p.jobQueue <- t:
	default:
		p.lock.Lock()
		defer p.lock.Unlock()
		if p.workerCount.Load() < p.workerPoolSize {
			w := Worker[T]{
				ticker: time.NewTicker(p.idleDuration),
				pool:   p,
			}
			p.wg.Add(1)
			p.workerCount.Add(1)
			go w.run()
			p.jobQueue <- t
			return
		}
		switch p.strategy {
		case QueueFullStrategyDropOldest:
			<-p.jobQueue
			p.jobQueue <- t
		case QueueFullStrategyBlock:
			p.jobQueue <- t
		default:
			return
		}
	}
}

// Close 关闭协程池 优雅退出
func (p *Pool[T]) Close() {
	close(p.jobQueue)
	close(p.done)
	p.wg.Wait()
}

type Worker[T any] struct {
	ticker *time.Ticker

	pool *Pool[T]
}

func (w *Worker[T]) run() {
	defer w.pool.wg.Done()
	defer w.pool.workerCount.Add(-1)
	for {
		select {
		case t := <-w.pool.jobQueue:
			if w.pool.workerFunc != nil {
				w.pool.workerFunc(t)
			}
			w.ticker.Reset(w.pool.idleDuration)
		case <-w.ticker.C:
			if w.pool.workerCount.Load() > w.pool.minIdleWorkerSize {
				w.ticker.Stop()
				return
			}
			w.ticker.Reset(w.pool.idleDuration)
		case <-w.pool.done:
			w.ticker.Stop()
			return
		}
	}
}
