package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/obs"
	userV1 "github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/pb/myorg/user/v1"
)

func (s Server) buildFieldForCreateUser(userType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        userType,
		Description: "Create new User",
		Args: graphql.FieldConfigArgument{
			argName: &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The User's name",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context
			logger := obs.ToLogger(ctx)

			args := p.Args
			logger.Info().Msgf("Received GraphQL request to %s with args: %v", p.Info.FieldName, args)

			// Build the request protobuf from the GraphQL args
			request, err := buildCreateUserRequestFromArgs(args)
			if err != nil {
				return nil, err
			}

			// Call the service
			res, err := s.service.CreateUser(ctx, request)
			if err != nil {
				return nil, err
			}

			// Build the response protobuf and return it
			return buildUser(res.User)
		},
	}
}

func buildCreateUserRequestFromArgs(args map[string]interface{}) (*userV1.CreateUserRequest, error) {
	if value, ok := args[argName]; ok {
		if val, ok := value.(string); ok {
			return &userV1.CreateUserRequest{
				Name: val,
			}, nil
		} else {
			return nil, fmt.Errorf("'%s' not a string", argName)
		}
	} else {
		return nil, fmt.Errorf("must specify '%s'", argName)
	}
}
