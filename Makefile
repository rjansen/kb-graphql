include Makefile.vars

.PHONY: clean
clean:
	@echo "$(ROOT_REPO)@$(BUILD) clean"
	-rm -Rf $(TMP_DIR)

$(TMP_DIR):
	mkdir -p $(TMP_DIR)

.PHONY: setup
setup: $(TMP_DIR)
	@echo "$(ROOT_REPO)@$(BUILD) setup"

.PHONY: install.gvm
install.gvm:
	@echo "$(ROOT_REPO)@$(BUILD) install.gvm"
	which gvm || \
		curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer | bash

.PHONY: setup.gvm
setup.gvm:
	@echo "$(ROOT_REPO)@$(BUILD) setup.gvm"
	gvm install go$(GOVERSION) -B
	@echo -e 'Please run:\n `gvm use go$(GOVERSION) --default`'

.PHONY: docker
docker.build:
	@echo "$(ROOT_REPO)@$(BUILD) docker"
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) \
		         -t $(DOCKER_NAME) -t $(DOCKER_NAME):$(VERSION) -f $(DOCKER_FILE) .

.PHONY: docker.bash
docker.bash:
	@echo "$(ROOT_REPO)@$(BUILD) docker.bash"
	docker run --rm --name $(NAME)-bash --entrypoint bash -it -u $(shell id -u):$(shell id -g) \
			   -v $(BASE_DIR):/go/src/$(ROOT_REPO) $(DOCKER_NAME)

docker.%:
	@echo "$(ROOT_REPO)@$(BUILD) docker.$*"
	docker run --rm --name $(NAME)-run -u $(shell id -u):$(shell id -g) \
    		    -v $(BASE_DIR):/go/src/$(ROOT_REPO) $(DOCKER_NAME) $*

include .make/*.makefile
