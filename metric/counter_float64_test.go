package metric

import "testing"

func TestCounterFloat64(t *testing.T) {
	c := NewCounterFloat64()
	if got := c.Snapshot().Count(); got != 0 {
		t.Errorf("initial count = %v; want 0", got)
	}
	c.Inc(2.5)
	if got := c.Snapshot().Count(); got != 2.5 {
		t.Errorf("after Inc(2.5), count = %v; want 2.5", got)
	}
	c.Dec(1.0)
	if got := c.Snapshot().Count(); got != 1.5 {
		t.Errorf("after Dec(1.0), count = %v; want 1.5", got)
	}
	c.Clear()
	if got := c.Snapshot().Count(); got != 0 {
		t.Errorf("after Clear(), count = %v; want 0", got)
	}
}
