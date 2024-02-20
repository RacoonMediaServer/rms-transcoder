package transcoder

import (
	"context"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	Profiles ProfileService
	Database Database
}

func (s Service) Initialize() error {
	return nil
}

func (s Service) AddJob(ctx context.Context, request *rms_transcoder.AddJobRequest, response *rms_transcoder.AddJobResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetJob(ctx context.Context, request *rms_transcoder.GetJobRequest, response *rms_transcoder.GetJobResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) CancelJob(ctx context.Context, request *rms_transcoder.CancelJobRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}
