package profiles

import rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"

type Database interface {
	LoadProfiles() ([]*rms_transcoder.Profile, error)
	AddProfile(profile *rms_transcoder.Profile) error
	GetProfile(id string) (*rms_transcoder.Profile, error)
	UpdateProfile(profile *rms_transcoder.Profile) error
	RemoveProfile(id string) error
}
