.PHONY: build
build:
	go build -o $(PWD)/build/server $(PWD)/cmd/server
	go build -o $(PWD)/build/client $(PWD)/cmd/client

.PHONY: run
run:
	docker compose -p faraway -f ./docker-compose.yaml up --build
