package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/obs"
	userV1 "github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/pb/myorg/user/v1"
)

func (s Server) buildFieldForDeleteUser(userType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Delete a User",
		Args: graphql.FieldConfigArgument{
			argID: &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The ID of the User the client wishes to delete",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context
			logger := obs.ToLogger(ctx)

			args := p.Args
			logger.Info().Msgf("Received GraphQL request to %s with args: %v", p.Info.FieldName, args)

			// Build the request protobuf from the GraphQL args
			request, err := buildDeleteUserRequestFromArgs(args)
			if err != nil {
				return false, err
			}

			// Call the service
			_, err = s.service.DeleteUser(ctx, request)
			if err != nil {
				return false, err
			}

			// Build the response protobuf and return it
			return true, nil
		},
	}
}

func buildDeleteUserRequestFromArgs(args map[string]interface{}) (*userV1.DeleteUserRequest, error) {
	if value, ok := args[argID]; ok {
		if val, ok := value.(string); ok {
			return &userV1.DeleteUserRequest{
				Id: val,
			}, nil
		} else {
			return nil, fmt.Errorf("'%s' not a string", argID)
		}
	} else {
		return nil, fmt.Errorf("must specify '%s'", argID)
	}
}
