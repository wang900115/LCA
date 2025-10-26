package metric

import (
	"math"
	"math/rand"
	"slices"
	"sync"
	"time"
)

const (
	rescaleThreshold = time.Hour
)

// Sample is an interface for calculating statistical samples.
type Sample interface {
	Snapshot() *sampleSnapshot
	Clear()
	Update(int64)
}

// sampleSnapshot is a snapshot of a Sample at a point in time.
type sampleSnapshot struct {
	count  int64
	values []int64

	max      int64
	min      int64
	mean     float64
	sum      int64
	variance float64
}

// Create a new sample snapshot with precalculated statistics.
func newSampleSnapshotPrecalculated(count int64, values []int64, min, max, sum int64) *sampleSnapshot {
	if len(values) == 0 {
		return &sampleSnapshot{
			count:  count,
			values: values,
		}
	}
	return &sampleSnapshot{
		count:  count,
		values: values,
		max:    max,
		min:    min,
		mean:   float64(sum) / float64(len(values)),
		sum:    sum,
	}
}

// Create a new sample snapshot calculating statistics from values.
func newSampleSnapshot(count int64, values []int64) *sampleSnapshot {
	var (
		max int64 = math.MinInt64
		min int64 = math.MaxInt64
		sum int64
	)
	for _, v := range values {
		sum += v
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return newSampleSnapshotPrecalculated(count, values, min, max, sum)
}

// Count returns the count of inputs at the time the snapshot was taken.
func (s *sampleSnapshot) Count() int64 { return s.count }

// Max returns the maximum value in the sample.
func (s *sampleSnapshot) Max() int64 { return s.max }

// Min returns the minimum value in the sample.
func (s *sampleSnapshot) Min() int64 { return s.min }

// Mean returns the mean of the values in the sample.
func (s *sampleSnapshot) Mean() float64 { return s.mean }

// Percentile returns the value at the given percentile.
func (s *sampleSnapshot) Percentile(p float64) float64 {
	return SamplePercentile(s.values, p)
}

// Percentiles returns the values at the given percentiles.
func (s *sampleSnapshot) Percentiles(ps []float64) []float64 {
	return CalculatePercentiles(s.values, ps)
}

// Size returns the number of values in the sample.
func (s *sampleSnapshot) Size() int { return len(s.values) }

// Sum returns the sum of the values in the sample.
func (s *sampleSnapshot) Sum() int64 { return s.sum }

// StdDev returns the standard deviation of the values in the sample.
func (s *sampleSnapshot) StdDev() float64 {
	if s.variance == 0.0 {
		s.variance = SampleVariance(s.mean, s.values)
	}
	return math.Sqrt(s.variance)
}

// Values returns a copy of the values in the sample.
func (s *sampleSnapshot) Values() []int64 {
	return slices.Clone(s.values)
}

// Variance returns the variance of the values in the sample.
func (s *sampleSnapshot) Variance() float64 {
	if s.variance == 0.0 {
		s.variance = SampleVariance(s.mean, s.values)
	}
	return s.variance
}

type ExpDecaySample struct {
	alpha         float64
	count         int64
	mutex         sync.Mutex
	reservoirSize int
	t0, t1        time.Time
	values        *expDecaySampleHeap
	rand          *rand.Rand
}

func NewExpDecaySample(size int, alpha float64) Sample {
	s := &ExpDecaySample{
		alpha:         alpha,
		reservoirSize: size,
		t0:            time.Now(),
		values:        NewExpDecaySampleHeap(size),
	}
	s.t1 = s.t0.Add(rescaleThreshold)
	return s
}

func (s *ExpDecaySample) SetRand(prng *rand.Rand) Sample {
	s.rand = prng
	return s
}

func (s *ExpDecaySample) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count = 0
	s.t0 = time.Now()
	s.t1 = s.t0.Add(rescaleThreshold)
	s.values.Clear()
}

func (s *ExpDecaySample) Snapshot() *sampleSnapshot {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var (
		samples       = s.values.Values()
		values        = make([]int64, len(samples))
		max     int64 = math.MinInt64
		min     int64 = math.MaxInt64
		sum     int64
	)
	for i, item := range samples {
		v := item.v
		values[i] = v
		sum += v
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return newSampleSnapshotPrecalculated(s.count, values, min, max, sum)
}

func (s *ExpDecaySample) Update(v int64) {
	s.update(time.Now(), v)
}

func (s *ExpDecaySample) update(t time.Time, v int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count++
	if s.values.Size() == s.reservoirSize {
		s.values.Pop()
	}
	var f64 float64
	if s.rand != nil {
		f64 = s.rand.Float64()
	} else {
		f64 = rand.Float64()
	}
	s.values.Push(expDecaySample{
		k: math.Exp(t.Sub(s.t0).Seconds()*s.alpha) / f64,
		v: v,
	})
	if t.After(s.t1) {
		values := s.values.Values()
		t0 := s.t0
		s.values.Clear()
		s.t0 = t
		s.t1 = s.t0.Add(rescaleThreshold)
		for _, v := range values {
			v.k = v.k * math.Exp(-s.alpha*s.t0.Sub(t0).Seconds())
			s.values.Push(v)
		}
	}
}

// Calculate the variance of a sample of values given the mean.
func SampleVariance(mean float64, values []int64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	var sum float64
	for _, v := range values {
		diff := float64(v) - mean
		sum += diff * diff
	}
	return sum / float64(len(values))
}

// Calculate a single percentile from a sample of values.
func SamplePercentile(values []int64, p float64) float64 {
	return CalculatePercentiles(values, []float64{p})[0]
}

// Hazen Method Interpolation to calculate percentiles.
func CalculatePercentiles(values []int64, ps []float64) []float64 {
	scores := make([]float64, len(ps))
	size := len(values)
	if size == 0 {
		return scores
	}
	slices.Sort(values)
	for i, p := range ps {
		pos := p * float64(size+1)
		if pos < 1.0 {
			scores[i] = float64(values[0])
		} else if pos >= float64(size) {
			scores[i] = float64(values[size-1])
		} else {
			lower := values[int(pos)-1]
			upper := values[int(pos)]
			fraction := pos - float64(int(pos))
			scores[i] = float64(lower) + fraction*float64(upper-lower)
		}
	}
	return scores
}

type UniformSample struct {
	count         int64
	mutex         sync.Mutex
	reservoirSize int
	values        []int64
	rand          *rand.Rand
}

func NewUniformSample(size int) Sample {
	return &UniformSample{
		reservoirSize: size,
		values:        make([]int64, 0, size),
	}
}

func (s *UniformSample) SetRand(prng *rand.Rand) Sample {
	s.rand = prng
	return s
}

func (s *UniformSample) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count = 0
	clear(s.values)
}

func (s *UniformSample) Snapshot() *sampleSnapshot {
	s.mutex.Lock()
	values := slices.Clone(s.values)
	count := s.count
	s.mutex.Unlock()
	return newSampleSnapshot(count, values)
}

func (s *UniformSample) Update(v int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count++
	if len(s.values) < s.reservoirSize {
		s.values = append(s.values, v)
		return
	}
	var r int64
	if s.rand != nil {
		r = s.rand.Int63n(s.count)
	} else {
		r = rand.Int63n(s.count)
	}
	if r < int64(len(s.values)) {
		s.values[int(r)] = v
	}
}

type expDecaySample struct {
	k float64
	v int64
}

type expDecaySampleHeap struct {
	s []expDecaySample
}

func NewExpDecaySampleHeap(reservoirSize int) *expDecaySampleHeap {
	return &expDecaySampleHeap{
		s: make([]expDecaySample, 0, reservoirSize),
	}
}

func (h *expDecaySampleHeap) Clear() {
	h.s = h.s[:0]
}

func (h *expDecaySampleHeap) Push(s expDecaySample) {
	n := len(h.s)
	h.s = h.s[0 : n+1]
	h.s[n] = s
	h.up(n)
}

func (h *expDecaySampleHeap) Pop() expDecaySample {
	n := len(h.s) - 1
	h.s[0], h.s[n] = h.s[n], h.s[0]
	h.down(0, n)

	n = len(h.s)
	s := h.s[n-1]
	h.s = h.s[0 : n-1]
	return s
}

func (h *expDecaySampleHeap) Size() int {
	return len(h.s)
}

func (h *expDecaySampleHeap) Values() []expDecaySample {
	return h.s
}

func (h *expDecaySampleHeap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !(h.s[j].k < h.s[i].k) {
			break
		}
		h.s[i], h.s[j] = h.s[j], h.s[i]
		j = i
	}
}

func (h *expDecaySampleHeap) down(i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !(h.s[j1].k < h.s[j2].k) {
			j = j2 // = 2*i + 2  // right child
		}
		if !(h.s[j].k < h.s[i].k) {
			break
		}
		h.s[i], h.s[j] = h.s[j], h.s[i]
		i = j
	}
}
