package worker

import (
	"context"
	"go-micro.dev/v4/logger"
	"sync"
)

const maxTasksPerWorker = 1000

type Service struct {
	l      logger.Logger
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
	q      chan *receipt
	done   chan string
}

func New(workers uint) *Service {
	s := &Service{
		q:    make(chan *receipt, workers*maxTasksPerWorker),
		done: make(chan string, workers*maxTasksPerWorker),
		l:    logger.DefaultLogger.Fields(map[string]interface{}{"from": "worker"}),
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.wg.Add(int(workers))
	for i := uint(0); i < workers; i++ {
		go func() {
			defer s.wg.Done()
			s.workerProcess()
		}()
	}

	return s
}

func (s *Service) Do(t Task) Receipt {
	r := s.newReceipt(t)
	s.q <- r
	return r
}

func (s *Service) DoneChannel() chan<- string {
	return s.done
}

func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	close(s.q)
	close(s.done)
}
