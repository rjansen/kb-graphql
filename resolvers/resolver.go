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
	collectionFmt = "environments/development/%s"
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
	bookProduct := struct {
		TypeName string
		types.BookWrite
	}{TypeName: "Book", BookWrite: book}

	err := r.setProduct(ctx, book.ID, bookProduct)
	if err != nil {
		return nil, err
	}

	return &types.Book{
		ID:          book.ID,
		Name:        book.Name,
		Description: book.Description,
		Value:       &book.Value,
		Isbn:        book.Isbn,
		Author:      book.Author,
		Flavor:      &book.Flavor,
	}, err
}

func (r *mutationResolver) SetAudio(ctx context.Context, audio types.AudioWrite) (*types.Audio, error) {
	audioProduct := struct {
		TypeName string
		types.AudioWrite
	}{TypeName: "Audio", AudioWrite: audio}

	err := r.setProduct(ctx, audio.ID, audioProduct)
	if err != nil {
		return nil, err
	}

	return &types.Audio{
		ID:          audio.ID,
		Name:        audio.Name,
		Description: audio.Description,
		Value:       &audio.Value,
		Singer:      audio.Singer,
		Compositor:  audio.Compositor,
		Duration:    &audio.Duration,
	}, err
}

func (r *mutationResolver) SetVideo(ctx context.Context, video types.VideoWrite) (*types.Video, error) {
	videoProduct := struct {
		TypeName string
		types.VideoWrite
	}{TypeName: "Video", VideoWrite: video}

	err := r.setProduct(ctx, video.ID, videoProduct)
	if err != nil {
		return nil, err
	}

	return &types.Video{
		ID:          video.ID,
		Name:        video.Name,
		Description: video.Description,
		Value:       &video.Value,
		Director:    video.Director,
		Writer:      video.Writer,
		Actors:      video.Actors,
		Duration:    &video.Duration,
	}, err
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

func (r *queryResolver) ProductBy(ctx context.Context, filter *types.ProductFilter) ([]types.Product, error) {
	var (
		client        = firestore.MustReference(r.tree)
		logger        = l.MustReference(r.tree)
		productRef    = fmt.Sprintf(collectionFmt, "products")
		products      []types.Product
		productsQuery firestore.Query = client.Collection(productRef)
	)

	logger.Info("resolve.product_by", l.NewValue("filter", filter))

	if filter.Name != nil {
		productsQuery = productsQuery.Where("Name", ">=", *filter.Name)
	}
	if filter.ID != nil {
		productsQuery = productsQuery.Where("ID", "==", *filter.ID)
	}
	if filter.Value != nil {
		productsQuery = productsQuery.Where("Value", "==", *filter.Value)
	}
	productsQuery = productsQuery.OrderBy("Name", firestore.Asc)

	logger.Info("resolve.product_by.query", l.NewValue("query", productsQuery))

	documents, err := productsQuery.Documents(ctx).GetAll()
	if err != nil {
		logger.Error("resolve.products_by.err", l.NewValue("error", err))
		return products, err
	}
	// TODO: Create function DataPath on raizel lib to prevent this struct usage
	type productType struct {
		TypeName string
	}

	for index, document := range documents {
		var targetType productType
		err := document.DataTo(&targetType)
		if err != nil {
			logger.Error(
				"resolve.product_by.fetch_type_err",
				l.NewValue("index", index),
				l.NewValue("error", err),
			)
			return products, err
		}

		var product types.Product
		switch targetType.TypeName {
		case "Book":
			var book types.Book
			err := document.DataTo(&book)
			if err != nil {
				logger.Error(
					"resolve.product_by.fetch_book_err",
					l.NewValue("index", index),
					l.NewValue("error", err),
				)
				return products, err
			}
			product = book
		case "Audio":
			var audio types.Audio
			err := document.DataTo(&audio)
			if err != nil {
				logger.Error(
					"resolve.product_by.fetch_audio_err",
					l.NewValue("index", index),
					l.NewValue("error", err),
				)
				return products, err
			}
			product = audio
		case "Video":
			var video types.Video
			err := document.DataTo(&video)
			if err != nil {
				logger.Error(
					"resolve.product_by.fetch_video_err",
					l.NewValue("index", index),
					l.NewValue("error", err),
				)
				return products, err
			}
			product = video
		default:
			logger.Info(
				"resolve.product_by.err_invalid_type",
				l.NewValue("type_name", targetType.TypeName),
			)
		}
		products = append(products, product)
	}

	logger.Info("resolve.produc_by.firestore.fetched", l.NewValue("products.len", len(products)))
	return products, err

}

func (r *queryResolver) Search(ctx context.Context, filter *types.ProductFilter) ([]types.SearchResult, error) {
	return nil, nil
}
