package graphql

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/graphql-go/graphql"
	userV1 "github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/pb/myorg/user/v1"
)

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
}

func buildUser(in *userV1.User) (User, error) {
	createdAt, err := ptypes.Timestamp(in.CreatedAt)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:        in.Id,
		CreatedAt: createdAt,
		Name:      in.Name,
	}, nil
}

func buildTypeForUser() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.String,
				Description: "The User's ID",
			},
			"createdAt": &graphql.Field{
				Type:        graphql.DateTime,
				Description: "The User's creation time",
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The User's name",
			},
		},
	})
}
