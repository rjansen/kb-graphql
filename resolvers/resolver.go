package resolvers

import (
	"context"
	"fmt"

	"github.com/rjansen/kb-graphql/graphql"
	"github.com/rjansen/kb-graphql/types"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel/firestore"
	"github.com/rjansen/yggdrasil"
)

var (
	collectionFmt = "environments/development/kb-graphql/%s"
)

type Resolver struct {
	tree yggdrasil.Tree
}

func NewResolver(tree yggdrasil.Tree) *Resolver {
	return &Resolver{tree: tree}
}

func (r *Resolver) Mutation() graphql.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) setProduct(ctx context.Context, id string, product interface{}) error {
	var (
		client     = firestore.MustReference(r.tree)
		logger     = l.MustReference(r.tree)
		productRef = fmt.Sprintf(collectionFmt, fmt.Sprintf("products/%s", id))
	)
	logger.Debug("resolve.set_product", l.NewValue("id", id), l.NewValue("product", product))

	err := client.Doc(productRef).Set(ctx, product)
	if err != nil {
		return err
	}

	logger.Debug("resolve.set_product.finished")
	return nil
}

func (r *mutationResolver) SetBook(ctx context.Context, book types.BookWrite) (*types.Book, error) {
	var result types.Book

	err := r.setProduct(ctx, book.ID, book)
	if err != nil {
		return &result, err
	}

	return &result, err
}

func (r *mutationResolver) SetAudio(ctx context.Context, audio types.AudioWrite) (*types.Audio, error) {
	var result types.Audio

	err := r.setProduct(ctx, audio.ID, audio)
	if err != nil {
		return &result, err
	}

	return &result, err
}

func (r *mutationResolver) SetVideo(ctx context.Context, video types.VideoWrite) (*types.Video, error) {
	var result types.Video

	err := r.setProduct(ctx, video.ID, video)
	if err != nil {
		return &result, err
	}

	return &result, err
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) product(ctx context.Context, id string, product interface{}) error {
	var (
		client     = firestore.MustReference(r.tree)
		logger     = l.MustReference(r.tree)
		productRef = fmt.Sprintf(collectionFmt, fmt.Sprintf("products/%s", id))
	)
	logger.Debug("resolve.product", l.NewValue("id", id))

	document, err := client.Doc(productRef).Get(ctx)
	if err != nil {
		return err
	}

	err = document.DataTo(product)
	if err != nil {
		return err
	}

	logger.Debug("resolve.product.fetched", l.NewValue("produc", product))
	return nil
}

func (r *queryResolver) Book(ctx context.Context, id string) (*types.Book, error) {
	var book types.Book

	err := r.product(ctx, id, &book)
	if err != nil {
		return nil, err
	}

	return &book, err
}

func (r *queryResolver) Audio(ctx context.Context, id string) (*types.Audio, error) {
	var audio types.Audio

	err := r.product(ctx, id, &audio)
	if err != nil {
		return nil, err
	}

	return &audio, err
}

func (r *queryResolver) Video(ctx context.Context, id string) (*types.Video, error) {
	var video types.Video

	err := r.product(ctx, id, &video)
	if err != nil {
		return nil, err
	}

	return &video, err
}

func (r *queryResolver) ProductBy(ctx context.Context, filter *types.ProductFilter) (types.Product, error) {
	return nil, nil
}

func (r *queryResolver) Search(ctx context.Context, filter *types.ProductFilter) (types.SearchResult, error) {
	return nil, nil
}
