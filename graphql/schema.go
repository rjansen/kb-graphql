package graphql

import (
	"github.com/99designs/gqlgen/graphql"
)

func NewSchema(resolver ResolverRoot) graphql.ExecutableSchema {
	return NewExecutableSchema(
		Config{
			Resolvers: resolver,
		},
	)
}
