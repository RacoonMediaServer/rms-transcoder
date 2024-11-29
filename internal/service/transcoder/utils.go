package transcoder

import (
	"net/url"
	"os"
	"path/filepath"

	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/RacoonMediaServer/rms-packages/pkg/worker"
	"go-micro.dev/v4/logger"
)

func (s *Service) runTranscodingTask(record *jobRecord) {
	job := record.job
	source := job.Source
	if isFileSource(source) {
		source = filepath.Join(s.Config.Directory, source)
	}
	destination := filepath.Join(s.Config.Directory, job.Destination)

	settings := record.job.Transcoding.Data()
	task := transcodingTask{
		l:           s.l.Fields(map[string]interface{}{"job": job.JobID}),
		id:          job.JobID,
		settings:    &settings,
		source:      source,
		destination: destination,
		dur:         job.Duration,
		offset:      job.Offset,
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

func clearContent(directory, destination string) error {
	path := filepath.Join(directory, destination)
	_, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return os.Remove(path)
}

func getFileSize(file string) int64 {
	fi, err := os.Stat(file)
	if err != nil {
		return 0
	}
	return fi.Size()
}

func isFileSource(source string) bool {
	u, err := url.Parse(source)
	if err != nil {
		logger.Warnf("Possible incorrect URL: %s", err)
		return true
	}

	return u.Scheme == "" || u.Scheme == "file"
}
