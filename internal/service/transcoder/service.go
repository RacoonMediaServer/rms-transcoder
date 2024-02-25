package transcoder

import (
	"context"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/RacoonMediaServer/rms-transcoder/internal/model"
	"github.com/google/uuid"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	l        logger.Logger
	Profiles ProfileService
	Database Database
	Workers  Workers
}

func (s *Service) Initialize() error {
	s.l = logger.DefaultLogger.Fields(map[string]interface{}{"from": "transcoder"})
	return nil
}

func (s *Service) AddJob(ctx context.Context, request *rms_transcoder.AddJobRequest, response *rms_transcoder.AddJobResponse) error {
	id, err := uuid.NewUUID()
	if err != nil {
		s.l.Logf(logger.ErrorLevel, "Generate id failed: %s", err)
		return err
	}
	job := model.Job{
		ID:           id.String(),
		Profile:      request.Profile,
		Source:       request.Source,
		Destination:  request.Destination,
		AutoComplete: request.AutoComplete,
		Done:         false,
	}
	// TODO: search profile
	if err = s.Database.AddJob(&job); err != nil {
		s.l.Logf(logger.ErrorLevel, "Add job to database failed: %s", err)
		return err
	}
	return nil
}

func (s *Service) GetJob(ctx context.Context, request *rms_transcoder.GetJobRequest, response *rms_transcoder.GetJobResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) CancelJob(ctx context.Context, request *rms_transcoder.CancelJobRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}
