package transcoder

import (
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/RacoonMediaServer/rms-transcoder/internal/worker"
)

func (s *Service) runTranscodingTask(record *jobRecord) {
	job := record.job
	task := transcodingTask{
		l:           s.l.Fields(map[string]interface{}{"job": job.JobID}),
		id:          job.JobID,
		profile:     record.job.Profile,
		source:      job.Source,
		destination: job.Destination,
	}
	record.receipt = s.Workers.Do(&task)
}

func convertStatus(status worker.Status) rms_transcoder.GetJobResponse_Status {
	switch status {
	case worker.Pending:
		return rms_transcoder.GetJobResponse_Pending
	case worker.Active:
		return rms_transcoder.GetJobResponse_Processing
	case worker.Failed:
		return rms_transcoder.GetJobResponse_Failed
	case worker.Done:
		return rms_transcoder.GetJobResponse_Done
	default:
		return rms_transcoder.GetJobResponse_Failed
	}
}
