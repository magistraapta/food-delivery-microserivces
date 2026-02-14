#!/usr/bin/env bash
# Apply Kubernetes manifests in the correct order.
# Prerequisite: NGINX Ingress Controller must be installed (e.g. in ingress-nginx namespace).
#   helm install ingress-nginx ingress-nginx/ingress-nginx -n ingress-nginx --create-namespace
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
kubectl apply -f k8s/ingress.yml && \
kubectl apply -f k8s/monitoring/