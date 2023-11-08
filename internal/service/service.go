package service

import "link-shortener/internal/domain"

type Service struct {
	Link
}

type Link interface {
	CreateToken(link domain.Link) (string, error)
	GetOriginalURL(token string) (string, error)
}

func NewService(linkService Link) *Service {
	return &Service{Link: linkService}
}
