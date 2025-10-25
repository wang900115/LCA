package metric

import "sync/atomic"

// GaugeSnapshot is a read-only copy of a Gauge metric.
type GaugeSnapshot int64

// Value returns the value of the GaugeSnapshot.
func (g GaugeSnapshot) Value() int64 { return int64(g) }

func NewGauge() *Gauge {
	return new(Gauge)
}

type Gauge atomic.Int64

func (g *Gauge) Update(value int64) {
	(*atomic.Int64)(g).Store(value)
}

// Snapshot returns a read-only copy of the current value of the Gauge.
func (g *Gauge) Snapshot() GaugeSnapshot {
	return GaugeSnapshot((*atomic.Int64)(g).Load())
}

// UpdateIfGt updates the gauge to v if v is greater than the current value.
func (g *Gauge) UpdateIfGt(v int64) {
	value := (*atomic.Int64)(g)
	for {
		exist := value.Load()
		if v <= exist {
			return
		}
		if value.CompareAndSwap(exist, v) {
			return
		}
	}
}

func (g *Gauge) Inc(delta int64) {
	(*atomic.Int64)(g).Add(delta)
}

func (g *Gauge) Dec(delta int64) {
	(*atomic.Int64)(g).Add(-delta)
}
