package grpc

import (
	"context"

	userV1 "github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/pb/myorg/user/v1"
)

func (s Server) CreateUser(ctx context.Context, request *userV1.CreateUserRequest) (*userV1.CreateUserResponse, error) {
	return s.service.CreateUser(ctx, request)
}

func (s Server) GetUser(ctx context.Context, request *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	return s.service.GetUser(ctx, request)
}

func (s Server) UpdateUser(ctx context.Context, request *userV1.UpdateUserRequest) (*userV1.UpdateUserResponse, error) {
	return s.service.UpdateUser(ctx, request)
}

func (s Server) DeleteUser(ctx context.Context, request *userV1.DeleteUserRequest) (*userV1.DeleteUserResponse, error) {
	return s.service.DeleteUser(ctx, request)
}
