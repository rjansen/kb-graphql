package resolvers

import (
	"context"
	"fmt"
	"testing"

	"github.com/rjansen/kb-graphql/types"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel/firestore"
	fmock "github.com/rjansen/raizel/firestore/mock"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testQueryResolver struct {
	name      string
	tree      yggdrasil.Tree
	mockSetup map[string]interface{}
	bookID    string
	audioID   string
	videoID   string
	productBy types.ProductFilter
	search    types.ProductFilter
}

func (scenario *testQueryResolver) setup(t *testing.T) {
	var (
		roots     = yggdrasil.NewRoots()
		errLogger = l.Register(&roots, l.NewZapLoggerDefault())
		client    = fmock.NewClientMock()
		errClient = firestore.Register(&roots, client)
	)
	require.Nil(t, errLogger, "setup logger error")
	require.Nil(t, errClient, "setup firestore error")

	if len(scenario.mockSetup) > 0 {

		for documentKey, document := range scenario.mockSetup {
			var (
				snapshot    = fmock.NewDocumentSnapshotMock()
				documentRef = fmock.NewDocumentRefMock()
				dataToMock  func(mock.Arguments)
			)

			switch realDocument := document.(type) {
			case types.Book:
				dataToMock = func(args mock.Arguments) {
					arg := args.Get(0)
					if arg != nil {
						product := arg.(*types.Book)
						*product = realDocument
					}
				}
			case types.Audio:
				dataToMock = func(args mock.Arguments) {
					arg := args.Get(0)
					if arg != nil {
						product := arg.(*types.Audio)
						*product = realDocument
					}
				}

			case types.Video:
				dataToMock = func(args mock.Arguments) {
					arg := args.Get(0)
					if arg != nil {
						product := arg.(*types.Video)
						*product = realDocument
					}
				}
			}

			snapshot.On("DataTo", mock.Anything).Run(dataToMock).Return(nil)
			documentRef.On("Get", mock.Anything).Return(snapshot, nil)
			client.On("Doc", fmt.Sprintf(collectionFmt, documentKey)).Return(documentRef)
		}
	}

	scenario.tree = roots.NewTreeDefault()
}

func (scenario *testQueryResolver) tearDown(*testing.T) {
	if scenario.tree != nil {
		scenario.tree.Close()
	}
}

func TestQueryResolver(test *testing.T) {
	scenarios := []testQueryResolver{
		{
			name: "Resolves all query fields successfully",
			mockSetup: map[string]interface{}{
				"products/mock_book":  types.Book{},
				"products/mock_audio": types.Audio{},
				"products/mock_video": types.Video{},
			},
			bookID:  "mock_book",
			audioID: "mock_audio",
			videoID: "mock_video",
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("%d-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				var (
					ctx           = context.Background()
					queryResolver = queryResolver{
						&Resolver{scenario.tree},
					}
				)
				require.NotNil(t, queryResolver, "query_resolver invalid instance")
				require.NotNil(t, queryResolver.Resolver, "resolver invalid instance")

				book, errBook := queryResolver.Book(ctx, scenario.bookID)
				audio, errAudio := queryResolver.Audio(ctx, scenario.audioID)
				video, errVideo := queryResolver.Video(ctx, scenario.videoID)

				errors := map[string]error{
					"Book": errBook, "Audio": errAudio, "Video": errVideo,
				}
				results := map[string]interface{}{
					"Book": book, "Audio": audio, "Video": video,
				}

				for key, err := range errors {
					assert.Nilf(t, err, "field resolve error: field=%s", key)
				}
				for key, result := range results {
					assert.NotZerof(t, result, "field invalid instance: field=%s", key)
				}
			},
		)
	}
}

type testMutationResolver struct {
	name      string
	tree      yggdrasil.Tree
	mockSetup map[string]interface{}
	book      types.BookWrite
	audio     types.AudioWrite
	video     types.VideoWrite
}

func (scenario *testMutationResolver) setup(t *testing.T) {
	var (
		roots     = yggdrasil.NewRoots()
		errLogger = l.Register(&roots, l.NewZapLoggerDefault())
		client    = fmock.NewClientMock()
		errClient = firestore.Register(&roots, client)
	)
	require.Nil(t, errLogger, "setup logger error")
	require.Nil(t, errClient, "setup firestore error")

	if len(scenario.mockSetup) > 0 {

		for documentKey, _ := range scenario.mockSetup {
			var documentRef = fmock.NewDocumentRefMock()

			documentRef.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			client.On("Doc", fmt.Sprintf(collectionFmt, documentKey)).Return(documentRef)
		}
	}

	scenario.tree = roots.NewTreeDefault()
}

func (scenario *testMutationResolver) tearDown(*testing.T) {
	if scenario.tree != nil {
		scenario.tree.Close()
	}
}

func TestMutationResolver(test *testing.T) {
	scenarios := []testMutationResolver{
		{
			name: "Resolves all mutation fields successfully",
			mockSetup: map[string]interface{}{
				"products/mock_book":  types.Book{},
				"products/mock_audio": types.Audio{},
				"products/mock_video": types.Video{},
			},
			book:  types.BookWrite{ID: "mock_book"},
			audio: types.AudioWrite{ID: "mock_audio"},
			video: types.VideoWrite{ID: "mock_video"},
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("%d-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				var (
					ctx              = context.Background()
					mutationResolver = mutationResolver{
						&Resolver{scenario.tree},
					}
				)
				require.NotNil(t, mutationResolver, "mutation_resolver invalid instance")
				require.NotNil(t, mutationResolver.Resolver, "resolver invalid instance")

				book, errBook := mutationResolver.SetBook(ctx, scenario.book)
				audio, errAudio := mutationResolver.SetAudio(ctx, scenario.audio)
				video, errVideo := mutationResolver.SetVideo(ctx, scenario.video)

				errors := map[string]error{
					"Book": errBook, "Audio": errAudio, "Video": errVideo,
				}
				results := map[string]interface{}{
					"Book": book, "Audio": audio, "Video": video,
				}

				for key, err := range errors {
					assert.Nilf(t, err, "field resolve error: field=%s", key)
				}
				for key, result := range results {
					assert.NotZerof(t, result, "field invalid instance: field=%s", key)
				}
			},
		)
	}
}
