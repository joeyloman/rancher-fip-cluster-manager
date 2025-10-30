# Image URL to use all building/pushing image targets
IMG ?= your-repo/rancher-fip-cluster-manager:latest
# Produce CRDs that work back to Kubernetes 1.11 (no pruning).
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

all: manager

# =================================================================================================
# Development
# =================================================================================================

## Run manager binary against the cluster specified in ~/.kube/config
run: generate
	go run ./cmd/cluster-manager/main.go --leader-elect=true

## Run tests
test: generate
	go test -v ./pkg/... ./cmd/... ./internal/...

# =================================================================================================
# Build
# =================================================================================================

## Build manager binary
manager: generate
	go build -o bin/rancher-fip-cluster-manager cmd/cluster-manager/main.go

## Build the docker image
docker-build: test
	docker build -f Dockerfile -t ${IMG} .

## Push the docker image
docker-push:
	docker push ${IMG}

# =================================================================================================
# Deployment
# =================================================================================================

## Deploy controller to the cluster
deploy:
	kubectl apply -f config/deployment/deployment.yaml

## Undeploy controller from the cluster
undeploy:
	kubectl delete -f config/deployment/deployment.yaml

.PHONY: all run test manager docker-build docker-push generate deploy undeploy
