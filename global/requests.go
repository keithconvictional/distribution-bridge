package global

import "time"

type RequestManager struct {
	CurrentTime time.Time
	Requests int
}

func (r *RequestManager) Wait() {
	now := time.Now()
	if r.CurrentTime.Second() != now.Second() {
		// New second so reset
		r.Requests = 0
		r.CurrentTime = now
	}
	if r.Requests > 3 {
		// Pause for a new second
		time.Sleep(1 * time.Second)
	}
}
