package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vadim8q258475/store-user-microservice/app"
	"github.com/vadim8q258475/store-user-microservice/config"
	"github.com/vadim8q258475/store-user-microservice/internal/cacher"
	grpcService "github.com/vadim8q258475/store-user-microservice/internal/grpc"
	"github.com/vadim8q258475/store-user-microservice/internal/interceptor"
	"github.com/vadim8q258475/store-user-microservice/internal/repo"
	"github.com/vadim8q258475/store-user-microservice/internal/service"
	"github.com/vadim8q258475/store-user-microservice/queue"
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

	// queue
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	))
	fmt.Println(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	))
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	_, err = ch.QueueDeclare(
		cfg.RabbitMQQueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)

	publisher := queue.NewRabbitMQQueue(ch, cfg)

	// service
	service := service.NewService(proxy, publisher)

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
