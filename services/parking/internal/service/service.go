package service

import (
	"context"

	"parking-service/internal/cache"
	"parking-service/internal/repository"
)

type Service struct {
	repo  *repository.PostgresRepo
	cache *cache.RedisCache
}

func New(repo *repository.PostgresRepo, cache *cache.RedisCache) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) AssignSlot(ctx context.Context) (int64, error) {
	id, err := s.repo.AssignSlot(ctx)
	if err != nil {
		return 0, err
	}
	// invalidate cache
	s.cache.InvalidateAvailable(ctx)
	return id, nil
}

func (s *Service) ReleaseSlot(ctx context.Context, id int64) error {
	if err := s.repo.ReleaseSlot(ctx, id); err != nil {
		return err
	}
	s.cache.InvalidateAvailable(ctx)
	return nil
}

func (s *Service) GetAvailable(ctx context.Context) (int32, error) {
	if v, err := s.cache.GetAvailable(ctx); err == nil {
		return v, nil
	}
	v, err := s.repo.CountFree(ctx)
	if err == nil {
		s.cache.SetAvailable(ctx, v)
	}
	return v, err
}
