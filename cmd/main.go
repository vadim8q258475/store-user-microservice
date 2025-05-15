package main

import (
	"fmt"

	"github.com/vadim8q258475/store-user-microservice/app"
	"github.com/vadim8q258475/store-user-microservice/config"
	grpcService "github.com/vadim8q258475/store-user-microservice/iternal/grpc"
	"github.com/vadim8q258475/store-user-microservice/iternal/interceptor"
	"github.com/vadim8q258475/store-user-microservice/iternal/repo"
	"github.com/vadim8q258475/store-user-microservice/iternal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// interceptor
	intterceptor := interceptor.NewInterceptor(logger)

	// load config
	cfg := config.MustLoadConfig()
	fmt.Println(cfg.String())

	// init db
	client, err := repo.InitDB(cfg)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// repo
	repo := repo.NewRepo(client)

	// service
	service := service.NewService(repo)

	// grpc service
	grpcService := grpcService.NewGrpcService(service)

	// grpc server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			intterceptor.UnaryServerInterceptor,
		),
	)

	// app
	app := app.NewApp(grpcService, server, logger, cfg)

	if err = app.Run(); err != nil {
		panic(err)
	}
}
