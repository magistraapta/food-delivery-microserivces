# Food Delivery Microservice

a microservice architecture for food delivery system built with Go.

## Architecture

![Architecture](./architecture-diagram.png)

## Services

- Api Gateway [8080]
- User Service [8081]
- Order Service [8082]
- Payment Service [8084]
- Food Service [8083]

## Tools

- RabbitMQ
- PostgreSQL
- Traefik
- Gorm

## Cons

- Every service need to validate the request from client
