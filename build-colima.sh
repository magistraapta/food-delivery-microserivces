#!/usr/bin/env bash
# Build all service images using Colima's Docker daemon so Kubernetes pods can find them.
# Ensure Docker context is Colima: docker context use colima
# Run before: ./run-k8s.sh
set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

echo "Building images (using current Docker context; use 'docker context use colima' for Colima)..."
echo "Building user-service..."
docker build -t user-service:latest ./services/user-service

echo "Building food-service..."
docker build -t food-service:latest ./services/food-service

echo "Building order-service..."
docker build -t order-service:latest ./services/order-service

echo "Building payment-service..."
docker build -t payment-service:latest ./services/payment-service

echo ""
echo "All services built. Images:"
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | grep -E "REPOSITORY|user-service|food-service|order-service|payment-service"
echo ""
echo "Restart deployments to use new images: kubectl rollout restart deployment -n food-delivery"
