#!/usr/bin/env bash
# Build all service images inside minikube's Docker daemon so Kubernetes pods can find them.
# Run after: minikube start
# Run before: ./run-k8s.sh
set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

if ! command -v minikube &>/dev/null; then
  echo "minikube is not installed or not in PATH"
  exit 1
fi

if ! minikube status &>/dev/null; then
  echo "minikube is not running. Start it with: minikube start"
  exit 1
fi

echo "Using minikube's Docker daemon..."
eval "$(minikube docker-env)"

echo "Building user-service..."
docker build -t user-service:latest ./services/user-service

echo "Building food-service..."
docker build -t food-service:latest ./services/food-service

echo "Building order-service..."
docker build -t order-service:latest ./services/order-service

echo "Building payment-service..."
docker build -t payment-service:latest ./services/payment-service

echo ""
echo "All services built. Images in minikube:"
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | grep -E "REPOSITORY|user-service|food-service|order-service|payment-service"
echo ""
echo "Apply manifests with: ./run-k8s.sh"
