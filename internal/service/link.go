package service

import (
	"link-shortener/internal/domain"
	"link-shortener/internal/repository"
	"link-shortener/pkg/encoder"
)

type LinkService struct {
	linkRepo repository.LinkRepo
	encoder  encoder.Encoder
}

func NewLinkService(linkRepo repository.LinkRepo, encoder encoder.Encoder) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
		encoder:  encoder,
	}
}

func (serv *LinkService) CreateToken(link domain.Link) (string, error) {
	id, err := serv.linkRepo.AddOriginalURL(link)
	if err != nil {
		return "", err
	}
	token := serv.encoder.Encode(id)
	return token, nil
}

func (serv *LinkService) GetOriginalURL(token string) (string, error) {
	link := domain.Link{
		ID: serv.encoder.Decode(token),
	}
	return serv.linkRepo.GetOrginalURL(link)
}
