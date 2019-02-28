package resolvers

import (
	"context"

	"github.com/rjansen/kb-graphql/graphql"
	"github.com/rjansen/kb-graphql/types"
)

type Resolver struct{}

func (r *Resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Book(ctx context.Context, id string) (*types.Book, error) {
	panic("not implemented")
}
func (r *queryResolver) Audio(ctx context.Context, id string) (*types.Audio, error) {
	panic("not implemented")
}
func (r *queryResolver) Video(ctx context.Context, id string) (*types.Video, error) {
	panic("not implemented")
}
func (r *queryResolver) ProductBy(ctx context.Context, filter *types.ProductFilter) (types.Product, error) {
	panic("not implemented")
}
func (r *queryResolver) Search(ctx context.Context, filter *types.ProductFilter) (types.SearchResult, error) {
	panic("not implemented")
}
