package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type CliqRepo struct {
	db *pgxpool.Pool
}

func NewCliqRepo(db *pgxpool.Pool) *CliqRepo {
	return &CliqRepo{db: db}
}
