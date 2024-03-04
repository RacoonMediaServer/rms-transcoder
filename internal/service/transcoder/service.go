package transcoder

import (
	"context"
	"errors"
	"fmt"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/RacoonMediaServer/rms-transcoder/internal/config"
	"github.com/RacoonMediaServer/rms-transcoder/internal/model"
	"github.com/RacoonMediaServer/rms-transcoder/internal/worker"
	"github.com/google/uuid"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/datatypes"
	"os"
	"sync"
)

type Service struct {
	Profiles  ProfileService
	Database  Database
	Workers   Workers
	Publisher micro.Event
	Config    config.Transcoding

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	l    logger.Logger
	mu   sync.RWMutex
	jobs map[string]*jobRecord
}

type jobRecord struct {
	job     *model.Job
	receipt worker.Receipt
}

func (s *Service) Initialize() error {
	fi, err := os.Stat(s.Config.Directory)
	if err != nil || !fi.IsDir() {
		return fmt.Errorf("problem with content directory: %w", err)
	}

	s.l = logger.DefaultLogger.Fields(map[string]interface{}{"from": "transcoder"})
	s.jobs = make(map[string]*jobRecord)
	s.ctx, s.cancel = context.WithCancel(context.Background())

	jobs, err := s.Database.LoadJobs()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	rerun := 0
	for _, job := range jobs {
		record := jobRecord{job: job}
		s.jobs[job.JobID] = &record
		if !job.Done {
			s.runTranscodingTask(&record)
			rerun++
		}
	}

	s.l.Logf(logger.InfoLevel, "Rerun jobs: %d / %d", rerun, len(jobs))

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.processReadyJobs()
	}()
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
		JobID:        id.String(),
		Transcoding:  datatypes.NewJSONType(*profile.Settings),
		Source:       request.Source,
		Destination:  request.Destination,
		AutoComplete: request.AutoComplete,
		Duration:     request.Duration,
	}
	if err = s.Database.AddJob(&job); err != nil {
		s.l.Logf(logger.ErrorLevel, "Add job to database failed: %s", err)
		return err
	}
	response.JobId = job.JobID

	s.mu.Lock()
	defer s.mu.Unlock()
	record := jobRecord{job: &job}
	s.jobs[job.JobID] = &record
	s.runTranscodingTask(&record)

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
		response.Status = convertStatus(r.Status())
	} else {
		response.Status = rms_transcoder.GetJobResponse_Done
	}

	return nil
}

func (s *Service) CancelJob(ctx context.Context, request *rms_transcoder.CancelJobRequest, empty *emptypb.Empty) error {
	record, err := s.getAndCancelJob(request.JobId)
	if err != nil {
		return err
	}
	if err := s.Database.RemoveJob(request.JobId); err != nil {
		s.l.Logf(logger.ErrorLevel, "Remove job %s from database failed: %s", request.JobId, err)
	}
	if request.RemoveFiles {
		if err = clearContent(s.Config.Directory, record.job.Destination); err != nil {
			s.l.Logf(logger.WarnLevel, "Clear content for %s failed: %s", request.JobId, err)
		}
	}
	return nil
}

func (s *Service) getAndCancelJob(id string) (*jobRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.jobs[id]
	if !ok {
		return nil, errors.New("job not found")
	}

	r := record.receipt
	if r != nil {
		r.Cancel()
	}

	delete(s.jobs, id)

	return record, nil
}

func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
}
