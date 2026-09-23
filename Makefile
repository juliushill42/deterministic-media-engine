.PHONY: all test build smoke run
all: test build smoke
test:
	go test ./...
	go vet ./...
build:
	mkdir -p dist
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o dist/dme ./cmd/dme
smoke: build
	./scripts/smoke.sh
run: build
	./dist/dme serve
