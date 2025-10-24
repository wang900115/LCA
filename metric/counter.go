package metric

import "sync/atomic"

func NewCounter() *Counter {
	return new(Counter)
}

type CounterSnapshot int64

func (c CounterSnapshot) Count() int64 {
	return int64(c)
}

type Counter atomic.Int64

func (c *Counter) Clear() {
	(*atomic.Int64)(c).Store(0)
}

func (c *Counter) Inc(delta int64) {
	(*atomic.Int64)(c).Add(delta)
}

func (c *Counter) Dec(delta int64) {
	(*atomic.Int64)(c).Add(-delta)
}

func (c *Counter) Snapshot() CounterSnapshot {
	return CounterSnapshot((*atomic.Int64)(c).Load())
}
