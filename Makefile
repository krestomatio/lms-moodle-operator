PROJECT_SHORTNAME ?= kio
VERSION ?= 0.3.35
OPERATOR_TYPE ?= go
PROJECT_TYPE ?= $(OPERATOR_TYPE)-operator

include hack/mk/main.mk
