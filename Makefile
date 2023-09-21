#====================
AUTHOR         ?= The sacloud/sacloud-router-usage Authors
COPYRIGHT_YEAR ?= 2022-2023

BIN            ?= sacloud-router-usage
GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
include includes/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
