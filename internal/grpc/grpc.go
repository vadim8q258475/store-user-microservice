package grpc

import (
	"context"

	gen "github.com/vadim8q258475/store-user-microservice/gen/v1"
	"github.com/vadim8q258475/store-user-microservice/internal/service"
)

type GrpcService struct {
	gen.UnimplementedUserServiceServer
	service service.Service
}

func NewGrpcService(service service.Service) *GrpcService {
	return &GrpcService{
		service: service,
	}
}

func (g *GrpcService) Create(ctx context.Context, req *gen.Create_Request) (*gen.Create_Response, error) {
	user, err := g.service.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &gen.Create_Response{Id: uint32(user.ID)}, nil
}

func (g *GrpcService) Delete(ctx context.Context, req *gen.Delete_Request) (*gen.Delete_Response, error) {
	err := g.service.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &gen.Delete_Response{}, nil
}

func (g *GrpcService) GetByEmail(ctx context.Context, req *gen.GetByEmail_Request) (*gen.GetByEmail_Response, error) {
	user, err := g.service.GetByEmail(ctx, req)
	if err != nil {
		return nil, err
	}
	return &gen.GetByEmail_Response{User: &gen.User{Id: uint32(user.ID), Email: user.Email, Password: user.Password}}, nil
}

func (g *GrpcService) GetByID(ctx context.Context, req *gen.GetByID_Request) (*gen.GetByID_Response, error) {
	user, err := g.service.GetByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return &gen.GetByID_Response{User: &gen.User{Id: uint32(user.ID), Email: user.Email, Password: user.Password}}, nil
}

func (g *GrpcService) List(ctx context.Context, req *gen.List_Request) (*gen.List_Response, error) {
	users, err := g.service.List(ctx)
	if err != nil {
		return nil, err
	}
	resp := &gen.List_Response{Users: make([]*gen.User, len(users))}
	for i, user := range users {
		resp.Users[i] = &gen.User{Id: uint32(user.ID), Email: user.Email, Password: user.Password}
	}
	return resp, nil
}

func (g *GrpcService) Update(ctx context.Context, req *gen.Update_Request) (*gen.Update_Response, error) {
	user, err := g.service.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	return &gen.Update_Response{Id: uint32(user.ID)}, nil
}
