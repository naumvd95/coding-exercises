GIT_HOST = github.com
PWD := $(shell pwd)
GOPATH_DEFAULT := $(shell go env GOPATH)
export GOPATH ?= $(GOPATH_DEFAULT)
GOBIN_DEFAULT := $(GOPATH)/bin
export GOBIN ?= $(GOBIN_DEFAULT)
PATH := $(PATH):$(PWD)/bin:$(GOBIN)
VERBOSE :=
ifndef VERBOSE
GOFLAGS   :=
DEPFLAGS  :=
else
GOFLAGS   := "-v"
DEPFLAGS  := "-v"
endif
GOOS ?= $(shell go env GOOS)

$(GOBIN):
	echo "create gobin"
	mkdir -p $(GOBIN)
work: $(GOBIN)

#go tools
HAS_DEP := $(shell command -v dep;)
depend: work
ifndef HAS_DEP
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
	dep ensure $(DEPFLAGS)

depend-update:
	dep ensure -update $(DEPFLAGS)

# binaries
machine:
	GOOS=$(GOOS) go build $(GOFLAGS) \
	-o bin/run_machine \
	machine.go && \
	chmod +x bin/run_machine
