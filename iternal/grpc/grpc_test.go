package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	userpbv1 "github.com/vadim8q258475/store-user-microservice/gen/v1"
	"github.com/vadim8q258475/store-user-microservice/iternal/grpc/tests/mocks"
	repoGen "github.com/vadim8q258475/store-user-microservice/iternal/repo/ent"
	"github.com/vadim8q258475/store-user-microservice/iternal/service"
	"go.uber.org/mock/gomock"
)

func TestGrpcService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepo(ctrl)
	service := service.NewService(mockRepo)
	grpcService := NewGrpcService(service)
	ctx := context.Background()
	email := "email1"
	password := "password1"

	// create ok
	mockRepo.EXPECT().
		Create(ctx, email, password).
		Return(&repoGen.User{ID: 0, Email: email, Password: password}, nil)

	response, err := grpcService.Create(ctx, &userpbv1.Create_Request{Email: email, Password: password})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "0", response.Id)

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
	id := 0
	strID := "0"
	email := "email1"
	password := "password1"

	// get ok
	mockRepo.EXPECT().
		GetByID(ctx, id).
		Return(&repoGen.User{ID: id, Email: email, Password: password}, nil)

	response, err := grpcService.GetByID(ctx, &userpbv1.GetByID_Request{Id: strID})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, strID, response.Id)
	assert.Equal(t, email, response.Email)
	assert.Equal(t, password, response.Password)

	// get err
	mockRepo.EXPECT().
		GetByID(ctx, -1).
		Return(nil, assert.AnError)
	response, err = grpcService.GetByID(ctx, &userpbv1.GetByID_Request{Id: "-1"})
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
	id := 0
	strID := "0"
	email := "email1"
	password := "password1"

	// get ok
	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(&repoGen.User{ID: id, Email: email, Password: password}, nil)

	response, err := grpcService.GetByEmail(ctx, &userpbv1.GetByEmail_Request{Email: email})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, strID, response.Id)
	assert.Equal(t, email, response.Email)
	assert.Equal(t, password, response.Password)

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
	assert.Equal(t, "0", response.Users[0].Id)
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
	id := 0
	strID := "0"
	email := "email1"
	password := "password1"
	newPassword := "newPassword1"

	// create ok
	mockRepo.EXPECT().
		Update(ctx, id, email, newPassword).
		Return(&repoGen.User{ID: 0, Email: email, Password: newPassword}, nil)

	response, err := grpcService.Update(ctx, &userpbv1.Update_Request{Id: strID, Email: email, Password: newPassword})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "0", response.Id)

	// Update err
	mockRepo.EXPECT().
		Update(ctx, -1, email, password).
		Return(nil, assert.AnError)
	response, err = grpcService.Update(ctx, &userpbv1.Update_Request{Id: "-1", Email: email, Password: password})
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
	id := 0
	email := "email1"
	password := "password1"

	// delete ok
	mockRepo.EXPECT().
		Delete(ctx, id).
		Return(nil)

	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(&repoGen.User{ID: id, Email: email, Password: password}, nil)

	response, err := grpcService.Delete(ctx, &userpbv1.Delete_Request{Email: email})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int32(200), response.Status)

	// delete err
	mockRepo.EXPECT().
		Delete(ctx, id).
		Return(assert.AnError)

	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(&repoGen.User{ID: id, Email: email, Password: password}, nil)

	response, err = grpcService.Delete(ctx, &userpbv1.Delete_Request{Email: email})

	assert.Error(t, err)
	assert.Nil(t, response)
}
