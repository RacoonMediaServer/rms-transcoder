package profiles

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
)

var predefinedProfiles = map[string]*rms_transcoder.Profile{
	"telegram": {
		Id: "telegram",
		Settings: &media.TranscodingSettings{
			Video: &media.VideoTranscodingSettings{
				Codec:  makePtr("h264"),
				Width:  makePtr[uint32](480),
				Height: makePtr[uint32](320),
			},
		},
		IsReserved: true,
		IsReadOnly: true,
	},
	"default": {
		Id:         "default",
		Settings:   &media.TranscodingSettings{},
		IsReserved: true,
	},
}

func isProfileReserved(id string) bool {
	p, ok := predefinedProfiles[id]
	if !ok {
		return false
	}
	return p.IsReserved
}

func isProfileReadOnly(id string) bool {
	p, ok := predefinedProfiles[id]
	if !ok {
		return false
	}
	return p.IsReadOnly
}

func makePtr[T any](value T) *T {
	return &value
}
