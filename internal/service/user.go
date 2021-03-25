package service

import (
	"context"

	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/db"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/obs"
	userV1 "github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/pb/myorg/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) DropAllData(ctx context.Context) error {
	return s.dbClient.DropAllData(ctx)
}

func (s Service) CreateUser(ctx context.Context, request *userV1.CreateUserRequest) (*userV1.CreateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (s Service) GetUser(ctx context.Context, request *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	ctx, span := obs.NewSpan(ctx, "GetUser")
	defer span.End()

	var userPB *userV1.User

	// Perform database query
	err := s.dbClient.RunInReadOnlyTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		if user, err := tx.GetUser(ctx, request.Id); err != nil {
			return err
		} else {
			userPB = user
		}

		return nil
	})

	// Handle error
	if err != nil {
		return nil, err
	}

	return &userV1.GetUserResponse{
		User: userPB,
	}, nil
}

func (s Service) UpdateUser(ctx context.Context, request *userV1.UpdateUserRequest) (*userV1.UpdateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (s Service) DeleteUser(ctx context.Context, request *userV1.DeleteUserRequest) (*userV1.DeleteUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}
