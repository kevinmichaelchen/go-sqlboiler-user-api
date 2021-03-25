package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/obs"
	userV1 "github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/pb/myorg/user/v1"
)

func (s Server) buildFieldForGetUser(userType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type: userType,
		Args: graphql.FieldConfigArgument{
			argID: &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "The User's ID",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context
			logger := obs.ToLogger(ctx)

			args := p.Args
			logger.Info().Msgf("Received GraphQL request to %s with args: %v", p.Info.FieldName, args)

			// Build the request protobuf from the GraphQL args
			request, err := buildGetUserRequestFromArgs(args)
			if err != nil {
				return nil, err
			}

			// Call the service
			res, err := s.service.GetUser(ctx, request)
			if err != nil {
				return nil, err
			}

			// Build the response protobuf and return it
			return buildUser(res.User)
		},
		Description: "Retrieve a User object",
	}
}

func buildGetUserRequestFromArgs(args map[string]interface{}) (*userV1.GetUserRequest, error) {
	if value, ok := args[argID]; ok {
		if val, ok := value.(string); ok {
			return &userV1.GetUserRequest{
				Id: val,
			}, nil
		} else {
			return nil, fmt.Errorf("'%s' not a string", argID)
		}
	} else {
		return nil, fmt.Errorf("must specify '%s'", argID)
	}
}
