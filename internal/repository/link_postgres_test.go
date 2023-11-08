package repository

import (
	"fmt"
	"link-shortener/internal/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestLinkPostgres_AddOriginalURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewLinkPostgresqlRepo(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   domain.Link
		want    int64
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", linksTable)).
					WithArgs("https://github.com/Atasik").WillReturnRows(rows)
			},
			input: domain.Link{
				OriginalURL: "https://github.com/Atasik",
			},
			want: 1,
		},
		{
			name: "Empty Input",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", linksTable)).
					WithArgs("").WillReturnRows(rows)
			},
			input: domain.Link{
				OriginalURL: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.AddOriginalURL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLinkPostgres_GetOriginalURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewLinkPostgresqlRepo(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   domain.Link
		want    string
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"original_url"}).AddRow("https://github.com/Atasik")
				mock.ExpectQuery(fmt.Sprintf("SELECT original_url FROM %s WHERE (.+)", linksTable)).
					WithArgs(1).WillReturnRows(rows)
			},
			input: domain.Link{
				ID: 1,
			},
			want: "https://github.com/Atasik",
		},
		{
			name: "Empty Input",
			mock: func() {
				rows := sqlmock.NewRows([]string{"original_url"})
				mock.ExpectQuery(fmt.Sprintf("SELECT original_url FROM %s WHERE (.+)", linksTable)).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetOrginalURL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
