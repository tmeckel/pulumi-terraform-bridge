LDFLAGS := -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore"

build: bin
	go build \
		-o bin/pulumi-resource-terraform-provider \
		${LDFLAGS} \
		github.com/pulumi/pulumi-terraform-bridge/dynamic

install: build
	pulumi plugin rm resource terraform-provider --all --yes
	pulumi plugin install resource terraform-provider v0.0.0-dev \
		--file bin/pulumi-resource-terraform-provider --reinstall

test_unit:
	cd internal/shim && go test ${LDFLAGS} ./...
	go test -short ${LDFLAGS} ./...

test:
	cd internal/shim && go test -v ${LDFLAGS} ./...
	go test -v ${LDFLAGS} ./...

test_accept:
	go test -v ${LDFLAGS} ./... -update

bin:
	mkdir bin

.PHONY: test
