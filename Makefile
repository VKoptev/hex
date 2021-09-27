default: test

.PHONY: test
test: GO_TEST_FLAGS := -race
test:
	go test -mod=vendor $(GO_TEST_FLAGS) --tags=$(GO_TEST_TAGS) ./...

.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml
	prototool lint $(PROTOTOOL_FLAGS)

.PHONY: generate
generate:
	prototool generate $(PROTOTOOL_FLAGS)
	go generate -mod=vendor -x ./...

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor


.PHONY: dev-test
dev-test: GO_TEST_FLAGS := -race
dev-test:
	$(call gotools, make test)

.PHONY: dev-lint
dev-lint:
	$(call gotools, make lint)

.PHONY: dev-generate
dev-generate:
	$(call gotools, make generate)


.PHONY: dev-vendor
dev-vendor:
	$(call gotools, make vendor)
	git add vendor

define gotools
	docker run --rm \
		-v $(PWD):/work \
		-w /work \
		dockerhub.jeshik.ru/library/gotools:0.0.4 $1
endef
