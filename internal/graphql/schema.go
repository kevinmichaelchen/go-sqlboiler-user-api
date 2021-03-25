package graphql

import (
	"github.com/rs/zerolog/log"

	"github.com/graphql-go/graphql"
)

func (s Server) buildSchema() *graphql.Schema {
	userType := buildTypeForUser()

	queryFields := graphql.Fields{
		"user": s.buildFieldForGetUser(userType),
	}

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: queryFields,
	})

	mutationFields := graphql.Fields{
		"createUser": s.buildFieldForCreateUser(userType),
		"updateUser": s.buildFieldForUpdateUser(userType),
		"deleteUser": s.buildFieldForDeleteUser(userType),
		"nuke":       s.buildFieldForNuke(),
	}

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: mutationFields,
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create GraphQL schema")
	}

	return &schema
}
