package service

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/entity"
	"createtodayapi/internal/logger"
	"createtodayapi/internal/storage"
)

type ProfileService struct {
	repo storage.Users
}

func (s *ProfileService) GetProfile(ctx context.Context, userId int) (*entity.Profile, error) {
	profile, err := s.repo.GetProfileByUserId(ctx, userId)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	return profile, nil
}

func NewProfileService(repo storage.Users) *ProfileService {
	return &ProfileService{
		repo: repo,
	}
}
