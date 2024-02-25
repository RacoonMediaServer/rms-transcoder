package worker

func (s *Service) workerProcess() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case r := <-s.q:
			r.run()
		}
	}
}
