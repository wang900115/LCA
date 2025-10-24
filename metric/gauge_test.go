package metric

import "testing"

func TestGaugeSnapShot(t *testing.T) {
	g := NewGauge()
	g.Update(10)
	snapshot := g.Snapshot()
	g.Update(0)
	if snapshot.Value() != 10 {
		t.Errorf("expected snapshot value 10, got %d", snapshot.Value())
	}
}
