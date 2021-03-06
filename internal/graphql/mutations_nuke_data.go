package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/obs"
)

func (s Server) buildFieldForNuke() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Drop all data",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context
			logger := obs.ToLogger(ctx)

			logger.Info().Msg("Received GraphQL request to drop all data")

			// Call the service
			if err := s.service.DropAllData(ctx); err != nil {
				return nil, err
			}

			// Build the response protobuf and return it
			return true, nil
		},
	}
}
