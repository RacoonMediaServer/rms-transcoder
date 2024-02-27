package db

import (
	"github.com/RacoonMediaServer/rms-transcoder/internal/model"
)

func (d *Database) LoadJobs() ([]*model.Job, error) {
	var result []*model.Job
	if err := d.conn.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (d *Database) AddJob(profile *model.Job) error {
	return d.conn.Create(profile).Error
}

func (d *Database) UpdateJob(profile *model.Job) error {
	return d.conn.Save(profile).Error
}

func (d *Database) RemoveJob(id string) error {
	return d.conn.Model(&model.Job{}).Unscoped().Delete(&model.Job{JobID: id}).Error
}
