package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CliqRepo struct {
	db *pgxpool.Pool
}

func NewCliqRepo(db *pgxpool.Pool) *CliqRepo {
	return &CliqRepo{db: db}
}

func (c *CliqRepo) CreateSlug(
	ctx context.Context,
	userID uuid.UUID,
	originLink string,
	slug string,
) error {
	_, err := c.db.Exec(ctx,
		`
		INSERT INTO links (
			user_id,
			origin_link,
			slug
		)
		VALUES ($1, $2, $3)
		`,
		userID,
		originLink,
		slug,
	)
	if err != nil {
		return fmt.Errorf("create slug: %w", err)
	}

	return nil
}
