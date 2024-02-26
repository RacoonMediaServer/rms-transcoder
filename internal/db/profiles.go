package db

import (
	"errors"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"gorm.io/gorm"
)

func (d *Database) LoadProfiles() ([]*rms_transcoder.Profile, error) {
	var result []*rms_transcoder.Profile
	if err := d.conn.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (d *Database) AddProfile(profile *rms_transcoder.Profile) error {
	return d.conn.Create(profile).Error
}

func (d *Database) GetProfile(id string) (*rms_transcoder.Profile, error) {
	var profile rms_transcoder.Profile
	if err := d.conn.First(&profile, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (d *Database) UpdateProfile(profile *rms_transcoder.Profile) error {
	return d.conn.Save(profile).Error
}

func (d *Database) RemoveProfile(id string) error {
	return d.conn.Model(&rms_transcoder.Profile{}).Unscoped().Delete(&rms_transcoder.Profile{Id: id}).Error
}
