.PHONY: lint
lint: go-lint ## Run linting tasks
	@echo "+ $@"

.PHONY: k8s
k8s: pr ## Run k8s tasks
	@echo "+ $@"

.PHONY: pr
pr: IMG = $(BUILD_IMAGE_TAG_BASE):$(BUILD_VERSION)
pr: image-build image-push deploy ## Run pr tasks
	@echo "+ $@"
