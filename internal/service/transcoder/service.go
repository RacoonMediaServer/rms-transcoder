package transcoder

import (
	"context"
	"errors"
	"fmt"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/RacoonMediaServer/rms-transcoder/internal/model"
	"github.com/RacoonMediaServer/rms-transcoder/internal/worker"
	"github.com/google/uuid"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type Service struct {
	l        logger.Logger
	Profiles ProfileService
	Database Database
	Workers  Workers

	mu   sync.RWMutex
	jobs map[string]*jobRecord
}

type jobRecord struct {
	job     *model.Job
	receipt worker.Receipt
}

func (s *Service) Initialize() error {
	s.l = logger.DefaultLogger.Fields(map[string]interface{}{"from": "transcoder"})
	s.jobs = make(map[string]*jobRecord)

	jobs, err := s.Database.LoadJobs()
	if err != nil {
		return err
	}

	s.l.Logf(logger.InfoLevel, "Rerun loaded jobs...")

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, job := range jobs {
		record := jobRecord{job: job}
		if job.Result == model.JobNotComplete {
			profile, err := s.Profiles.GetProfile(job.Profile)
			if err != nil || profile == nil {
				if err == nil {
					err = errors.New("profile not found")
				}
				s.l.Logf(logger.ErrorLevel, "Cannot rerun task %s: %s", job.ID, err)
				job.Result = model.JobFailed
			} else {
				s.runTranscodingTask(&record, profile)
				continue
			}
		}
		s.jobs[job.ID] = &record
	}
	return nil
}

func (s *Service) AddJob(ctx context.Context, request *rms_transcoder.AddJobRequest, response *rms_transcoder.AddJobResponse) error {
	id, err := uuid.NewUUID()
	if err != nil {
		s.l.Logf(logger.ErrorLevel, "Generate id failed: %s", err)
		return err
	}
	profile, err := s.Profiles.GetProfile(request.Profile)
	if err != nil || profile == nil {
		return fmt.Errorf("cannot use profile '%s': %w", request.Profile, err)
	}
	job := model.Job{
		ID:           id.String(),
		Profile:      request.Profile,
		Source:       request.Source,
		Destination:  request.Destination,
		AutoComplete: request.AutoComplete,
	}
	if err = s.Database.AddJob(&job); err != nil {
		s.l.Logf(logger.ErrorLevel, "Add job to database failed: %s", err)
		return err
	}
	response.JobId = job.ID

	s.mu.Lock()
	defer s.mu.Unlock()
	s.runTranscodingTask(&jobRecord{job: &job}, profile)
	return nil
}

func (s *Service) GetJob(ctx context.Context, request *rms_transcoder.GetJobRequest, response *rms_transcoder.GetJobResponse) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.jobs[request.JobId]
	if !ok {
		return errors.New("job not found")
	}

	j, r := record.job, record.receipt
	response.Destination = j.Destination
	if r != nil {
		status, _ := r.Status()
		switch status {
		case worker.Pending:
			response.Status = rms_transcoder.GetJobResponse_Pending
		case worker.Active:
			response.Status = rms_transcoder.GetJobResponse_Processing
		case worker.Done:
			response.Status = rms_transcoder.GetJobResponse_Done
		case worker.Cancelled:
			fallthrough
		case worker.Failed:
			response.Status = rms_transcoder.GetJobResponse_Failed
		}
	} else {
		switch j.Result {
		case model.JobNotComplete:
			response.Status = rms_transcoder.GetJobResponse_Pending
		case model.JobDone:
			response.Status = rms_transcoder.GetJobResponse_Done
		case model.JobFailed:
			response.Status = rms_transcoder.GetJobResponse_Failed
		}
	}

	return nil
}

func (s *Service) CancelJob(ctx context.Context, request *rms_transcoder.CancelJobRequest, empty *emptypb.Empty) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.jobs[request.JobId]
	if !ok {
		return errors.New("job not found")
	}

	r := record.receipt
	if r != nil {
		r.Cancel()
	}

	return nil
}
