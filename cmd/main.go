package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/vadim8q258475/store-user-microservice/app"
	"github.com/vadim8q258475/store-user-microservice/config"
	"github.com/vadim8q258475/store-user-microservice/iternal/cacher"
	grpcService "github.com/vadim8q258475/store-user-microservice/iternal/grpc"
	"github.com/vadim8q258475/store-user-microservice/iternal/interceptor"
	"github.com/vadim8q258475/store-user-microservice/iternal/repo"
	"github.com/vadim8q258475/store-user-microservice/iternal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// TODO
// add cacher

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
	repos := repo.NewRepo(client)

	// init cache
	cacheClient := redis.NewClient(&redis.Options{
		Addr: cfg.CacheHost + ":" + cfg.CachePort,
	})
	cacher := cacher.NewCacher(cacheClient, cfg)

	// repo proxy
	proxy := repo.NewProxy(repos, cacher)

	// service
	service := service.NewService(proxy)

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
