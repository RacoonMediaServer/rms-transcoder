package worker

func (s *Service) workerProcess() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case r := <-s.q:
			// TODO: add deadline
			r.run()
			s.done <- r.t.ID()
		}
	}
}
