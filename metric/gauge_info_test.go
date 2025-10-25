package metric

import "testing"

func TestGaugeInfoJsonString(t *testing.T) {
	g := NewGaugeInfo()
	g.Update(GaugeInfoValue{"key1": "value1", "key2": "value2"})
	expected := `{"key1":"value1","key2":"value2"}`

	original := g.Snapshot()
	g.Update(GaugeInfoValue{"key1": "value_fix1"})

	if have := original.Value().String(); have != expected {
		t.Errorf("GaugeInfo.String() = %q, want %q", have, expected)
	}

	if have, want := g.Snapshot().Value().String(), `{"key1":"value_fix1"}`; have != want {
		t.Errorf("GaugeInfo.String() after update = %q, want %q", have, want)
	}
}
