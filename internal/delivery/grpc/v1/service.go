package v1

import (
	"link-shortener/internal/delivery/grpc"
	"link-shortener/internal/service"
)

type GrpcService struct {
	grpc.UnimplementedLinkServer
	services *service.Service
}

func NewGrpcService(services *service.Service) *GrpcService {
	return &GrpcService{
		services: services,
	}
}
