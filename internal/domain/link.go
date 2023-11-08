package domain

import (
	"fmt"
	"net/url"
	"regexp"
)

const (
	tokenRegex = "^[a-zA-Z0-9_]{10}$"
)

type Link struct {
	ID          int64  `db:"id"`
	OriginalURL string `db:"original_url"`
}

type CreateTokenRequest struct {
	OriginalURL string `json:"original_url"`
}

type GetOriginalURLRequest struct {
	Token string `json:"token"`
}

func (r CreateTokenRequest) ValidateURL() error {
	if r.OriginalURL == "" {
		return fmt.Errorf("empty URL")
	}

	if _, err := url.ParseRequestURI(r.OriginalURL); err != nil {
		return fmt.Errorf("invalid URL")
	}

	u, err := url.Parse(r.OriginalURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid URL")
	}

	return nil
}

func (r GetOriginalURLRequest) ValidateToken() error {
	if r.Token == "" {
		return fmt.Errorf("empty token")
	}

	valid, err := regexp.Match(tokenRegex, []byte(r.Token))
	if err != nil {
		return fmt.Errorf("regex error")
	}
	if !valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
