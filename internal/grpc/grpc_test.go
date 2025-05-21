package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	userpbv1 "github.com/vadim8q258475/store-user-microservice/gen/v1"
	"github.com/vadim8q258475/store-user-microservice/internal/grpc/tests/mocks"
	repoGen "github.com/vadim8q258475/store-user-microservice/internal/repo/ent"
	"github.com/vadim8q258475/store-user-microservice/internal/service"
	"go.uber.org/mock/gomock"
)

func TestGrpcService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepo(ctrl)
	service := service.NewService(mockRepo)
	grpcService := NewGrpcService(service)
	ctx := context.Background()
	var id uint32 = 0
	email := "email1"
	password := "password1"

	// create ok
	mockRepo.EXPECT().
		Create(ctx, email, password).
		Return(&repoGen.User{ID: 0, Email: email, Password: password}, nil)

	response, err := grpcService.Create(ctx, &userpbv1.Create_Request{Email: email, Password: password})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, id, response.Id)

	// create err
	mockRepo.EXPECT().
		Create(ctx, "", password).
		Return(nil, assert.AnError)
	response, err = grpcService.Create(ctx, &userpbv1.Create_Request{Email: "", Password: password})
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGrpcService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepo(ctrl)
	service := service.NewService(mockRepo)
	grpcService := NewGrpcService(service)
	ctx := context.Background()
	var id uint32 = 0
	var notValidId uint32 = 404
	email := "email1"
	password := "password1"

	// get ok
	mockRepo.EXPECT().
		GetByID(ctx, id).
		Return(&repoGen.User{ID: int(id), Email: email, Password: password}, nil)

	response, err := grpcService.GetByID(ctx, &userpbv1.GetByID_Request{Id: id})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, id, response.User.Id)
	assert.Equal(t, email, response.User.Email)
	assert.Equal(t, password, response.User.Password)

	// get err
	mockRepo.EXPECT().
		GetByID(ctx, notValidId).
		Return(nil, assert.AnError)
	response, err = grpcService.GetByID(ctx, &userpbv1.GetByID_Request{Id: notValidId})
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGrpcService_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepo(ctrl)
	service := service.NewService(mockRepo)
	grpcService := NewGrpcService(service)
	ctx := context.Background()
	var id uint32 = 0
	email := "email1"
	password := "password1"

	// get ok
	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(&repoGen.User{ID: int(id), Email: email, Password: password}, nil)

	response, err := grpcService.GetByEmail(ctx, &userpbv1.GetByEmail_Request{Email: email})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, id, response.User.Id)
	assert.Equal(t, email, response.User.Email)
	assert.Equal(t, password, response.User.Password)

	// get err
	mockRepo.EXPECT().
		GetByEmail(ctx, "").
		Return(nil, assert.AnError)
	response, err = grpcService.GetByEmail(ctx, &userpbv1.GetByEmail_Request{Email: ""})
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGrpcService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepo(ctrl)
	service := service.NewService(mockRepo)
	grpcService := NewGrpcService(service)
	ctx := context.Background()
	var id uint32 = 0

	repoList := []*repoGen.User{
		{ID: 0, Email: "email1", Password: "pass1"},
		{ID: 1, Email: "email2", Password: "pass2"},
	}

	// list ok
	mockRepo.EXPECT().
		List(ctx).
		Return(repoList, nil)
	response, err := grpcService.List(ctx, &userpbv1.List_Request{})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, id, response.Users[0].Id)
	assert.Equal(t, "email1", response.Users[0].Email)
	assert.Equal(t, "pass1", response.Users[0].Password)

	// list err
	mockRepo.EXPECT().
		List(ctx).
		Return(nil, assert.AnError)
	response, err = grpcService.List(ctx, &userpbv1.List_Request{})

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGrpcService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepo(ctrl)
	service := service.NewService(mockRepo)
	grpcService := NewGrpcService(service)
	ctx := context.Background()
	var id uint32 = 0
	var notValidId uint32 = 404
	email := "email1"
	newPassword := "newPassword1"

	// create ok
	mockRepo.EXPECT().
		Update(ctx, id, email, newPassword).
		Return(&repoGen.User{ID: 0, Email: email, Password: newPassword}, nil)

	response, err := grpcService.Update(ctx, &userpbv1.Update_Request{User: &userpbv1.User{Id: id, Email: email, Password: newPassword}})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, id, response.Id)

	// Update err
	mockRepo.EXPECT().
		Update(ctx, notValidId, email, newPassword).
		Return(nil, assert.AnError)
	response, err = grpcService.Update(ctx, &userpbv1.Update_Request{User: &userpbv1.User{Id: notValidId, Email: email, Password: newPassword}})
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGrpcService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepo(ctrl)
	service := service.NewService(mockRepo)
	grpcService := NewGrpcService(service)
	ctx := context.Background()
	var id uint32 = 404

	// delete ok
	mockRepo.EXPECT().
		Delete(ctx, id).
		Return(nil)

	_, err := grpcService.Delete(ctx, &userpbv1.Delete_Request{Id: id})
	assert.NoError(t, err)

	// delete err
	mockRepo.EXPECT().
		Delete(ctx, id).
		Return(assert.AnError)

	_, err = grpcService.Delete(ctx, &userpbv1.Delete_Request{Id: id})

	assert.Error(t, err)
}
