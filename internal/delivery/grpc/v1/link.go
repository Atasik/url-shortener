package v1

import (
	"context"
	"link-shortener/internal/delivery/grpc"
	gen "link-shortener/internal/delivery/grpc"
	"link-shortener/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// для grpc мне стало лень уже тесты писать
func (srv *GrpcService) CreateToken(ctx context.Context, in *gen.OriginalURL) (*gen.Token, error) {
	inp := domain.CreateTokenRequest{OriginalURL: in.Url}
	if err := inp.ValidateURL(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid URL")
	}

	link := domain.Link{
		OriginalURL: inp.OriginalURL,
	}

	token, err := srv.services.CreateToken(link)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}
	return &grpc.Token{Token: token}, nil
}

func (srv *GrpcService) GetOriginalURL(ctx context.Context, in *gen.Token) (*gen.OriginalURL, error) {
	inp := domain.GetOriginalURLRequest{Token: in.Token}
	if err := inp.ValidateToken(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid Token")
	}

	longURL, err := srv.services.Link.GetOriginalURL(inp.Token)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}
	return &grpc.OriginalURL{Url: longURL}, err
}
