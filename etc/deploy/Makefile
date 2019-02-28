NAME 		:= fivecolors-deploy
APP	    	?= etl
MAKEFILE    := $(word $(words $(MAKEFILE_LIST)), $(MAKEFILE_LIST))
BASE_DIR    := $(shell cd $(dir $(MAKEFILE)); pwd)
TEMP_DIR    ?= $(BASE_DIR)/tmp
VAR_DIR     ?= $(BASE_DIR)/../../var

.PHONY: deploy
deploy: prepare
	@#go build -ldflags "-linkmode external -extldflags -static -w -s -X main.version=$(VERSION)-$(BUILD)" -o $(MAIN) $(REPO)
	GOOS=linux GOARCH=amd64 LD_FLAGS=" -linkmode external -extldflags -static" MAIN=$(APP) $(MAKE) -C ../../ docker.build
	cp ../../$(APP) $(TEMP_DIR)/$(APP)/
	@#ln -s $(PRJ)/Elysion/FiveColors/storage/db/assets/data $(TEMP_DIR)/$(APP)/data
	cp -R $(VAR_DIR)/$(APP)/data $(TEMP_DIR)/$(APP)/

.PHONY: prepare
prepare:
	-rm -Rf $(TEMP_DIR)/$(APP)
	mkdir -p $(TEMP_DIR)/$(APP)
