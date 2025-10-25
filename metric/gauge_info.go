package metric

import (
	"encoding/json"
	"sync"
)

type GaugeInfoValue map[string]string

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

func (g GaugeInfoSnapshot) Value() GaugeInfoValue { return GaugeInfoValue(g) }

type GaugeInfo struct {
	mutex sync.Mutex
	value GaugeInfoValue
}

func (g *GaugeInfo) Snapshot() GaugeInfoSnapshot {
	return GaugeInfoSnapshot(g.value)
}

func (g *GaugeInfo) Update(value GaugeInfoValue) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value = value
}
