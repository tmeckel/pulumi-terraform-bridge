PROJECT_NAME := Pulumi Terraform Bridge
include build/common.mk

PACK             := terraform
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-terraform-bridge
TESTPARALLELISM  := 10

build::
	go build ${PROJECT}/v2/pkg/...
	go build ${PROJECT}/v2/internal/...
	go install ${PROJECT}/v2/cmd/...

lint::
	golangci-lint run

test_fast::
	$(GO_TEST) ./...

test_all:: test_fast

.PHONY: check_clean_worktree
check_clean_worktree:
	$$(go env GOPATH)/src/github.com/pulumi/scripts/ci/check-worktree-is-clean.sh

# The travis_* targets are entrypoints for CI.
.PHONY: travis_cron travis_push travis_pull_request travis_api
travis_cron: all
travis_push: all check_clean_worktree only_test
travis_pull_request: all
travis_api: all
