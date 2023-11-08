package app

import (
	"context"
	"fmt"
	"link-shortener/internal/config"
	pb "link-shortener/internal/delivery/grpc"
	v1 "link-shortener/internal/delivery/grpc/v1"
	delivery "link-shortener/internal/delivery/http"
	"link-shortener/internal/repository"
	"link-shortener/internal/server"
	"link-shortener/internal/service"
	"link-shortener/pkg/database/postgres"
	"link-shortener/pkg/encoder"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

const (
	timeout     = 5 * time.Second
	encodingLen = 10
)

// @title Ozon Backend Test
// @version 2.0
// @description Тех. Задание Ozon

// @host localhost:8080
// @BasePath /
func Run(configsDir string) {
	var (
		flags Flags
		repos *repository.Repository
		db    *sqlx.DB
	)

	flags.getFlags()
	cfg, err := config.InitConfig(configsDir)
	if err != nil {
		log.Fatal("Error occurred while loading config: ", err.Error())
	}

	if flags.Postgres {
		db, err = postgres.NewPostgresqlDB(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username,
			cfg.Postgres.DBName, cfg.Postgres.Password, cfg.Postgres.SSLMode)
		if err != nil {
			log.Fatalf("Error occurred while loading DB: %s\n", err.Error())
			return
		}
		repos = repository.NewRepository(
			repository.NewLinkPostgresqlRepo(db),
		)
	} else {
		repos = repository.NewRepository(
			repository.NewLinkMemoryRepo(),
		)
	}

	b63Encoder := encoder.NewBase63Encoder(encodingLen)

	services := service.NewService(
		service.NewLinkService(repos, b63Encoder),
	)

	h := delivery.NewHandler(services)
	mux := h.Init()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		log.Fatalln("cant listen to port", err)
	}
	defer lis.Close()

	linkService := v1.NewGrpcService(services)
	grpcServer := grpc.NewServer()
	pb.RegisterLinkServer(grpcServer, linkService)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	srv := server.NewServer(cfg, mux)
	go func() {
		if err := grpcServer.Serve(lis); err != net.ErrClosed {
			log.Fatalln("error happened: ", err.Error())
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Fatalf("panic occurred: %s\n", err)
			}
		}()
		if err := srv.Run(); err != http.ErrServerClosed {
			log.Fatalln("error happened: ", err.Error())
		}
	}()

	log.Println("Application is running")

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	log.Println("Application is shutting down")

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error occurred on http server shutting down: %s", err.Error())
	}

	if flags.Postgres {
		if err := db.Close(); err != nil {
			log.Printf("error occurred on db connection close: %s", err.Error())
		}
	}
}
