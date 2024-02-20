package profiles

import (
	"context"
	"errors"
	"fmt"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	Database Database
}

func (s Service) Initialize() error {
	for _, profile := range predefinedProfiles {
		if err := s.addReservedProfile(profile); err != nil {
			return fmt.Errorf("add default profile failed: %w", err)
		}
	}
	return nil
}

func (s Service) addReservedProfile(profile *rms_transcoder.Profile) error {
	exists, err := s.Database.GetProfile(profile.Id)
	if err != nil {
		return err
	}
	if exists != nil {
		return nil
	}
	return s.Database.AddProfile(profile)
}

func (s Service) UpdateProfile(ctx context.Context, profile *rms_transcoder.Profile, empty *emptypb.Empty) error {
	if isProfileReadOnly(profile.Id) {
		return errors.New("profile is read only")
	}
	profile.IsReserved = isProfileReserved(profile.Id)
	profile.IsReadOnly = false
	return s.Database.UpdateProfile(profile)
}

func (s Service) AddProfile(ctx context.Context, profile *rms_transcoder.Profile, empty *emptypb.Empty) error {
	profile.IsReserved = false
	profile.IsReadOnly = false
	return s.Database.AddProfile(profile)
}

func (s Service) GetProfiles(ctx context.Context, empty *emptypb.Empty, response *rms_transcoder.GetProfilesResponse) error {
	profiles, err := s.Database.LoadProfiles()
	if err != nil {
		return err
	}
	response.Profiles = profiles
	return nil
}

func (s Service) RemoveProfile(ctx context.Context, request *rms_transcoder.RemoveProfileRequest, empty *emptypb.Empty) error {
	if isProfileReserved(request.Id) {
		return errors.New("cannot delete reserved profile")
	}
	return s.Database.RemoveProfile(request.Id)
}
