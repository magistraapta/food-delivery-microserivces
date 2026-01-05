# Food Delivery Microservices - Makefile
# ======================================
# This Makefile provides commands to run all services locally

.PHONY: help run-all run-user run-food run-order run-payment \
        build-all build-user build-food build-order build-payment \
        test-all test-user test-food test-order test-payment \
        infra infra-down clean logs

# Colors for terminal output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
CYAN   := \033[0;36m
RESET  := \033[0m

# Service directories
SERVICES_DIR := ./services
USER_SERVICE := $(SERVICES_DIR)/user-service
FOOD_SERVICE := $(SERVICES_DIR)/food-service
ORDER_SERVICE := $(SERVICES_DIR)/order-service
PAYMENT_SERVICE := $(SERVICES_DIR)/payment-service

# Default target
help:
	@echo ""
	@echo "$(CYAN)╔══════════════════════════════════════════════════════════════╗$(RESET)"
	@echo "$(CYAN)║      Food Delivery Microservices - Development Commands      ║$(RESET)"
	@echo "$(CYAN)╚══════════════════════════════════════════════════════════════╝$(RESET)"
	@echo ""
	@echo "$(GREEN)Infrastructure:$(RESET)"
	@echo "  make infra          - Start infrastructure (RabbitMQ, Traefik)"
	@echo "  make infra-down     - Stop infrastructure"
	@echo ""
	@echo "$(GREEN)Run Services (requires infrastructure):$(RESET)"
	@echo "  make run-all        - Run all services concurrently"
	@echo "  make run-user       - Run user-service     (port 8081)"
	@echo "  make run-food       - Run food-service     (port 8082)"
	@echo "  make run-order      - Run order-service    (port 8083)"
	@echo "  make run-payment    - Run payment-service  (port 8084)"
	@echo ""
	@echo "$(GREEN)Build:$(RESET)"
	@echo "  make build-all      - Build all services"
	@echo "  make build-user     - Build user-service"
	@echo "  make build-food     - Build food-service"
	@echo "  make build-order    - Build order-service"
	@echo "  make build-payment  - Build payment-service"
	@echo ""
	@echo "$(GREEN)Test:$(RESET)"
	@echo "  make test-all       - Test all services"
	@echo "  make test-user      - Test user-service"
	@echo "  make test-food      - Test food-service"
	@echo "  make test-order     - Test order-service"
	@echo "  make test-payment   - Test payment-service"
	@echo ""
	@echo "$(GREEN)Docker:$(RESET)"
	@echo "  make docker-up      - Start all services with Docker Compose"
	@echo "  make docker-down    - Stop all Docker services"
	@echo "  make docker-build   - Build all Docker images"
	@echo "  make logs           - View Docker logs"
	@echo ""
	@echo "$(GREEN)Other:$(RESET)"
	@echo "  make clean          - Clean build artifacts"
	@echo ""

# ==========================================
# Infrastructure Commands
# ==========================================

# Start infrastructure services (RabbitMQ, Traefik) using Docker Compose
infra:
	@echo "$(CYAN)Starting infrastructure services...$(RESET)"
	docker-compose up -d rabbitmq traefik
	@echo "$(GREEN)Infrastructure started!$(RESET)"
	@echo "  - RabbitMQ Management: http://localhost:15672 (guest/guest)"
	@echo "  - Traefik Dashboard:   http://localhost:8080"

# Stop infrastructure services
infra-down:
	@echo "$(YELLOW)Stopping infrastructure services...$(RESET)"
	docker-compose down
	@echo "$(GREEN)Infrastructure stopped!$(RESET)"

# ==========================================
# Run Services Locally
# ==========================================

# Run all services concurrently (each in background)
# Note: This uses trap to handle Ctrl+C gracefully
run-all:
	@echo "$(CYAN)Starting all services...$(RESET)"
	@echo "$(YELLOW)Press Ctrl+C to stop all services$(RESET)"
	@echo ""
	@trap 'kill 0' SIGINT; \
	(cd $(USER_SERVICE) && go run main.go) & \
	(cd $(FOOD_SERVICE) && go run main.go) & \
	(cd $(ORDER_SERVICE) && go run main.go) & \
	(cd $(PAYMENT_SERVICE) && go run main.go) & \
	wait

# Run individual services
run-user:
	@echo "$(CYAN)Starting user-service on port 8081...$(RESET)"
	cd $(USER_SERVICE) && go run main.go

run-food:
	@echo "$(CYAN)Starting food-service on port 8082...$(RESET)"
	cd $(FOOD_SERVICE) && go run main.go

run-order:
	@echo "$(CYAN)Starting order-service on port 8083...$(RESET)"
	cd $(ORDER_SERVICE) && go run main.go

run-payment:
	@echo "$(CYAN)Starting payment-service on port 8084...$(RESET)"
	cd $(PAYMENT_SERVICE) && go run main.go

# ==========================================
# Build Commands
# ==========================================

# Build all services
build-all: build-user build-food build-order build-payment
	@echo "$(GREEN)All services built successfully!$(RESET)"

build-user:
	@echo "$(CYAN)Building user-service...$(RESET)"
	cd $(USER_SERVICE) && go build -o bin/user-service ./...

build-food:
	@echo "$(CYAN)Building food-service...$(RESET)"
	cd $(FOOD_SERVICE) && go build -o bin/food-service ./...

build-order:
	@echo "$(CYAN)Building order-service...$(RESET)"
	cd $(ORDER_SERVICE) && go build -o bin/order-service ./...

build-payment:
	@echo "$(CYAN)Building payment-service...$(RESET)"
	cd $(PAYMENT_SERVICE) && go build -o bin/payment-service ./...

# ==========================================
# Test Commands
# ==========================================

# Test all services
test-all: test-user test-food test-order test-payment
	@echo "$(GREEN)All tests completed!$(RESET)"

test-user:
	@echo "$(CYAN)Testing user-service...$(RESET)"
	cd $(USER_SERVICE) && go test -v ./...

test-food:
	@echo "$(CYAN)Testing food-service...$(RESET)"
	cd $(FOOD_SERVICE) && go test -v ./...

test-order:
	@echo "$(CYAN)Testing order-service...$(RESET)"
	cd $(ORDER_SERVICE) && go test -v ./...

test-payment:
	@echo "$(CYAN)Testing payment-service...$(RESET)"
	cd $(PAYMENT_SERVICE) && go test -v ./...

# ==========================================
# Docker Commands
# ==========================================

# Start all services with Docker Compose
docker-up:
	@echo "$(CYAN)Starting all services with Docker Compose...$(RESET)"
	docker-compose up -d
	@echo "$(GREEN)All services started!$(RESET)"

# Stop all Docker services
docker-down:
	@echo "$(YELLOW)Stopping all Docker services...$(RESET)"
	docker-compose down
	@echo "$(GREEN)All services stopped!$(RESET)"

# Build all Docker images
docker-build:
	@echo "$(CYAN)Building all Docker images...$(RESET)"
	docker-compose build
	@echo "$(GREEN)All images built!$(RESET)"

# View Docker logs
logs:
	docker-compose logs -f

# ==========================================
# Utility Commands
# ==========================================

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(RESET)"
	rm -rf $(USER_SERVICE)/bin
	rm -rf $(FOOD_SERVICE)/bin
	rm -rf $(ORDER_SERVICE)/bin
	rm -rf $(PAYMENT_SERVICE)/bin
	@echo "$(GREEN)Clean complete!$(RESET)"

# Install dependencies for all services
deps:
	@echo "$(CYAN)Installing dependencies for all services...$(RESET)"
	cd $(USER_SERVICE) && go mod tidy
	cd $(FOOD_SERVICE) && go mod tidy
	cd $(ORDER_SERVICE) && go mod tidy
	cd $(PAYMENT_SERVICE) && go mod tidy
	@echo "$(GREEN)Dependencies installed!$(RESET)"

# Format code in all services
fmt:
	@echo "$(CYAN)Formatting code in all services...$(RESET)"
	cd $(USER_SERVICE) && go fmt ./...
	cd $(FOOD_SERVICE) && go fmt ./...
	cd $(ORDER_SERVICE) && go fmt ./...
	cd $(PAYMENT_SERVICE) && go fmt ./...
	@echo "$(GREEN)Code formatted!$(RESET)"

# Run linter on all services
lint:
	@echo "$(CYAN)Linting all services...$(RESET)"
	cd $(USER_SERVICE) && golangci-lint run ./...
	cd $(FOOD_SERVICE) && golangci-lint run ./...
	cd $(ORDER_SERVICE) && golangci-lint run ./...
	cd $(PAYMENT_SERVICE) && golangci-lint run ./...
	@echo "$(GREEN)Linting complete!$(RESET)"
