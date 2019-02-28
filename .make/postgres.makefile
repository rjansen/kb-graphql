.PHONY: postgres.setup
postgres.setup: $(MIGRATE)
	@echo "$(ROOT_REPO)@$(BUILD) postgres.setup"

$(MIGRATE_FILE): | $(DEPS_DIR)
	@echo "$(ROOT_REPO)@$(BUILD) $(MIGRATE_FILE)"
	curl -o $(MIGRATE_FILE) -L $(MIGRATE_URL)

$(MIGRATE): | $(MIGRATE_FILE)
	@echo "$(ROOT_REPO)@$(BUILD) $(MIGRATE)"
	@cd $(DEPS_DIR) && tar xf $(MIGRATE_NAME) && mv -f $(MIGRATE_FILENAME) $(TOOLS_DIR)/$(MIGRATE_BIN)
	$(MIGRATE) -help > /dev/null 2>&1

.PHONY: postgres.run
postgres.run:
	@echo "$(REPO) postgres.run"
	docker run --rm -d --name postgres-run --net=host \
		-v "$(POSTGRES_SCRIPTS_DIR):/docker-entrypoint-initdb.d" postgres:$(POSTGRES_VERSION)
	@sleep 7 #wait until database is ready

.PHONY: postgres.kill
postgres.kill:
	@echo "$(REPO) postgres.kill"
	docker kill postgres-run

.PHONY: postgres.scripts
postgres.scripts:
	@echo "$(REPO) postgres.scripts"
	sleep 5
	cat $(POSTGRES_SCRIPTS_DIR)/* | \
		psql -h $(POSTGRES_HOST) -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE) -1 -f -

.PHONY: postgres.psql
postgres.psql:
	@echo "$(REPO) postgres.psql"
	docker run --rm -it --name psql-run --net=host postgres:$(POSTGRES_VERSION) \
		psql -h $(POSTGRES_HOST) -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE)

.PHONY: postgres.migrations.run
postgres.migrations.run: postgres.run postgres.migrations.up

postgres.migrations.%:
	migrate -source file://$(POSTGRES_MIGRATIONS_DIR) -database $(SQL_DSN) $*
