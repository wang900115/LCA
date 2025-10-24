package metric

import "testing"

func TestCounterClear(t *testing.T) {
	c := NewCounter()
	c.Inc(10)
	c.Clear()
	if count := c.Snapshot().Count(); count != 0 {
		t.Errorf("expected count to be 0 after Clear, got %d", count)
	}
}

func TestCounter(t *testing.T) {
	c := NewCounter()
	if count := c.Snapshot().Count(); count != 0 {
		t.Errorf("expected initial count to be 0, got %d", count)
	}
	c.Dec(1)
	if count := c.Snapshot().Count(); count != -1 {
		t.Errorf("expected count to be -1 after Decrementing by 1, got %d", count)
	}
	c.Inc(5)
	if count := c.Snapshot().Count(); count != 4 {
		t.Errorf("expected count to be 4 after Incrementing by 5, got %d", count)
	}
	c.Dec(2)
	if count := c.Snapshot().Count(); count != 2 {
		t.Errorf("expected count to be 2 after Decrementing by 2, got %d", count)
	}
	c.Inc(3)
	if count := c.Snapshot().Count(); count != 5 {
		t.Errorf("expected count to be 5 after Incrementing by 3, got %d", count)
	}
}

func TestCounterSnapTest(t *testing.T) {
	c := NewCounter()
	c.Inc(7)
	snapshot := c.Snapshot()
	c.Inc(1)
	if snapshot.Count() != 7 {
		t.Errorf("expected snapshot count to be 7, got %d", snapshot.Count())
	}
	if c.Snapshot().Count() != 8 {
		t.Errorf("expected current count to be 8 after Incrementing by 1, got %d", c.Snapshot().Count())
	}
}
