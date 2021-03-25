package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/db/models"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/obs"
	userV1 "github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/pb/myorg/user/v1"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserTransaction interface {
	GetUser(ctx context.Context, id string) (*userV1.User, error)
	CreateUser(ctx context.Context, item *userV1.User) error
	UpdateUser(ctx context.Context, request *userV1.UpdateUserRequest) (*userV1.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, id string) (*userV1.DeleteUserResponse, error)
}

type userTransactionImpl struct {
	tx          *sql.Tx
	redisClient RedisClient
}

func (tx *userTransactionImpl) GetUser(ctx context.Context, id string) (*userV1.User, error) {
	// Create new tracing span
	ctx, span := obs.NewSpan(ctx, "GetUser")
	defer span.End()

	// Perform query
	user, err := models.FindUser(ctx, tx.tx, id)

	// Handle error
	if err != nil {
		return nil, err
	}

	// Convert time.Time to Timestamp
	createdAt, err := ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("found user with invalid createdAt timestamp: %w", err)
	}

	// Return payload
	return &userV1.User{
		Id:        user.ID,
		CreatedAt: createdAt,
		Name:      user.Name,
	}, nil
}

func (tx *userTransactionImpl) CreateUser(ctx context.Context, item *userV1.User) error {
	ctx, span := obs.NewSpan(ctx, "CreateUser")
	defer span.End()

	createdAt, err := ptypes.Timestamp(item.CreatedAt)
	if err != nil {
		return fmt.Errorf("found user with invalid createdAt timestamp: %w", err)
	}

	user := models.User{
		ID:        item.Id,
		CreatedAt: createdAt,
		Name:      item.Name,
	}

	return user.Insert(ctx, tx.tx, boil.Infer())
}

func (tx *userTransactionImpl) UpdateUser(ctx context.Context, request *userV1.UpdateUserRequest) (*userV1.UpdateUserResponse, error) {
	ctx, span := obs.NewSpan(ctx, "UpdateUser")
	defer span.End()

	exec := tx.tx

	// TODO respect FieldMask. Do not update user's name if none is provided

	user, _ := models.FindUser(ctx, exec, request.Id)
	user.Name = request.Name
	if _, err := user.Update(ctx, exec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	createdAt, err := ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		return nil, err
	}

	userPB := &userV1.User{
		Id:        user.ID,
		Name:      user.Name,
		CreatedAt: createdAt,
	}

	return &userV1.UpdateUserResponse{
		User: userPB,
	}, nil

}

func (tx *userTransactionImpl) DeleteUser(ctx context.Context, id string) (*userV1.DeleteUserResponse, error) {
	ctx, span := obs.NewSpan(ctx, "DeleteUser")
	defer span.End()

	exec := tx.tx

	user, _ := models.FindUser(ctx, exec, id)
	if _, err := user.Delete(ctx, exec); err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &userV1.DeleteUserResponse{}, nil
}

func (tx *userTransactionImpl) cacheUser(ctx context.Context, item *userV1.User) error {
	ctx, span := obs.NewSpan(ctx, "cacheUser")
	defer span.End()

	longevity := time.Hour * 24
	return tx.redisClient.Set(ctx, redisKeyForUser(item.Id), item, longevity)
}

func redisKeyForUser(id string) string {
	return fmt.Sprintf("user-%s", id)
}
