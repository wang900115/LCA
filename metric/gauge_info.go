package metric

import (
	"encoding/json"
	"sync"
)

type GaugeInfoValue map[string]string

// String returns the JSON representation of the GaugeInfoValue.
func (g GaugeInfoValue) String() string {
	data, _ := json.Marshal(g)
	return string(data)
}

func NewGaugeInfo() *GaugeInfo {
	return &GaugeInfo{
		value: GaugeInfoValue{},
	}
}

type GaugeInfoSnapshot GaugeInfoValue

// Value returns the value of the GaugeInfoSnapshot.
func (g GaugeInfoSnapshot) Value() GaugeInfoValue { return GaugeInfoValue(g) }

type GaugeInfo struct {
	mutex sync.Mutex
	value GaugeInfoValue
}

// Snapshot returns a read-only copy of the current value of the GaugeInfo.
func (g *GaugeInfo) Snapshot() GaugeInfoSnapshot {
	return GaugeInfoSnapshot(g.value)
}

// Update sets the value of the GaugeInfo.
func (g *GaugeInfo) Update(value GaugeInfoValue) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value = value
}
