.PHONY: ci
ci: vet coverage.text bench

.PHONY: clearcache
clearcache:
	@echo "$(ROOT_REPO)@$(BUILD) clearcache"
	-rm -Rf $(DEPS_DIR)
	-$(HOME)/.cache/go-build
	-$(HOME)/gopath/pkg/mod
	-$(foreach path,$(BIN_LIST),ls -ld $(BASE_DIR)/$(path);)

$(DEPS_DIR):
	mkdir -p $(DEPS_DIR)

.PHONY: bin.setup
bin.setup: $(GOTESTSUM) $(CODECOV)
	@echo "$(ROOT_REPO)@$(BUILD) bin.setup bin=$(CODECOV) file=$(CODECOV_FILE)"

$(GOTESTSUM_FILE): | $(DEPS_DIR)
	@echo "$(ROOT_REPO)@$(BUILD) $(GOTESTSUM_FILE)"
	curl -o $(GOTESTSUM_FILE) -L $(GOTESTSUM_URL)

$(GOTESTSUM): | $(GOTESTSUM_FILE)
	@echo "$(ROOT_REPO)@$(BUILD) $(GOTESTSUM)"
	@cd $(DEPS_DIR) && tar xf $(GOTESTSUM_NAME) && mv -f $(GOTESTSUM_BIN) $(TOOLS_DIR)
	$(GOTESTSUM) --help > /dev/null 2>&1

$(CODECOV_FILE): | $(DEPS_DIR)
	@echo "$(ROOT_REPO)@$(BUILD) $(CODECOV_FILE)"
	curl -o $(CODECOV_FILE) -L $(CODECOV_URL)

$(CODECOV): | $(CODECOV_FILE)
	@echo "$(ROOT_REPO)@$(BUILD) $(CODECOV)"
	cd $(DEPS_DIR) && chmod a+x $(CODECOV_NAME) && cp -f $(CODECOV_BIN) $(TOOLS_DIR)
	$(CODECOV) -h > /dev/null 2>&1

.PHONY: bin.setup.debug
bin.setup.debug: dlv.setup
	@echo "$(ROOT_REPO)@$(BUILD) module.setup.debug"

.PHONY: delv.setup
dlv.setup:
	@echo "$(ROOT_REPO)@$(BUILD) dlv.setup"
	which dlv || go get -u github.com/derekparker/delve/cmd/dlv
	dlv version > /dev/null 2>&1

.PHONY: build
build: $(TMP_DIR)
	@echo "$(REPO)@$(BUILD) build"
	go build -ldflags '-X main.version=$(BUILD)' -o $(BIN_NAME) $(REPO)

.PHONY: run
run: build
	@echo "$(REPO)@$(BUILD) run"
	$(BIN_NAME)

.PHONY: vendor
vendor:
	@echo "$(REPO)@$(BUILD) vendor"
	go mod verify && go mod vendor

.PHONY: debug
debug:
	@echo "$(REPO)@$(BUILD) debug"
	dlv debug $(REPO)

.PHONY: debugtest
debugtest:
	@echo "$(REPO)@$(BUILD) debugtest"
	dlv test $(DEBUG_PKG) --build-flags="-ldflags '-X $(ITEST_ROOT_DOC)' -tags 'integration'" -- \
		-test.run $(TESTS)

.PHONY: vet
vet:
	@echo "$(REPO)@$(BUILD) vet"
	go vet $(TEST_PKGS)

.PHONY: test
test:
	@echo "$(REPO)@$(BUILD) test $(REPO) - $(TEST_PKGS)"
	$(GOTESTSUM) -f $(TEST_VERBOSITY) -- -v -race -run $(TESTS) $(TEST_PKGS)

.PHONY: itest
itest:
	@echo "$(REPO)@$(BUILD) itest"
	$(GOTESTSUM) -f $(TEST_VERBOSITY) -- -ldflags '-X $(ITEST_ROOT_DOC)' -tags="integration" \
		-v -race -run $(TESTS) $(TEST_PKGS)

.PHONY: bench
bench:
	@echo "$(REPO)@$(BUILD) bench"
	$(GOTESTSUM) -f $(TEST_VERBOSITY) -- -ldflags '-X $(ITEST_ROOT_DOC)' \
		-bench=. -run="^$$" -benchmem $(TEST_PKGS)

.PHONY: coverage
coverage: $(TMP_DIR)
	@echo "$(REPO)@$(BUILD) coverage"
	ls -ld $(TMP_DIR)
	@touch $(COVERAGE_FILE)
	$(GOTESTSUM) -f $(TEST_VERBOSITY) -- -ldflags '-X $(ITEST_ROOT_DOC)' -tags="integration" \
		-v -run $(TESTS) -covermode=atomic -coverpkg=$(PKGS) -coverprofile=$(COVERAGE_FILE) $(TEST_PKGS)

.PHONY: coverage.text
coverage.text: coverage
	@echo "$(REPO)@$(BUILD) coverage.text"
	go tool cover -func=$(COVERAGE_FILE)

.PHONY: coverage.html
coverage.html: coverage
	@echo "$(REPO)@$(BUILD) coverage.html"
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@open $(COVERAGE_HTML) || google-chrome $(COVERAGE_HTML) || google-chrome-stable $(COVERAGE_HTML)

.PHONY: coverage.push
coverage.push:
	@echo "$(REPO) coverage.push"
	$(CODECOV) -f $(COVERAGE_FILE)$(if $(CODECOV_TOKEN), -t $(CODECOV_TOKEN),)
