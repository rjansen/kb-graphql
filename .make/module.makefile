.PHONY: ci
ci: vet coverage.text bench

.PHONY: clearcache
clearcache:
	@echo "$(ROOT_REPO)@$(BUILD) clearcache"
	-rm -Rf $(DEPS_DIR)
	-$(HOME)/.cache/go-build
	-$(HOME)/gopath/pkg/mod
	-$(foreach path,$(MODULE_LIST),ls -ld $(BASE_DIR)/$(path)/vendor;)

$(DEPS_DIR):
	mkdir -p $(DEPS_DIR)

.PHONY: module.setup
module.setup: $(GOTESTSUM) $(CODECOV)
	@echo "$(ROOT_REPO)@$(BUILD) module.setup bin=$(CODECOV) file=$(CODECOV_FILE)"

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

.PHONY: setup.debug
module.setup.debug: dlv.setup
	@echo "$(ROOT_REPO)@$(BUILD) module.setup.debug"

.PHONY: delv.setup
dlv.setup:
	@echo "$(ROOT_REPO)@$(BUILD) dlv.setup"
	which dlv || go get -u github.com/derekparker/delve/cmd/dlv
	dlv version > /dev/null 2>&1

.PHONY: build
build: $(TMP_DIR)
	@echo "$(REPO)@$(BUILD) build"
	cd $(MODULE_DIR) && go build -ldflags '-X main.version=$(BUILD)' -o $(MODULE_BIN) $(REPO)

.PHONY: run
run: build
	@echo "$(REPO)@$(BUILD) run"
	$(MODULE_BIN)

.PHONY: vendor
vendor:
	@echo "$(REPO)@$(BUILD) vendor"
	cd $(MODULE_DIR) && go mod verify && go mod vendor

.PHONY: debug
debug:
	@echo "$(REPO)@$(BUILD) debug"
	cd $(MODULE_DIR) dlv debug $(REPO)

.PHONY: debugtest
debugtest:
	@echo "$(REPO)@$(BUILD) debugtest"
	cd $(MODULE_DIR) && dlv test $(DEBUG_PKG) \
		--build-flags="-ldflags '-X $(ITEST_ROOT_DOC)' -tags 'integration'" -- -test.run $(TESTS)

.PHONY: vet
vet:
	@echo "$(REPO)@$(BUILD) vet"
	cd $(MODULE_DIR) && go vet $(TEST_PKGS)

.PHONY: test
test:
	@echo "$(REPO)@$(BUILD) test $(REPO) - $(TEST_PKGS)"
	cd $(MODULE_DIR) && $(GOTESTSUM) -f $(TEST_VERBOSITY) -- -v -race -run $(TESTS) $(TEST_PKGS)

.PHONY: itest
itest:
	@echo "$(REPO)@$(BUILD) itest"
	cd $(MODULE_DIR) && $(GOTESTSUM) -f $(TEST_VERBOSITY) -- -ldflags '-X $(ITEST_ROOT_DOC)' -tags="integration" \
		-v -race -run $(TESTS) $(TEST_PKGS)

.PHONY: bench
bench:
	@echo "$(REPO)@$(BUILD) bench"
	cd $(MODULE_DIR) && $(GOTESTSUM) -f $(TEST_VERBOSITY) -- -ldflags '-X $(ITEST_ROOT_DOC)' \
		-bench=. -run="^$$" -benchmem $(TEST_PKGS)

.PHONY: coverage
coverage: $(TMP_DIR)
	@echo "$(REPO)@$(BUILD) coverage"
	ls -ld $(TMP_DIR)
	@touch $(COVERAGE_FILE)
	cd $(MODULE_DIR) && $(GOTESTSUM) -f $(TEST_VERBOSITY) -- -ldflags '-X $(ITEST_ROOT_DOC)' -tags="integration" \
		-v -run $(TESTS) -covermode=atomic -coverpkg=$(PKGS) -coverprofile=$(COVERAGE_FILE) $(TEST_PKGS)

.PHONY: coverage.text
coverage.text: coverage
	@echo "$(REPO)@$(BUILD) coverage.text"
	cd $(MODULE_DIR) && go tool cover -func=$(COVERAGE_FILE)

.PHONY: coverage.html
coverage.html: coverage
	@echo "$(REPO)@$(BUILD) coverage.html"
	cd $(MODULE_DIR) && go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@open $(COVERAGE_HTML) || google-chrome $(COVERAGE_HTML) || google-chrome-stable $(COVERAGE_HTML)

.PHONY: coverage.push
coverage.push:
	@echo "$(REPO) coverage.push"
	$(CODECOV) -f $(COVERAGE_FILE)$(if $(CODECOV_TOKEN), -t $(CODECOV_TOKEN),)
