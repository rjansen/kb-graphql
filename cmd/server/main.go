package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rjansen/kb-graphql/api"
	"github.com/rjansen/kb-graphql/graphql"
	"github.com/rjansen/kb-graphql/resolvers"
	"github.com/rjansen/l"
	"github.com/rjansen/migi"
	"github.com/rjansen/raizel/firestore"
	"github.com/rjansen/yggdrasil"
)

var (
	version string
)

type options struct {
	bindAddress string
	projectID   string
}

func newOptions() options {
	var (
		env     = migi.NewOptions(migi.NewEnvironmentSource())
		options options
	)
	env.StringVar(
		&options.bindAddress, "server_bindaddress", ":8080", "Server bind address, ip:port",
	)
	env.StringVar(
		&options.projectID, "project_id", "project-id", "GCP project identifier",
	)
	env.Parse()
	return options
}

func newTree(options options) yggdrasil.Tree {
	var (
		logger = l.NewZapLoggerDefault()
		roots  = yggdrasil.NewRoots()
		err    error
	)

	err = l.Register(&roots, logger)
	if err != nil {
		panic(err)
	}

	err = firestore.Register(&roots, newFirestoreClient(options))
	if err != nil {
		panic(err)
	}

	return roots.NewTreeDefault()
}

func newSchema(tree yggdrasil.Tree, options options) graphql.Schema {
	return graphql.NewSchema(
		resolvers.NewResolver(tree),
	)
}

func newFirestoreClient(options options) firestore.Client {
	return firestore.NewClient(options.projectID)
}

func httpRouterHandler(handler http.HandlerFunc) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handler(w, r)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "alive")
}

func main() {
	var (
		options        = newOptions()
		tree           = newTree(options)
		schema         = newSchema(tree, options)
		logger         = l.MustReference(tree)
		graphqlHandler = httpRouterHandler(api.NewGraphQLHandler(tree, schema))
		router         = httprouter.New()
	)

	logger.Info("server.router.init")

	router.GET("/api/healthcheck", httpRouterHandler(healthCheck))
	router.GET("/api/query", graphqlHandler)
	router.POST("/api/query", graphqlHandler)

	server := &http.Server{
		Addr:    options.bindAddress,
		Handler: router,
	}

	logger.Info("server.router.created")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	go func() {
		logger.Info("server.starting", l.NewValue("address", options.bindAddress))

		if err := server.ListenAndServe(); err != nil {
			logger.Error(
				"server.err", l.NewValue("error", err), l.NewValue("address", options.bindAddress),
			)
		}
	}()

	<-stop

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger.Info("server.shutdown")
	server.Shutdown(shutDownCtx)
	logger.Info("server.shutdown.gracefully")
}
