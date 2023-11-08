package repository

import (
	"database/sql"
	"link-shortener/internal/domain"

	"github.com/lib/pq"
)

const (
	linksTable = "links"
)

type Repository struct {
	LinkRepo
}

type LinkRepo interface {
	GetOrginalURL(link domain.Link) (string, error)
	AddOriginalURL(link domain.Link) (int64, error)
}

func NewRepository(linkRepo LinkRepo) *Repository {
	return &Repository{LinkRepo: linkRepo}
}

func ParsePostgresError(err error) error {
	if err == nil {
		return nil
	}

	pgErr, ok := err.(*pq.Error)
	if ok {
		if pgErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}
	}

	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	return err
}
