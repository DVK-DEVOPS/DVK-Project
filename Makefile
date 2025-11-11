GO_DIR := go
BIN := app
LOCAL_DB_PORT := 55432
SSH_TUNNEL_PID_FILE := /tmp/ssh_tunnel.pid
SSH_KEY := ~/.ssh/id_rsa

ifneq ("$(wildcard $(GO_DIR)/.env)","")
	include $(GO_DIR)/.env
	export
endif

ssh-tunnel:
	if ! lsof -i:$(LOCAL_DB_PORT) >/dev/null 2>&1; then \
		echo "Starting SSH tunnel..."; \
		sh -c 'ssh -f -N -L $(LOCAL_DB_PORT):10.0.0.5:5432 adminuser@$(APP_VM_PUBLIC_IP) -i $(SSH_KEY) && echo $$! > $(SSH_TUNNEL_PID_FILE)'; \
	else \
		echo "SSH tunnel already running on port $(LOCAL_DB_PORT)"; \
	fi


ssh-tunnel-stop:
	@if lsof -ti:$(LOCAL_DB_PORT) >/dev/null 2>&1; then \
		echo "Killing SSH tunnel on port $(LOCAL_DB_PORT)..."; \
		pkill -f "ssh -f -N -L $(LOCAL_DB_PORT):localhost:5432"; \
	else \
		echo "No SSH tunnel running"; \
	fi

run: ssh-tunnel
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
