# Build stage - no BuildKit required (works with minikube legacy Docker)
ARG GO_VERSION=1.25.1
FROM golang:${GO_VERSION} AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/server .

# Final stage
FROM alpine:latest AS final

RUN apk --no-cache add ca-certificates tzdata && update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

COPY --from=build /bin/server /bin/
COPY .env.local.docker /app/.env.local.docker
WORKDIR /app

USER appuser

EXPOSE 8083
ENTRYPOINT [ "/bin/server" ]
