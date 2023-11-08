package repository

import (
	"fmt"
	"link-shortener/internal/domain"

	"github.com/jmoiron/sqlx"
)

type LinkPostgresqlRepo struct {
	db *sqlx.DB
}

func NewLinkPostgresqlRepo(db *sqlx.DB) *LinkPostgresqlRepo {
	return &LinkPostgresqlRepo{db: db}
}

func (repo *LinkPostgresqlRepo) AddOriginalURL(link domain.Link) (int64, error) {
	var id int64
	query := fmt.Sprintf(`INSERT INTO %s (original_url) VALUES ($1)
						  RETURNING id
						`, linksTable)

	err := repo.db.QueryRow(query, link.OriginalURL).Scan(&id)

	if err == nil {
		return id, nil
	}

	if err = ParsePostgresError(err); err == domain.ErrAlreadyExists {
		query = fmt.Sprintf(`SELECT id FROM %s 
							 WHERE original_url = $1
							`, linksTable)
		if err = repo.db.QueryRow(query, link.OriginalURL).Scan(&id); err != nil {
			return 0, ParsePostgresError(err)
		}
		return id, nil
	}

	return 0, ParsePostgresError(err)
}

func (repo *LinkPostgresqlRepo) GetOrginalURL(link domain.Link) (string, error) {
	var url string
	query := fmt.Sprintf(`SELECT original_url 
						  FROM %s WHERE id = $1
						`, linksTable)

	if err := repo.db.Get(&url, query, link.ID); err != nil {
		return "", ParsePostgresError(err)
	}
	return url, nil
}
