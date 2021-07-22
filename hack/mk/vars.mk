CONTAINER_BUILDER ?= docker
OPERATOR_NAME ?= operator
REPO_NAME ?= operator
REPO_OWNER ?= krestomatio
VERSION ?= 0.0.1

# Image
REGISTRY ?= quay.io
REGISTRY_PATH ?= $(REGISTRY)/$(REPO_OWNER)
IMAGE_TAG_BASE ?= $(REGISTRY_PATH)/$(OPERATOR_NAME)
IMG ?= $(IMAGE_TAG_BASE):$(VERSION)

# requirements
OPERATOR_VERSION ?= 1.11.0
KUSTOMIZE_VERSION ?= 4.1.3
OPM_VERSION ?= 1.15.1

# JX
JOB_NAME ?= pr
PULL_NUMBER ?= 0
BUILD_ID ?= 0

# Build
BUILD_REGISTRY_PATH ?= docker-registry.jx.krestomat.io/krestomatio
BUILD_OPERATOR_NAME ?= $(OPERATOR_NAME)
BUILD_IMAGE_TAG_BASE ?= $(BUILD_REGISTRY_PATH)/$(BUILD_OPERATOR_NAME)
ifeq ($(JOB_NAME),release)
BUILD_VERSION ?= $(shell git rev-parse HEAD^2 &>/dev/null && git rev-parse HEAD^2 || echo)
else
BUILD_VERSION ?= $(shell git rev-parse HEAD 2> /dev/null  || echo)
endif

# CI
SKIP_MSG := skip.ci
RUN_PIPELINE ?= $(shell git log -1 --pretty=%B | cat | grep -q "\[$(SKIP_MSG)\]" && echo || echo 1)
ifeq ($(RUN_PIPELINE),)
SKIP_PIPELINE = true
$(info RUN_PIPELINE not set, skipping...)
endif
ifeq ($(BUILD_VERSION),)
SKIP_PIPELINE = true
$(info BUILD_VERSION not set, skipping...)
endif
ifeq ($(origin PULL_BASE_SHA),undefined)
CHANGELOG_FROM ?= HEAD~1
else
CHANGELOG_FROM ?= $(PULL_BASE_SHA)
endif

# molecule
MOLECULE_SEQUENCE ?= test
MOLECULE_SCENARIO ?= default
export OPERATOR_IMAGE ?= $(IMG)
export TEST_OPERATOR_NAMESPACE ?= $(OPERATOR_NAME)-$(JOB_NAME)-$(PULL_NUMBER)-$(BUILD_ID)

# skopeo
SKOPEO_SRC_TLS ?= True
SKOPEO_DEST_TLS ?= true

# Release
GIT_REMOTE ?= origin
GIT_BRANCH ?= master
GIT_ADD_FILES ?= Makefile
CHANGELOG_FILE ?= CHANGELOG.md
