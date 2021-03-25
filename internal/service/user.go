package service

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/db"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/obs"
	userV1 "github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/pb/myorg/user/v1"
	"github.com/rs/xid"
)

func (s Service) DropAllData(ctx context.Context) error {
	return s.dbClient.DropAllData(ctx)
}

func (s Service) CreateUser(ctx context.Context, request *userV1.CreateUserRequest) (*userV1.CreateUserResponse, error) {
	// TODO add validation

	user := &userV1.User{
		Id:        xid.New().String(),
		CreatedAt: ptypes.TimestampNow(),
		Name:      request.Name,
	}

	if err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return tx.CreateUser(ctx, user)
	}); err != nil {
		return nil, err
	}

	return &userV1.CreateUserResponse{
		User: user,
	}, nil
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
	var response *userV1.UpdateUserResponse

	// TODO add validation

	if err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		if res, err := tx.UpdateUser(ctx, request); err != nil {
			return err
		} else {
			response = res
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return response, nil
}

func (s Service) DeleteUser(ctx context.Context, request *userV1.DeleteUserRequest) (*userV1.DeleteUserResponse, error) {
	var response *userV1.DeleteUserResponse
	if err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		if res, err := tx.DeleteUser(ctx, request.Id); err != nil {
			return err
		} else {
			response = res
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return response, nil
}
