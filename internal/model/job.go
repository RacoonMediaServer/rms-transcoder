package model

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	"gorm.io/datatypes"
)

type JobResult int

const (
	JobNotComplete JobResult = iota
	JobDone
	JobFailed
	JobCancelled
)

type Job struct {
	JobID        string `gorm:"primaryKey"`
	Transcoding  datatypes.JSONType[media.TranscodingSettings]
	Source       string
	Destination  string `gorm:"unique"`
	AutoComplete bool
	Done         bool
}
