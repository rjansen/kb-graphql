NAME 		:= fivecolors-deploy
APP	    	?= asset
VERSION     ?= 0.0.1
PORT        ?= 8081
TARGET_PORT ?= $(PORT)
TYPE        ?= NodePort
MAKEFILE    := $(word $(words $(MAKEFILE_LIST)), $(MAKEFILE_LIST))
BASE_DIR    := $(shell cd $(dir $(MAKEFILE)); pwd)
TEMP_DIR    ?= $(BASE_DIR)/tmp
VAR_DIR     ?= $(BASE_DIR)/../../var


.PHONY: default
default: build

.PHONY: build
build: prepare
	@#go build -ldflags "-linkmode external -extldflags -static -w -s -X main.version=$(VERSION)-$(BUILD)" -o $(MAIN) $(REPO)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 LD_FLAGS='-linkmode external -extldflags -static' MAIN=$(APP) $(MAKE) -C ../../ docker.build
	cp ../../$(APP) $(TEMP_DIR)/$(APP)/

.PHONY: deploy
deploy: build
	docker build -t fivecolors-$(APP) -t fivecolors-$(APP):$(VERSION) -f Dockerfile.$(APP) $(TEMP_DIR)/$(APP)

.PHONY: prepare
prepare:
	-rm -Rf $(TEMP_DIR)/$(APP)
	mkdir -p $(TEMP_DIR)/$(APP)

# just store the command
expose:
	kubectl expose deploy fivecolors-$(APP) --name=fivecolors-$(APP) --port=$(PORT) --target-port=$(TARGET_PORT) --type=$(TYPE)

other:
	# i guess that this command overrides the image value to another
	kubectl create -f deployment.yaml --set spec.template.spec.containers[0].image=image_url

