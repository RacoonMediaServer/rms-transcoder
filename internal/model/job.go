package model

type Job struct {
	ID           string `gorm:"primaryKey"`
	Profile      string
	Source       string
	Destination  string `gorm:"unique"`
	AutoComplete bool
	Done         bool
}
