package worker

import (
	"context"
	"go-micro.dev/v4/logger"
	"sync"
	"time"
)

type Status int

const (
	Pending Status = iota
	Active
	Done
	Failed
)

type Receipt interface {
	Status() Status
	Cancel()
}

type receipt struct {
	t      Task
	l      logger.Logger
	ctx    context.Context
	cancel context.CancelFunc

	mu     sync.RWMutex
	status Status
}

func (s *Service) newReceipt(t Task) *receipt {
	r := receipt{
		t: t,
		l: s.l.Fields(map[string]interface{}{"task": t.ID()}),
	}
	r.ctx, r.cancel = context.WithCancel(s.ctx)
	return &r
}

func (r *receipt) Cancel() {
	r.cancel()
}

func (r *receipt) Status() Status {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.status
}

func (r *receipt) setStatus(status Status) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.l.Logf(logger.DebugLevel, "Set status = %d", status)
	r.status = status
}

func (r *receipt) run(timeout time.Duration) {
	defer r.cancel() // prevent context leak
	defer func() {
		if err := recover(); err != nil {
			r.l.Logf(logger.ErrorLevel, "Panic: %s", err)
			r.setStatus(Failed)
		}
	}()

	select {
	case <-r.ctx.Done():
		r.setStatus(Failed)
		r.l.Logf(logger.DebugLevel, "Skip cancelled task")
		return
	default:
	}

	ctx, cancel := context.WithTimeout(r.ctx, timeout)
	defer cancel()

	r.setStatus(Active)
	if err := r.t.Do(ctx); err != nil {
		r.l.Logf(logger.ErrorLevel, "Job failed: %s", err)
		r.setStatus(Failed)
	} else {
		r.setStatus(Done)
	}
}
