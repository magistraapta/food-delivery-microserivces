#!/usr/bin/env bash
# Apply Kubernetes manifests in the correct order
set -e

echo "Applying namespace..."
kubectl apply -f k8s/namespace.yml && \
echo "Applying database..."
kubectl apply -f k8s/database/ && \
echo "Applying rabbitmq..."
kubectl apply -f k8s/rabbitmq/ && \
kubectl apply -f k8s/shared/ && \
kubectl apply -f k8s/user/ && \
kubectl apply -f k8s/food/ && \
kubectl apply -f k8s/order/ && \
kubectl apply -f k8s/payment/ && \
kubectl apply -f k8s/ingress.yml