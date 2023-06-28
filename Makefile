VERSION ?= 0.0.1
BUF_VERSION:=1.1.0

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

KUSTOMIZE_VERSION ?= v3.8.7
CONTROLLER_TOOLS_VERSION ?= v0.9.2

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"

.PHONY: vendor
vendor:
	go mod vendor

mocks: vendor
	go run vendor/github.com/vektra/mockery/v2/main.go --all --dir ./internal --keeptree --case underscore

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s $(LOCALBIN); }

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./apis/..."

manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./apis/..." output:crd:artifacts:config=install/crd

build-job:
	docker build -t mayconjrpacheco/cloudx-runner:latest -f Dockerfile.runner .
	docker push mayconjrpacheco/cloudx-runner:latest

build-controller:
	docker build -t mayconjrpacheco/cloudx-controller:latest -f Dockerfile.controller .
	docker push mayconjrpacheco/cloudx-controller:latest

install-manifests:
	kubectl apply -f install/crd
	kubectl apply -f install/default
	kubectl apply -f install/rbac
	kubectl apply -f install/controller

install: manifests build-job build-controller
	kubectl apply -f install/crd
	kubectl apply -f install/default
	kubectl apply -f install/rbac
	kubectl apply -f install/controller

controller:
	go run cmd/controller/*.go

apiserver:
	go run cmd/apiserver/*.go
