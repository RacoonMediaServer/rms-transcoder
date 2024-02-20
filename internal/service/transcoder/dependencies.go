package transcoder

import (
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/RacoonMediaServer/rms-transcoder/internal/model"
)

type ProfileService interface {
	GetProfile(id string) (*rms_transcoder.Profile, error)
}

type Database interface {
	LoadJobs() ([]*model.Job, error)
	AddJob(profile *model.Job) error
	UpdateJob(profile *model.Job) error
	RemoveJob(id string) error
}
