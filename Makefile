OPERATOR_NAME ?= kio-operator
REPO_NAME ?= kio-operator
VERSION ?= 0.0.1

MK_PIPELINES_PR_FILE ?= $(MK_INCLUDE_DIR)/pipelines-pr-go.mk

include hack/mk/main.mk
