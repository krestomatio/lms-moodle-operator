##@ Pipelines

MK_PIPELINES_PR_FILE ?= $(MK_INCLUDE_DIR)/pipelines-pr.mk
MK_PIPELINES_PR_SKIP_FILE ?= $(MK_INCLUDE_DIR)/pipelines-pr-skip.mk
MK_PIPELINES_RELEASE_FILE ?= $(MK_INCLUDE_DIR)/pipelines-release.mk
MK_PIPELINES_RELEASE_SKIP_FILE ?= $(MK_INCLUDE_DIR)/pipelines-release-skip.mk

## start if not SKIP_PIPELINE
ifeq ($(origin SKIP_PIPELINE),undefined)

## Pull request
include $(MK_PIPELINES_PR_FILE)

## Release
include $(MK_PIPELINES_RELEASE_FILE)

## else if not SKIP_PIPELINE
else
$(info SKIP_PIPELINE set:)

## Pull request
include $(MK_PIPELINES_PR_SKIP_FILE)

## Release
include $(MK_PIPELINES_RELEASE_SKIP_FILE)

## end if not SKIP_PIPELINE
endif
