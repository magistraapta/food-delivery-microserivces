#!/usr/bin/env bash
# Delete manifests in reverse order of application
set -e

echo "Deleting services and deployments..."
[ -d k8s/rabbitmq ] && kubectl delete -f k8s/rabbitmq/ --ignore-not-found
[ -d k8s/database ] && kubectl delete -f k8s/database/ --ignore-not-found
kubectl delete -f k8s/user/ --ignore-not-found
kubectl delete -f k8s/payment/ --ignore-not-found
kubectl delete -f k8s/order/ --ignore-not-found
kubectl delete -f k8s/food/ --ignore-not-found

echo "Deleting ingress..."
kubectl delete -f k8s/ingress.yml --ignore-not-found

echo "Deleting shared config and secrets..."
kubectl delete -f k8s/shared/ --ignore-not-found

echo "Deleting namespace..."
kubectl delete -f k8s/namespace.yml --ignore-not-found

echo "Done."
