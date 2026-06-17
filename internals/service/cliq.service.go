package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rivando-al-rasyid/cliq/internals/dto"
)

type CliqRepository interface {
	CreateSlug(ctx context.Context, userID uuid.UUID, originLink string, slug string) error
}

type CliqService struct {
	repo CliqRepository
	rdb  *redis.Client
}

func NewCliqService(repo CliqRepository, rdb *redis.Client) *CliqService {
	return &CliqService{repo: repo, rdb: rdb}
}

func (c *CliqService) CreateSlug(ctx context.Context, userID uuid.UUID, link dto.Link) (string, error) {
	if err := c.repo.CreateSlug(ctx, userID, link.OriginLink, link.Slug); err != nil {
		return "", err
	}

	return link.Slug, nil
}
