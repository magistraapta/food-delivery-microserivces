# Food Delivery Microservice

a microservice architecture for food delivery system built with Go.

## Architecture

![Architecture](./architecture-diagram.png)

## Features

- User Service

  - Register
  - Login
  - Get User

- Order Service

  - Create Order
  - Get Order

- Payment Service

  - Create Payment
  - Get Payment

- Food Service

  - Create Food
  - Get Food

- Delivery Service
  - Get delivery status
  - Update delivery status

## Service Structure

Each service is using layered architecture that contains repository, service, and controller layer.

## Tools

- RabbitMQ
- PostgreSQL
- Traefik
- Gorm

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (for local development)

### Quick Start with Makefile

This project includes a comprehensive Makefile for easy development. Run `make help` to see all available commands:

```bash
make help
```

#### 1. Clone and Setup

```bash
git clone <repository-url>
cd food-delivery
```

#### 2. Configure Environment

Copy `.env.example` to `.env` in each service directory and update values as needed:

```bash
cp services/user-service/.env.example services/user-service/.env
cp services/food-service/.env.example services/food-service/.env
cp services/order-service/.env.example services/order-service/.env
cp services/payment-service/.env.example services/payment-service/.env
```

#### 3. Create Databases

Start PostgreSQL and create the required databases:

```sql
CREATE DATABASE "food-service-db";
CREATE DATABASE "user-service-db";
CREATE DATABASE "order-service-db";
CREATE DATABASE "payment-service-db";
```

---

### Running with Docker Compose (Recommended)

The easiest way to run all services:

```bash
# Start all services (including RabbitMQ and Traefik)
make docker-up

# View logs
make logs

# Stop all services
make docker-down
```

**Access Points:**

- RabbitMQ Management: http://localhost:15672 (guest/guest)
- Traefik Dashboard: http://localhost:8080

---

### Running Locally (Development Mode)

For local development with hot-reload capability:

#### Option 1: Run All Services Together

```bash
# Start infrastructure (RabbitMQ, Traefik)
make infra

# Run all services concurrently
make run-all
```

Press `Ctrl+C` to stop all services.

#### Option 2: Run Services Individually

```bash
# Start infrastructure first
make infra

# Run individual services (each in separate terminal)
make run-user      # User Service    → http://localhost:8081
make run-food      # Food Service    → http://localhost:8082
make run-order     # Order Service   → http://localhost:8083
make run-payment   # Payment Service → http://localhost:8084
```

---

### Build & Test Commands

```bash
# Build all services
make build-all

# Build individual service
make build-user
make build-food
make build-order
make build-payment

# Run all tests
make test-all

# Test individual service
make test-user
make test-food
make test-order
make test-payment
```

---

### Utility Commands

```bash
# Install/update dependencies
make deps

# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Clean build artifacts
make clean
```

---

### Makefile Commands Reference

| Command             | Description                              |
| ------------------- | ---------------------------------------- |
| `make help`         | Show all available commands              |
| `make infra`        | Start infrastructure (RabbitMQ, Traefik) |
| `make infra-down`   | Stop infrastructure                      |
| `make run-all`      | Run all services concurrently            |
| `make run-user`     | Run user-service (port 8081)             |
| `make run-food`     | Run food-service (port 8082)             |
| `make run-order`    | Run order-service (port 8083)            |
| `make run-payment`  | Run payment-service (port 8084)          |
| `make build-all`    | Build all services                       |
| `make test-all`     | Test all services                        |
| `make docker-up`    | Start all services with Docker Compose   |
| `make docker-down`  | Stop all Docker services                 |
| `make docker-build` | Build all Docker images                  |
| `make logs`         | View Docker logs                         |
| `make deps`         | Install dependencies for all services    |
| `make fmt`          | Format code in all services              |
| `make lint`         | Run linter on all services               |
| `make clean`        | Clean build artifacts                    |

## API Endpoints

All endpoints are accessible via Traefik on port 80:

| Service         | Path Prefix    | Example                                        |
| --------------- | -------------- | ---------------------------------------------- |
| Food Service    | `/api/food`    | `GET http://localhost/api/food/restaurant/:id` |
| User Service    | `/api/user`    | `POST http://localhost/api/user/register`      |
| Order Service   | `/api/order`   | `POST http://localhost/api/order`              |
| Payment Service | `/api/payment` | `POST http://localhost/api/payment`            |

### Direct Access (Without Traefik)

| Service         | Port | Example                                       |
| --------------- | ---- | --------------------------------------------- |
| User Service    | 8081 | `GET http://localhost:8081/health`            |
| Food Service    | 8082 | `GET http://localhost:8082/restaurant/:id`    |
| Order Service   | 8083 | `GET http://localhost:8083/order/:id`         |
| Payment Service | 8084 | `GET http://localhost:8084/checkout/:orderId` |

## Payment Flow

The payment system uses **Stripe Checkout** for secure payment processing:

```
1. Client creates order → Order Service
2. Order Service publishes order.created event
3. Payment Service creates Stripe Checkout Session
4. Client gets checkout URL from Payment Service
5. User is redirected to Stripe Checkout page
6. User completes payment on Stripe
7. Stripe sends webhook to Payment Service
8. Payment Service publishes payment.success event
9. Order Service updates order status to CONFIRMED
```

### Payment Endpoints

| Endpoint                         | Method | Description                          |
| -------------------------------- | ------ | ------------------------------------ |
| `/api/payment/checkout/:orderId` | GET    | Get Stripe Checkout URL for an order |
| `/api/payment/status/:orderId`   | GET    | Get payment status for an order      |
| `/api/payment/webhook/stripe`    | POST   | Stripe webhook endpoint              |

### Payment Timeout

If payment is not completed within **5 minutes**, the order is automatically cancelled.

## Cons

- Every service need to validate the request from client

## To-do List

- [x] Return Stripe Payment URL
- [ ] Setup Unit test
- [ ] Setup Integration test
- [ ] Setup Kubernetes
- [ ] Setup CI/CD
- [ ] Setup S3 from minio

## Contact

email: magistaagis@gmail.com
twitter: @magistrapta
