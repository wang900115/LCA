package metric

import (
	"testing"
	"time"
)

func TestResettingTimer(t *testing.T) {
	tests := []struct {
		values   []int64
		start    int
		end      int
		wantP50  float64
		wantP95  float64
		wantP99  float64
		wantMean float64
		wantMin  int64
		wantMax  int64
	}{
		{
			values:  []int64{},
			start:   1,
			end:     11,
			wantP50: 5.5, wantP95: 10, wantP99: 10,
			wantMin: 1, wantMax: 10, wantMean: 5.5,
		},
		{
			values:  []int64{},
			start:   1,
			end:     101,
			wantP50: 50.5, wantP95: 95.94999999999999, wantP99: 99.99,
			wantMin: 1, wantMax: 100, wantMean: 50.5,
		},
		{
			values:  []int64{1},
			start:   0,
			end:     0,
			wantP50: 1, wantP95: 1, wantP99: 1,
			wantMin: 1, wantMax: 1, wantMean: 1,
		},
		{
			values:  []int64{0},
			start:   0,
			end:     0,
			wantP50: 0, wantP95: 0, wantP99: 0,
			wantMin: 0, wantMax: 0, wantMean: 0,
		},
		{
			values:  []int64{},
			start:   0,
			end:     0,
			wantP50: 0, wantP95: 0, wantP99: 0,
			wantMin: 0, wantMax: 0, wantMean: 0,
		},
		{
			values:  []int64{1, 10},
			start:   0,
			end:     0,
			wantP50: 5.5, wantP95: 10, wantP99: 10,
			wantMin: 1, wantMax: 10, wantMean: 5.5,
		},
	}
	for i, tt := range tests {
		timer := NewResettingTimer()

		for i := tt.start; i < tt.end; i++ {
			tt.values = append(tt.values, int64(i))
		}

		for _, v := range tt.values {
			timer.Update(time.Duration(v))
		}
		snap := timer.Snapshot()

		ps := snap.Percentiles([]float64{0.50, 0.95, 0.99})

		if have, want := snap.Min(), tt.wantMin; have != want {
			t.Fatalf("%d: min: have %d, want %d", i, have, want)
		}
		if have, want := snap.Max(), tt.wantMax; have != want {
			t.Fatalf("%d: max: have %d, want %d", i, have, want)
		}
		if have, want := snap.Mean(), tt.wantMean; have != want {
			t.Fatalf("%d: mean: have %v, want %v", i, have, want)
		}
		if have, want := ps[0], tt.wantP50; have != want {
			t.Errorf("%d: p50: have %v, want %v", i, have, want)
		}
		if have, want := ps[1], tt.wantP95; have != want {
			t.Errorf("%d: p95: have %v, want %v", i, have, want)
		}
		if have, want := ps[2], tt.wantP99; have != want {
			t.Errorf("%d: p99: have %v, want %v", i, have, want)
		}
	}
}
