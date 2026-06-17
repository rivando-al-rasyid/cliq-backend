package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"

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
	slug := link.Slug

	if slug == "" {
		generatedSlug, err := generateSlug(8)
		if err != nil {
			return "", err
		}

		slug = generatedSlug
	}

	if err := c.repo.CreateSlug(ctx, userID, link.OriginLink, slug); err != nil {
		return "", err
	}

	return slug, nil
}

func generateSlug(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	slug := base64.RawURLEncoding.EncodeToString(b)
	slug = strings.ReplaceAll(slug, "_", "")
	slug = strings.ReplaceAll(slug, "-", "")

	if len(slug) > length {
		slug = slug[:length]
	}

	return slug, nil
}
