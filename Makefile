SHELL := /bin/bash


VERSION := 1.0
KIND_CLUSTER := conduit-api-cluster
API_DOCKER_IMAGE := conduit-api-amd64

tidy:
	go mod tidy && go mod vendor

all: real-api

run:
	go run app/cli/conduit/main.go api serve

conduit-api:
	docker build \
		-f infra/docker/Dockerfile.conduit-api \
		-t $(API_DOCKER_IMAGE):$(VERSION) \
		--build-arg VCS_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.


# ==============================================================================
# Running from within k8s/kind
kind-up:
	kind create cluster \
		--image kindest/node:v1.22.0 \
		--name $(KIND_CLUSTER) \
		--config infra/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=conduit-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-restart:
	kubectl rollout restart deployment conduit-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-conduit:
	kubectl get pods -o wide --watch

kind-load:
	cd infra/k8s/kind/realworld-pod; kustomize edit set image conduit-api-image=$(API_DOCKER_IMAGE):$(VERSION)

	kind load docker-image $(API_DOCKER_IMAGE):$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build infra/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod

	kustomize build infra/k8s/kind/conduit-pod | kubectl apply -f -

kind-describe:
	kubectl describe pod -l app=conduit

kind-logs:
	kubectl logs -l app=conduit --all-containers=true -f --tail=100

kind-service-delete:
	kustomize build infra/k8s/kind/conduit-pod | kubectl delete -f -
	kustomize build infra/k8s/kind/database-pod | kubectl delete -f -
