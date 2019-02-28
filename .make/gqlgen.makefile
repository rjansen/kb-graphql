gqlgen.setup:
	dep ensure

gqlgen:
	GO111MODULE=off go run $(GQLGEN_SCRIPT) -v -c $(GQLGEN_CONFIG)
