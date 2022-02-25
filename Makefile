# Copyright 2021 Linka Cloud  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MODULE = go.linka.cloud/totp


PROTO_BASE_PATH = $(PWD)

INCLUDE_PROTO_PATH = -I$(PROTO_BASE_PATH) \
	-I $(shell go list -m -f {{.Dir}} github.com/alta/protopatch)

PROTO_OPTS = paths=source_relative


$(shell mkdir -p .bin)

export GOBIN=$(PWD)/.bin

export PATH := $(GOBIN):$(PATH)

install:
	@go build -o /usr/local/bin/totp ./cmd/totp

bin:
	@go install github.com/alta/protopatch/cmd/protoc-gen-go-patch

clean:
	@rm -rf .bin
	@find $(PROTO_BASE_PATH) -name '*.pb*.go' -type f -exec rm {} \;


.PHONY: proto
proto: gen-proto lint

.PHONY: gen-proto
gen-proto: bin
	@find $(PROTO_BASE_PATH) -name '*.proto' -type f -exec \
    	protoc $(INCLUDE_PROTO_PATH) \
    		--go-patch_out=plugin=go,$(PROTO_OPTS):. {} \;

.PHONY: lint
lint:
	@goimports -w -local $(MODULE) $(PWD)
	@gofmt -w $(PWD)

.PHONY: tests
tests: proto
	@go test -v ./...
