package profiles

import rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"

var predefinedProfiles = map[string]*rms_transcoder.Profile{
	"telegram": {
		Id:          "telegram",
		VideoMute:   false,
		VideoCodec:  makePtr("h264"),
		VideoWidth:  makePtr[uint32](480),
		VideoHeight: makePtr[uint32](320),
		AudioMute:   true,
		IsReserved:  true,
		IsReadOnly:  true,
	},
	"default": {
		Id:         "default",
		VideoCodec: nil,
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
