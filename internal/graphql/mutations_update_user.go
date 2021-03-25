package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/obs"
	userV1 "github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/pb/myorg/user/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (s Server) buildFieldForUpdateUser(userType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        userType,
		Description: "Update a User",
		Args: graphql.FieldConfigArgument{
			argID: &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The ID of the User the client wishes to update",
			},
			argName: &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "The User's new name. This argument is optional.",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context
			logger := obs.ToLogger(ctx)

			args := p.Args
			logger.Info().Msgf("Received GraphQL request to %s with args: %v", p.Info.FieldName, args)

			// TODO get selection set and use FieldMask

			// Build the request protobuf from the GraphQL args
			request, err := buildUpdateUserRequestFromArgs(args)
			if err != nil {
				return nil, err
			}

			// Call the service
			res, err := s.service.UpdateUser(ctx, request)
			if err != nil {
				return nil, err
			}

			// Build the response protobuf and return it
			return buildUser(res.User)
		},
	}
}

func buildUpdateUserRequestFromArgs(args map[string]interface{}) (*userV1.UpdateUserRequest, error) {
	var paths []string
	request := &userV1.UpdateUserRequest{}

	if value, ok := args[argID]; ok {
		// TODO do we need these type assertions or does the graphql library take care of that for us?
		if val, ok := value.(string); ok {
			request.Id = val
		} else {
			return nil, fmt.Errorf("'%s' not a string", argID)
		}
	} else {
		return nil, fmt.Errorf("must specify '%s'", argID)
	}

	if value, ok := args[argName]; ok {
		if val, ok := value.(string); ok {
			request.Name = val
			paths = append(paths, argName)
		} else {
			return nil, fmt.Errorf("'%s' not a string", argName)
		}
	} else {
		return nil, fmt.Errorf("must specify '%s'", argName)
	}

	request.FieldMask = &fieldmaskpb.FieldMask{
		Paths: paths,
	}

	return request, nil
}
