.PHONY:server
server:
	@echo "Running server"
	@go run ./cmd/server/main.go --race

.PHONY:client
client:
	@echo "Running client"
	@go run ./cmd/client/main.go --race

.PHONY:ddos
ddos:
	@echo "Running client in DDOS mode"
	@sh ./scripts/ddos.sh

.PHONY:docker-server
build-docker-server:
	@echo "Building server docker image"
	@docker build -f ./deployments/docker/server/Dockerfile -t hashcash-server .

.PHONY:docker-client
build-docker-client:
	@echo "Building client docker image"
	@docker build -f ./deployments/docker/client/Dockerfile -t hashcash-client .

.PHONY:test-locally
test-locally:
	@echo "Running containers"
	@docker-compose -f docker-compose.local.yml up

.PHONY:test
test:
	@echo "Running containers"
	@docker-compose up
