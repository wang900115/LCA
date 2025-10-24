package metric

func NewHealthcheck(f func(*Healthcheck)) *Healthcheck {
	return &Healthcheck{nil, f}
}

type Healthcheck struct {
	err error
	f   func(*Healthcheck)
}

func (h *Healthcheck) Check() {
	h.f(h)
}

func (h *Healthcheck) Err() error {
	return h.err
}

func (h *Healthcheck) Healthy() {
	h.err = nil
}

func (h *Healthcheck) Unhealthy(err error) {
	h.err = err
}
