.DEFAULT_GOAL := default

KIND_CLUSTER_NAME=opengamifylms

.PHONY: default
default: docker-compose-up ## Run default target

.PHONY: frontend
frontend: ## Start frontend
ifeq ($(shell which npm || echo false),false)
	@echo Error: npm not found in \$$PATH
	@exit 1
endif
	REACT_APP_CUSTOM_BACKEND="http://localhost:3000" npm run start --prefix ./frontend

.PHONY: docker-compose-up
docker-compose-up: ## Start application in Docker
ifeq ($(shell which docker || echo false),false)
	@echo Error: docker not found in \$$PATH
	@exit 1
endif
	cd backend && CGO_ENABLED=0 go build -ldflags "-s -w" -o ./bin/opengamifylms ./cmd/opengamifylms/main.go
	docker compose up --build

.PHONY: lint
lint: ## Run linter
ifeq ($(shell which docker || echo false),false)
	@echo Error: docker not found in \$$PATH
	@exit 1
endif
	docker run --rm -v $(PWD)/backend:/app -w /app golangci/golangci-lint:v1.58.1-alpine golangci-lint run ./...

.PHONY: test
test: ## Run tests
ifeq ($(shell which go || echo false),false)
	@echo Error: go not found in \$$PATH
	@exit 1
endif
	cd backend && go test -race -cover -count=1 ./...

.PHONY: test-cover
test-cover: ## Run tests and display coverage info in HTML viewer
ifeq ($(shell which go || echo false),false)
	@echo Error: go not found in \$$PATH
	@exit 1
endif
	t="/tmp/go-cover.\$\$\.tmp"; \
	go test -coverprofile=$$t ./backend/... && go tool cover -html=$$t && unlink $$t

.PHONY: test-e2e
test-e2e: ## Run Cypress e2e tests
ifeq ($(shell which docker || echo false),false)
	@echo Error: docker not found in \$$PATH
	@exit 1
endif
	docker compose -f e2e/docker-compose.yml -p opengamifylms_cypress down --remove-orphans --volumes
	docker compose -f e2e/docker-compose.yml -p opengamifylms_cypress up --abort-on-container-exit --exit-code-from cypress-e2e-tests --build
	docker compose -f e2e/docker-compose.yml -p opengamifylms_cypress down --remove-orphans --volumes

.PHONY: kind-up
kind-up: ## Creates a kind cluster if it doesn't exist
ifeq ($(shell which kind || echo false),false)
	@echo Error: kind not found in \$$PATH
	@exit 1
endif
	@if [ -z "`kind get clusters | grep $(KIND_CLUSTER_NAME)`" ]; then \
		kind create cluster --name $(KIND_CLUSTER_NAME); \
	else \
		echo "Kind cluster '$(KIND_CLUSTER_NAME)' already exists."; \
	fi

.PHONY: kind-down
kind-down: ## Deletes the kind cluster
ifeq ($(shell which kind || echo false),false)
	@echo Error: kind not found in \$$PATH
	@exit 1
endif
	@if [ ! -z "`kind get clusters | grep $(KIND_CLUSTER_NAME)`" ]; then \
		kind delete cluster --name $(KIND_CLUSTER_NAME); \
	else \
		echo "Kind cluster '$(KIND_CLUSTER_NAME)' does not exist."; \
	fi

.PHONY: helm-install
helm-install: kind-up ## Install Helm chart
ifeq ($(shell which helm || echo false),false)
	@echo Error: helm not found in \$$PATH
	@exit 1
endif
	helm install opengamifylms ./helm -n opengamifylms --create-namespace

.PHONY: helm-uninstall
helm-uninstall: ## Uninstall Helm chart
ifeq ($(shell which helm || echo false),false)
	@echo Error: helm not found in \$$PATH
	@exit 1
endif
	helm uninstall opengamifylms -n opengamifylms
	@echo "To delete the kind cluster, run 'make kind-down'"

.PHONY: sqlc
sqlc: ## Generate Go code from SQL queries using sqlc
ifeq ($(shell which docker || echo false),false)
	@echo Error: docker not found in \$$PATH
	@exit 1
endif
	docker run --rm -v $(PWD)/backend/store:/src -w /src sqlc/sqlc:1.26.0 generate

.PHONY: help
help: ## Makefile help
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
