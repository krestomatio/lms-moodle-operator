# Makefiles
MK_INCLUDE_DIR ?= hack/mk
MK_VARS_FILE ?= $(MK_INCLUDE_DIR)/vars.mk
MK_TARGET_FILE ?= $(MK_INCLUDE_DIR)/targets.mk
MK_DIST_FILE ?= Makefile-dist.mk
MK_PIPELINES_FILE ?= $(MK_INCLUDE_DIR)/pipelines.mk

include $(MK_VARS_FILE)
include $(MK_DIST_FILE)
include $(MK_TARGET_FILE)
include $(MK_PIPELINES_FILE)
