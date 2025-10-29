GO_DIR := go
GO_TAGS := fts5
BIN := app

ifneq ("$(wildcard $(GO_DIR)/.env)","")
	include $(GO_DIR)/.env
	export
endif

run:
	cd $(GO_DIR) && go run -tags "$(GO_TAGS)" .

build:
	cd $(GO_DIR) && go build -tags "$(GO_TAGS)" -o ../$(BIN) .

test:
	cd $(GO_DIR) && go test -tags "$(GO_TAGS)" ./...

up:
	docker-compose -f docker-compose.dev.yml up --build

down:
	docker-compose -f docker-compose.dev.yml down

clean:
	rm -f $(BIN)
