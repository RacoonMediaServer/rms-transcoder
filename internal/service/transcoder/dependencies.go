package transcoder

import (
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/RacoonMediaServer/rms-transcoder/internal/model"
	"github.com/RacoonMediaServer/rms-transcoder/internal/worker"
)

type ProfileService interface {
	GetProfile(id string) (*rms_transcoder.Profile, error)
}

type Database interface {
	LoadJobs() ([]*model.Job, error)
	AddJob(job *model.Job) error
	UpdateJob(job *model.Job) error
	RemoveJob(id string) error
}

type Workers interface {
	Do(t worker.Task) worker.Receipt
}
