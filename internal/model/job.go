package model

import rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"

type JobResult int

const (
	JobNotComplete JobResult = iota
	JobDone
	JobFailed
	JobCancelled
)

type Job struct {
	JobID        string                  `gorm:"primaryKey"`
	Profile      *rms_transcoder.Profile `gorm:"embedded"`
	Source       string
	Destination  string `gorm:"unique"`
	AutoComplete bool
	Done         bool
}
