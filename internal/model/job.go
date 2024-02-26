package model

type JobResult int

const (
	JobNotComplete JobResult = iota
	JobDone
	JobFailed
)

type Job struct {
	ID           string `gorm:"primaryKey"`
	Profile      string
	Source       string
	Destination  string `gorm:"unique"`
	AutoComplete bool
	Result       JobResult
}
