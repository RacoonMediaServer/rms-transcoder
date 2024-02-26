package transcoder

import (
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
)

func (s *Service) runTranscodingTask(record *jobRecord, profile *rms_transcoder.Profile) {
	job := record.job
	task := transcodingTask{
		l:           s.l.Fields(map[string]interface{}{"job": job.ID}),
		id:          job.ID,
		profile:     profile,
		source:      job.Source,
		destination: job.Destination,
	}
	record.receipt = s.Workers.Do(&task)
	s.jobs[job.ID] = record
}
