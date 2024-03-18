PROJECT_SHORTNAME ?= kio
VERSION ?= 0.3.42
OPERATOR_TYPE ?= go
PROJECT_TYPE ?= $(OPERATOR_TYPE)-operator

include hack/mk/main.mk
