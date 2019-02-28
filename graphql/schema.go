package model

import (
	"github.com/rjansen/fivecolors/core/errors"
	"github.com/rjansen/fivecolors/core/graphql"
)

var (
	ErrInvalidState = errors.New("ErrInvalidState")
)

func NewSchema(resolver *Resolver) graphql.Schema {
	return graphql.NewSchema(
		NewExecutableSchema(
			Config{
				Resolvers: resolver,
			},
		),
	)
}
