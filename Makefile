# Image URL to use all building/pushing image targets
IMG ?= timgcpsm-operator:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./... -coverprofile cover.out

##@ Build

.PHONY: manifests
manifests: ## Generate CRDs from Go types
	@echo "📝 Generating CRDs..."
	@bash scripts/generate-crds.sh

.PHONY: generate
generate: ## Generate deepcopy code
	@echo "📝 Generating deepcopy code..."
	@controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: generate-all
generate-all: generate manifests ## Generate all code and manifests

.PHONY: build
build: fmt vet ## Build manager binary.
	go build -o bin/manager cmd/main.go

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run cmd/main.go

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

##@ Deployment

.PHONY: install
install: ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	kubectl apply -f config/crd/timgcpsmsecret-crd.yaml
	kubectl apply -f config/crd/timgcpsmsecretconfig-crd.yaml
	kubectl apply -f config/crd/timgcpsmclusterconfig-crd.yaml

.PHONY: uninstall
uninstall: ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	kubectl delete -f config/crd/timgcpsmsecret-crd.yaml
	kubectl delete -f config/crd/timgcpsmsecretconfig-crd.yaml
	kubectl delete -f config/crd/timgcpsmclusterconfig-crd.yaml

.PHONY: deploy
deploy: ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	kubectl apply -f config/manager/namespace.yaml
	kubectl apply -f config/rbac/role.yaml
	kubectl apply -f config/rbac/service_account.yaml
	kubectl apply -f config/rbac/role_binding.yaml
	kubectl apply -f config/manager/deployment.yaml

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	kubectl delete -f config/manager/deployment.yaml
	kubectl delete -f config/rbac/role_binding.yaml
	kubectl delete -f config/rbac/role.yaml
	kubectl delete -f config/rbac/service_account.yaml

