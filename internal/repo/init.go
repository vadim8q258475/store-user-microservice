package repo

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/vadim8q258475/store-user-microservice/config"
	"github.com/vadim8q258475/store-user-microservice/internal/repo/ent"
)

func InitDB(cfg config.Config) (*ent.Client, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	client, err := ent.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}
	return client, err
}
