package metric

import (
	"math"
	"sync/atomic"
)

type GaugeFloat64Snapshot float64

func (g GaugeFloat64Snapshot) Value() float64 { return float64(g) }

func NewGaugeFloat64() *GaugeFloat64 {
	return new(GaugeFloat64)
}

type GaugeFloat64 atomic.Uint64

func (g *GaugeFloat64) Snapshot() GaugeFloat64Snapshot {
	return GaugeFloat64Snapshot(math.Float64frombits((*atomic.Uint64)(g).Load()))
}

func (g *GaugeFloat64) Update(v float64) {
	(*atomic.Uint64)(g).Store(math.Float64bits(v))
}
