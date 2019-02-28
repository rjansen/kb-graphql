package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rjansen/fivecolors/core/graphql"
	"github.com/rjansen/fivecolors/core/graphql/mockschema"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	raizelmock "github.com/rjansen/raizel/mock"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
)

type testGraphQL struct {
	name           string
	tree           yggdrasil.Tree
	repository     *raizelmock.MockRepository
	method         string
	path           string
	body           string
	contentType    string
	request        *http.Request
	response       *httptest.ResponseRecorder
	responseStatus int
}

func (scenario *testGraphQL) setup(t *testing.T) {
	var (
		roots         = yggdrasil.NewRoots()
		repository    = raizelmock.NewMockRepository()
		errLogger     = l.Register(&roots, l.NewZapLoggerDefault())
		errRepository = raizel.Register(&roots, repository)
		errSchema     = graphql.Register(&roots, graphql.NewSchema(mockschema.New()))
	)
	require.Nil(t, errLogger, "setup logger error")
	require.Nil(t, errRepository, "setup repository error")
	require.Nil(t, errSchema, "setup schema error")
	require.Nil(t, errSchema, "setup schema error")

	r := httptest.NewRequest(
		scenario.method, scenario.path, strings.NewReader(scenario.body),
	)
	r.Header.Set("content-type", scenario.contentType)

	scenario.tree = roots.NewTreeDefault()
	scenario.repository = repository
	scenario.request = r
	scenario.response = httptest.NewRecorder()
}

func (scenario *testGraphQL) tearDown(*testing.T) {}

func TestGraphQL(test *testing.T) {
	scenarios := []testGraphQL{
		{
			name:           "When handler receives a HEAD request returns method not allowed",
			method:         http.MethodHead,
			path:           "/query",
			responseStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "When handler receives a PUT request returns method not allowed",
			method:         http.MethodPut,
			path:           "/query",
			responseStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "When handler receives a PATCH request returns method not allowed",
			method:         http.MethodPatch,
			path:           "/query",
			responseStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "When handler receives a DELETE request returns method not allowed",
			method:         http.MethodDelete,
			path:           "/query",
			responseStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "When handler receives a CONNECT request returns method not allowed",
			method:         http.MethodConnect,
			path:           "/query",
			responseStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "When handler receives a OPTIONS request returns method not allowed",
			method:         http.MethodOptions,
			path:           "/query",
			responseStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "When handler receives a TRACE request returns method not allowed",
			method:         http.MethodTrace,
			path:           "/query",
			responseStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "When a GET request has a blank query parameter returns bad request",
			method:         http.MethodGet,
			path:           "/query?q=",
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "When a POST content type is invalid returns bad request",
			method:         http.MethodPost,
			path:           "/query",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/xml",
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "When a POST request body has a invalid graphql content returns internal server error",
			method:         http.MethodPost,
			path:           "/query",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/graphql",
			responseStatus: http.StatusInternalServerError,
		},
		{
			name:           "When a POST request body has a invalid json content returns bad request",
			method:         http.MethodPost,
			path:           "/query",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/json",
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "When a POST request body has a invalid graphql query returns internal server error",
			method:         http.MethodPost,
			path:           "/query",
			body:           `{"query": "<xml><id>xmlid</id></name>Invalid Body</name></xml>"}`,
			contentType:    "application/json",
			responseStatus: http.StatusInternalServerError,
		},
		{
			name:           "When a POST request body is graphql executes the query and returns ok with query results",
			method:         http.MethodPost,
			path:           "/query",
			body:           "{me{tid,user{id,name}}}",
			contentType:    "application/graphql",
			responseStatus: http.StatusOK,
		},
		{
			name:           "When a POST request body is graphql executes the query and returns ok with query results",
			method:         http.MethodPost,
			path:           "/query",
			body:           "{mockEntity{tid,entity{id,string,float,integer,dateTime}}}",
			contentType:    "application/graphql",
			responseStatus: http.StatusOK,
		},
		{
			name:           "When a POST request body is a valid json executes the query and returns ok with query results",
			method:         http.MethodPost,
			path:           "/query",
			body:           `{"query": "{me{tid,user{id,name}}}"}`,
			contentType:    "application/json",
			responseStatus: http.StatusOK,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				graphql := NewGraphQLHandler(scenario.tree)
				graphql(scenario.response, scenario.request)
				result := scenario.response.Result()
				/*
					var (
						resultMap map[string]interface{}
						errDecode = json.NewDecoder(result.Body).Decode(&resultMap)
					)
					require.Equal(t, map[string]interface{}{}, resultMap, "result invalid instance")
					require.Nil(t, errDecode, "response body invalid")
					require.NotZero(t, resultMap, "result invalid")
				*/
				require.Equal(t, scenario.responseStatus, result.StatusCode, "invalid response statuscode")
			},
		)
	}
}
