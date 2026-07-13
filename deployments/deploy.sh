#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${CLUSTER_NAME:-queueflow-cluster}"
REGION="${AWS_REGION:-us-east-1}"

echo "Configuring kubectl..."
aws eks update-kubeconfig --name "$CLUSTER_NAME" --region "$REGION"

echo "Applying namespace..."
kubectl apply -f deployments/k8s/namespace.yaml

echo "Deploying services..."
kubectl apply -f deployments/k8s/job-service.yaml
kubectl apply -f deployments/k8s/worker-service.yaml
kubectl apply -f deployments/k8s/scheduler-service.yaml

echo "Waiting for deployments..."
kubectl rollout status deployment/job-service -n queueflow --timeout=5m
kubectl rollout status deployment/worker-service -n queueflow --timeout=5m
kubectl rollout status deployment/scheduler-service -n queueflow --timeout=5m

echo "Deployment complete. Services are running in namespace 'queueflow'."
