package metric

import (
	"errors"
	"sync"
)

// todo
var ErrDuplicateMetric = errors.New("duplicate metric")

type Registry interface {

	// Each call the given function for each registered metric.
	Each(func(string, interface{}))

	// Get retrieves a registered metric by name.
	Get(name string) interface{}

	// GetAll retrieves all registered metrics.
	GetAll() map[string]map[string]interface{}

	// GetOrRegister retrieves a registered metric by name or registers a new one.
	GetOrRegister(name string, metric interface{}) interface{}

	// Register registers a new metric by name.
	Register(name string, metric interface{}) error

	// HealthChecks performs health checks on all registered metrics.
	HealthChecks()

	// UnRegister removes a registered metric by name.
	UnRegister(name string)
}

type StandardRegistry struct {
	metrics sync.Map
}

func (r *StandardRegistry) Each(f func(string, interface{})) {
	for name, i := range r.registered() {
		f(name, i)
	}
}

func (r *StandardRegistry) Get(name string) interface{} {
	value, _ := r.metrics.Load(name)
	return value
}

func (r *StandardRegistry) loadOrRegister(name string, i interface{}) (interface{}, bool, bool) {
	return nil, false, false
}

func (r *StandardRegistry) registered() map[string]interface{} {
	metrics := make(map[string]interface{})
	r.metrics.Range(func(key, value any) bool {
		metrics[key.(string)] = value
		return true
	})
	return metrics
}
