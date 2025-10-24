package metric

import (
	"math"
	"sync/atomic"
)

func NewCounterFloat64() *CounterFloat64 {
	return new(CounterFloat64)
}

type CounterFloat64Snapshot float64

func (c CounterFloat64Snapshot) Count() float64 {
	return float64(c)
}

type CounterFloat64 atomic.Uint64

func (c *CounterFloat64) Clear() {
	(*atomic.Uint64)(c).Store(0)
}

func (c *CounterFloat64) Inc(delta float64) {
	atomicAddFloat((*atomic.Uint64)(c), delta)
}

func (c *CounterFloat64) Dec(delta float64) {
	atomicAddFloat((*atomic.Uint64)(c), -delta)
}

func (c *CounterFloat64) Snapshot() CounterFloat64Snapshot {
	return CounterFloat64Snapshot(math.Float64frombits((*atomic.Uint64)(c).Load()))
}

func atomicAddFloat(fbits *atomic.Uint64, delta float64) {
	for {
		oldBits := fbits.Load()
		newBits := math.Float64bits(math.Float64frombits(oldBits) + delta)
		if fbits.CompareAndSwap(oldBits, newBits) {
			return
		}
	}
}
