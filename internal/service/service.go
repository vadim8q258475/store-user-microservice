package service

import (
	"context"

	gen "github.com/vadim8q258475/store-user-microservice/gen/v1"
	"github.com/vadim8q258475/store-user-microservice/iternal/repo"
	"github.com/vadim8q258475/store-user-microservice/iternal/repo/ent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if ent.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return user, err
}

func (s *service) GetByID(ctx context.Context, req *gen.GetByID_Request) (*ent.User, error) {
	user, err := s.repo.GetByID(ctx, req.Id)
	if ent.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return user, err
}

func (s *service) Update(ctx context.Context, req *gen.Update_Request) (*ent.User, error) {
	return s.repo.Update(ctx, req.User.Id, req.User.Email, req.User.Password)
}

func (s *service) Delete(ctx context.Context, req *gen.Delete_Request) error {
	return s.repo.Delete(ctx, req.Id)
}
