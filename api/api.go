package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rjansen/kb-graphql/graphql"
	"github.com/rjansen/kb-graphql/validator"
	"github.com/rjansen/l"
	"github.com/rjansen/yggdrasil"
)

func NewGraphQLHandler(tree yggdrasil.Tree, schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GraphQL(tree, schema, w, r)
	}
}

func GraphQL(tree yggdrasil.Tree, schema graphql.Schema, w http.ResponseWriter, r *http.Request) {
	var (
		logger      = l.MustReference(tree)
		contentType = r.Header.Get("Content-Type")
		q           graphql.Request
	)

	logger.Info("api.request.try")

	switch r.Method {
	case http.MethodGet:
		q.Query = r.URL.Query().Get("query")
		q.OperationName = r.URL.Query().Get("operationName")
		if variables := r.URL.Query().Get("variables"); variables != "" {
			if err := json.NewDecoder(strings.NewReader(variables)).Decode(&q.Variables); err != nil {
				logger.Error("api.request.variables.err", l.NewValue("error", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	case http.MethodPost:
		switch {
		case strings.HasPrefix("application/graphql", contentType):
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logger.Error("api.request.err", l.NewValue("error", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			q = graphql.Request{Query: string(body)}
		default:
			err := json.NewDecoder(r.Body).Decode(&q)
			if err != nil {
				logger.Error("graphql.json.request.err", l.NewValue("error", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := validator.IsBlank(q.Query); err != nil {
		logger.Error("graphql.request.err", l.NewValue("error", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	logger.Debug("graphql.query.try", l.NewValue("query", q))
	result := graphql.Execute(tree, schema, q)
	if len(result.Errors) > 0 {
		logger.Error("graphql.query.err", l.NewValue("query", q), l.NewValue("result", result))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		logger.Debug("graphql.query.result", l.NewValue("query", q), l.NewValue("result", result))
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(result)
}
