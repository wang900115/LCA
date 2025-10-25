package metric

import "testing"

func TestGaugeFloat64Snapshot(t *testing.T) {
	g := NewGaugeFloat64()
	g.Update(3.14)
	snapshot := g.Snapshot()
	g.Update(2.71)

	if snapshot.Value() != 3.14 {
		t.Errorf("expected snapshot value 3.14, got %v", snapshot.Value())
	}
	if g.Snapshot().Value() != 2.71 {
		t.Errorf("expected current gauge value 2.71, got %v", g.Snapshot().Value())
	}
}
