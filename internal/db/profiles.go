package db

import (
	"encoding/json"
	"errors"
	"fmt"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"gorm.io/gorm"
)

type profileRecord struct {
	ID      string `gorm:"primaryKey"`
	Profile string
}

func encodeProfile(profile *rms_transcoder.Profile) (string, error) {
	content, err := json.Marshal(profile)
	return string(content), err
}

func decodeProfile(encoded string) (*rms_transcoder.Profile, error) {
	profile := rms_transcoder.Profile{}
	err := json.Unmarshal([]byte(encoded), &profile)
	return &profile, err
}

func (d *Database) LoadProfiles() ([]*rms_transcoder.Profile, error) {
	var records []*profileRecord
	if err := d.conn.Find(&records).Error; err != nil {
		return nil, err
	}
	result := make([]*rms_transcoder.Profile, 0, len(records))
	for _, record := range records {
		profile, err := decodeProfile(record.Profile)
		if err != nil {
			return nil, fmt.Errorf("decode profile error: %w", err)
		}
		result = append(result, profile)
	}
	return result, nil
}

func (d *Database) AddProfile(profile *rms_transcoder.Profile) error {
	encoded, err := encodeProfile(profile)
	if err != nil {
		return fmt.Errorf("encode profile failed: %w", err)
	}
	record := profileRecord{
		ID:      profile.Id,
		Profile: encoded,
	}
	return d.conn.Create(&record).Error
}

func (d *Database) GetProfile(id string) (*rms_transcoder.Profile, error) {
	var record profileRecord
	if err := d.conn.First(&record, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	profile, err := decodeProfile(record.Profile)
	return profile, err
}

func (d *Database) UpdateProfile(profile *rms_transcoder.Profile) error {
	encoded, err := encodeProfile(profile)
	if err != nil {
		return fmt.Errorf("encode profile failed: %w", err)
	}
	record := profileRecord{
		ID:      profile.Id,
		Profile: encoded,
	}
	return d.conn.Save(&record).Error
}

func (d *Database) RemoveProfile(id string) error {
	return d.conn.Model(&profileRecord{}).Unscoped().Delete(&profileRecord{ID: id}).Error
}
