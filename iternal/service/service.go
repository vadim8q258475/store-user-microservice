package service

import (
	"context"
	"strconv"

	gen "github.com/vadim8q258475/store-user-microservice/gen/v1"
	"github.com/vadim8q258475/store-user-microservice/iternal/repo"
	"github.com/vadim8q258475/store-user-microservice/iternal/repo/ent"
)

type Service interface {
	Create(ctx context.Context, req *gen.Create_Request) (*ent.User, error)
	List(ctx context.Context) ([]*ent.User, error)
	GetByEmail(ctx context.Context, req *gen.GetByEmail_Request) (*ent.User, error)
	GetByID(ctx context.Context, req *gen.GetByID_Request) (*ent.User, error)
	Update(ctx context.Context, req *gen.Update_Request) (*ent.User, error)
	Delete(ctx context.Context, req *gen.Delete_Request) error
}

type service struct {
	repo repo.Repo
}

func NewService(repo repo.Repo) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *gen.Create_Request) (*ent.User, error) {
	return s.repo.Create(ctx, req.Email, req.Password)
}

func (s *service) List(ctx context.Context) ([]*ent.User, error) {
	return s.repo.List(ctx)
}

func (s *service) GetByEmail(ctx context.Context, req *gen.GetByEmail_Request) (*ent.User, error) {
	return s.repo.GetByEmail(ctx, req.Email)
}

func (s *service) GetByID(ctx context.Context, req *gen.GetByID_Request) (*ent.User, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, req *gen.Update_Request) (*ent.User, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, req.Email, req.Password)
}

func (s *service) Delete(ctx context.Context, req *gen.Delete_Request) error {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, user.ID)
}
