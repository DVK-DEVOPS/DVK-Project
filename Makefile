GO_DIR := go
BIN := app

ifneq ("$(wildcard $(GO_DIR)/.env)","")
	include $(GO_DIR)/.env
	export
endif

run:
	cd $(GO_DIR) && go run .

build:
	cd $(GO_DIR) && go build -o ../$(BIN) .

test:
	cd $(GO_DIR) && go test ./...

up:
	docker-compose -f docker-compose.dev.yml up --build

down:
	docker-compose -f docker-compose.dev.yml down

clean:
	rm -f $(BIN)
