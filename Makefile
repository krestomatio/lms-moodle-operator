PROJECT_SHORTNAME ?= lms-moodle
VERSION ?= 0.4.0
MOODLE_OPERATOR_VERSION ?= 0.6.12
POSTGRES_OPERATOR_VERSION ?= 0.3.7
NFS_OPERATOR_VERSION ?= 0.4.7
KEYDB_OPERATOR_VERSION ?= 0.3.7
OPERATOR_TYPE ?= go
PROJECT_TYPE ?= $(OPERATOR_TYPE)-operator
COMMUNITY_OPERATOR_NAME ?= lms-moodle-operator

include hack/mk/main.mk
