CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.62.2
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PROJNAME=urler
PACKAGE=github.com/dbeleon/urler/urler/cmd/app


all: format build test lint

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

install-go-deps:
	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1
	GOBIN=$(BINDIR) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(BINDIR) go install github.com/envoyproxy/protoc-gen-validate@v1.1.0
	GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.25.0

get-go-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

vendor-proto:
		mkdir -p vendor-proto
		@if [ ! -d vendor-proto/google ]; then \
			mkdir -p vendor-proto/google &&\
			cd vendor-proto/google &&\
			git init &&\
			echo '[remote "origin"]' >> .git/config &&\
			echo 'url = https://github.com/googleapis/googleapis.git' >> .git/config &&\
			echo 'fetch = +refs/heads/master:refs/remotes/origin/master' >> .git/config &&\
			git config core.sparseCheckout true &&\
			echo 'google/api/' >> .git/info/sparse-checkout &&\
			git fetch --depth=1 &&\
			git checkout origin/master -B master &&\
			mkdir api &&\
			mv google/api/*.proto ./api &&\
			rm -rf .git &&\
			rm -rf google ;\
		fi
		@if [ ! -d vendor-proto/google/protobuf ]; then\
			echo '123' &&\
			mkdir -p vendor-proto/google/protobuf &&\
			cd vendor-proto/google/protobuf &&\
			git init &&\
			echo '[remote "origin"]' >> .git/config &&\
			echo 'url = https://github.com/protocolbuffers/protobuf.git' >> .git/config &&\
			echo 'fetch = +refs/heads/main:refs/remotes/origin/main' >> .git/config &&\
			git config core.sparseCheckout true &&\
			echo 'src/google/protobuf/' >> .git/info/sparse-checkout &&\
			git fetch --depth=1 &&\
			git checkout origin/main -B master &&\
			mv src/google/protobuf/*.proto ./ &&\
			rm -rf src &&\
			rm -rf .git ;\
		fi
		@if [ ! -d vendor-proto/validate ]; then \
			mkdir -p vendor-proto/validate &&\
			cd vendor-proto/validate &&\
			git init &&\
			echo '[remote "origin"]' >> .git/config &&\
			echo 'url = https://github.com/envoyproxy/protoc-gen-validate.git' >> .git/config &&\
			echo 'fetch = +refs/heads/main:refs/remotes/origin/main' >> .git/config &&\
			git config core.sparseCheckout true &&\
			echo 'validate/' >> .git/info/sparse-checkout &&\
			git fetch --depth=1 &&\
			git checkout origin/main -B master &&\
			mv validate/*.proto ./ &&\
			rm -rf validate &&\
			rm -rf .git ;\
		fi

generate:
	mkdir -p pkg/$(PROJNAME)/v1
	protoc -I api/$(PROJNAME)/v1 -I vendor-proto \
	--go_out=pkg/$(PROJNAME)/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/$(PROJNAME)/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out="lang=go:pkg/$(PROJNAME)/v1" --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	--grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
	--grpc-gateway_out pkg/$(PROJNAME)/v1 \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	api/$(PROJNAME)/v1/$(PROJNAME).proto
