package worker

import (
	"context"
	"errors"
	"fmt"
	"go-micro.dev/v4/logger"
	"sync"
)

type Status int

const (
	Pending = iota
	Active
	Done
	Failed
	Cancelled
)

type Receipt interface {
	Status() (Status, error)
	Cancel()
}

type receipt struct {
	t      Task
	l      logger.Logger
	ctx    context.Context
	cancel context.CancelFunc

	mu     sync.RWMutex
	status Status
	err    error
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
	r.setStatus(Cancelled, nil)
}

func (r *receipt) Status() (Status, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.status, r.err
}

func (r *receipt) setStatus(s Status, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.l.Logf(logger.DebugLevel, "Set status = %d, err = %s", s, err)
	r.status, r.err = s, err
}

func (r *receipt) run() {
	defer func() {
		if err := recover(); err != nil {
			r.setStatus(Failed, fmt.Errorf("%+v", err))
		}
	}()

	select {
	case <-r.ctx.Done():
		r.l.Logf(logger.DebugLevel, "Skip cancelled task")
		return
	default:
	}

	r.setStatus(Active, nil)
	if err := r.t.Do(r.ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			r.setStatus(Failed, err)
		}
	} else {
		r.setStatus(Done, nil)
	}
	r.cancel()
}
