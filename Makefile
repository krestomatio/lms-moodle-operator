PROJECT_SHORTNAME ?= kio
VERSION ?= 0.3.11
OPERATOR_TYPE ?= go
PROJECT_TYPE ?= $(OPERATOR_TYPE)-operator

include hack/mk/main.mk
