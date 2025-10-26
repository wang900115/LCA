package metric

type HistogramSnapshot interface {
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(p float64) float64
	Percentiles(ps []float64) []float64
	Size() int
	Sum() int64
	StdDev() float64
	Variance() float64
}

type Histogram interface {
	Clear()
	Snapshot() HistogramSnapshot
	Update(value int64)
}

func NewHistogram(s Sample) Histogram {
	return &StandardHistogram{s}
}

type StandardHistogram struct {
	sample Sample
}

func (h *StandardHistogram) Clear() {
	h.sample.Clear()
}

func (h *StandardHistogram) Snapshot() HistogramSnapshot {
	return h.sample.Snapshot()
}

func (h *StandardHistogram) Update(value int64) {
	h.sample.Update(value)
}
