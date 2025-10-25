package metric

import (
	"sync"
	"time"
)

// ResettingTimer is a timer that resets its values after each snapshot.
type ResettingTimer struct {
	// values holds the recorded durations.
	values []int64
	// sum holds the total duration.
	sum   int64
	mutex sync.Mutex
}

func NewResettingTimer() *ResettingTimer {
	return &ResettingTimer{
		values: make([]int64, 0, 10),
	}
}

// Snapshot returns a snapshot of the current timer and resets its values.
func (t *ResettingTimer) Snapshot() *ResettingTimerSnapshot {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	snapshot := &ResettingTimerSnapshot{}
	if len(t.values) > 0 {
		snapshot.mean = float64(t.sum) / float64(len(t.values))
		snapshot.values = t.values
		t.values = make([]int64, 0, 10)
	}
	t.sum = 0
	return snapshot
}

// Time measures the duration of the function f and records it.
func (t *ResettingTimer) Time(f func()) {
	start := time.Now()
	f()
	t.Update(time.Since(start))
}

// Update records a duration d.
func (t *ResettingTimer) Update(d time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.values = append(t.values, int64(d))
	t.sum += int64(d)
}

// UpdateSince records the duration since the given timestamp ts.
func (t *ResettingTimer) UpdateSince(ts time.Time) {
	t.Update(time.Since(ts))
}

type ResettingTimerSnapshot struct {
	// values holds the recorded durations.
	values []int64
	// mean holds the mean duration.
	mean float64
	// max holds the maximum duration.
	max int64
	// min holds the minimum duration.
	min int64
	// thresholdBoundaries holds the calculated percentiles.
	thresholdBoundaries []float64
	calculated          bool
}

func (t *ResettingTimerSnapshot) Count() int {
	return len(t.values)
}

func (t *ResettingTimerSnapshot) Percentiles(percentiles []float64) []float64 {
	t.calc(percentiles)
	return t.thresholdBoundaries
}

func (t *ResettingTimerSnapshot) Mean() float64 {
	if !t.calculated {
		t.calc(nil)
	}
	return t.mean
}

func (t *ResettingTimerSnapshot) Max() int64 {
	if !t.calculated {
		t.calc(nil)
	}
	return t.max
}

func (t *ResettingTimerSnapshot) Min() int64 {
	if !t.calculated {
		t.calc(nil)
	}
	return t.min
}

func (t *ResettingTimerSnapshot) calc(percentiles []float64) {
	scores := CalculatePercentiles(t.values, percentiles)
	t.thresholdBoundaries = scores
	if len(t.values) == 0 {
		return
	}
	t.min = t.values[0]
	t.max = t.values[len(t.values)-1]
	t.calculated = true
}
