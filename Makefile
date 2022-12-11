.DEFAULT_GOAL := help
DOCKER_RUN_BUF_FLAGS = --rm --volume "$(shell pwd):/workspace" --workdir /workspace
# Buf CLI versions:
# https://hub.docker.com/r/bufbuild/buf/tags
DOCKER_BUF = bufbuild/buf:1.9.0

.PHONY: help
help : Makefile
	@sed -n 's/^##//p' $<

## buf-mod-update   : update your modules' dependencies to their latest versions
.PHONY: buf-mod-update
buf-mod-update:
	@for i in $(shell fd buf.yaml | xargs dirname) ; do \
	  docker run $(DOCKER_RUN_BUF_FLAGS) $(DOCKER_BUF) mod update $$i ; \
	done

## buf-gen          : generate your protos locally (won't need this if we're using remote generation)
.PHONY: buf-gen
buf-gen:
	docker run $(DOCKER_RUN_BUF_FLAGS) $(DOCKER_BUF) generate

## buf-lint         : lint and format protos
.PHONY: buf-lint
buf-lint:
	docker run $(DOCKER_RUN_BUF_FLAGS) $(DOCKER_BUF) lint
	docker run $(DOCKER_RUN_BUF_FLAGS) $(DOCKER_BUF) format -w
	# docker run --rm --volume "$(shell pwd):/workspace" --workdir /workspace bufbuild/buf breaking --against 'https://github.com/kevinmichaelchen/food-app.git#branch=main,subdir=backend/idl/proto'

## buf-login        : authenticate with the Buf Schema Registry (BSR)
.PHONY: buf-login
buf-login:
	@echo ${BUF_API_TOKEN} | buf registry login --username ${BUF_USER} --token-stdin

## buf-push         : push your protos to the Buf Schema Registry (BSR)
.PHONY: buf-push
buf-push:
	@for i in $(shell fd buf.yaml | xargs dirname) ; do \
	  echo $$i ; \
	  pushd . ; \
	  cd $$i ; \
	  pwd ; \
	  buf push ; \
	  popd ; \
	  echo "" ; \
	done