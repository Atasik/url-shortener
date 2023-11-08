package service

import (
	"errors"
	"link-shortener/internal/domain"
	mock_repository "link-shortener/internal/repository/mocks"
	mock_encoder "link-shortener/pkg/encoder/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLinkService_CreateToken(t *testing.T) {
	type mockBehaviour func(rl *mock_repository.MockLinkRepo, e *mock_encoder.MockEncoder, link domain.Link)

	tests := []struct {
		name          string
		link          domain.Link
		mockBehaviour mockBehaviour
		want          string
		wantErr       bool
	}{
		{
			name: "DB OK",
			link: domain.Link{
				OriginalURL: "https://github.com/Atasik",
			},
			mockBehaviour: func(rl *mock_repository.MockLinkRepo, e *mock_encoder.MockEncoder, link domain.Link) {
				rl.EXPECT().AddOriginalURL(link).Return(int64(1), nil)
				e.EXPECT().Encode(int64(1)).Return("token")
			},
			want: "token",
		}, {
			name: "DB Error",
			link: domain.Link{
				OriginalURL: "https://github.com/Atasik",
			},
			mockBehaviour: func(rl *mock_repository.MockLinkRepo, e *mock_encoder.MockEncoder, link domain.Link) {
				rl.EXPECT().AddOriginalURL(link).Return(int64(0), errors.New("something went wrong"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		c := gomock.NewController(t)
		defer c.Finish()

		repoLink := mock_repository.NewMockLinkRepo(c)
		encoder := mock_encoder.NewMockEncoder(c)
		test.mockBehaviour(repoLink, encoder, test.link)

		linkService := NewLinkService(repoLink, encoder)

		got, err := linkService.CreateToken(test.link)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		}
	}
}

func TestLinkService_GetOriginalURL(t *testing.T) {
	type mockBehaviour func(rl *mock_repository.MockLinkRepo, e *mock_encoder.MockEncoder, token string)

	tests := []struct {
		name          string
		token         string
		mockBehaviour mockBehaviour
		want          string
		wantErr       bool
	}{
		{
			name:  "DB OK",
			token: "token",
			mockBehaviour: func(rl *mock_repository.MockLinkRepo, e *mock_encoder.MockEncoder, token string) {
				e.EXPECT().Decode(token).Return(int64(1))
				link := domain.Link{ID: int64(1)}
				rl.EXPECT().GetOrginalURL(link).Return("https://github.com/Atasik", nil)
			},
			want: "https://github.com/Atasik",
		}, {
			name:  "DB Error",
			token: "token",
			mockBehaviour: func(rl *mock_repository.MockLinkRepo, e *mock_encoder.MockEncoder, token string) {
				e.EXPECT().Decode(token).Return(int64(1))
				link := domain.Link{ID: int64(1)}
				rl.EXPECT().GetOrginalURL(link).Return("", errors.New("something went wrong"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		c := gomock.NewController(t)
		defer c.Finish()

		repoLink := mock_repository.NewMockLinkRepo(c)
		encoder := mock_encoder.NewMockEncoder(c)
		test.mockBehaviour(repoLink, encoder, test.token)

		linkService := NewLinkService(repoLink, encoder)

		got, err := linkService.GetOriginalURL(test.token)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		}
	}
}
