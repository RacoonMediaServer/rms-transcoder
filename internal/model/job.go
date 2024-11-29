package model

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	"gorm.io/datatypes"
)

type Job struct {
	JobID        string `gorm:"primaryKey"`
	Transcoding  datatypes.JSONType[media.TranscodingSettings]
	Source       string
	Destination  string `gorm:"unique"`
	AutoComplete bool
	Done         bool
	Duration     *uint32
	Offset       *uint32
}
