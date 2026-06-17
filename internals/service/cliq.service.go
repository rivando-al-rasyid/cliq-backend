package service

import (
	"github.com/redis/go-redis/v9"
)

type CliqRepository interface {
}

type CliqService struct {
	repo CliqRepository
	rdb  *redis.Client
}

func NewCliqService(repo CliqRepository, rdb *redis.Client) *CliqService {
	return &CliqService{repo: repo, rdb: rdb}
}
