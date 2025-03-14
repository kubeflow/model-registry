detected_OS := $(shell uname -s)
real_OS := $(detected_OS)
arch := $(shell uname -m)
ifeq ($(detected_OS),Darwin)
	detected_OS := mac
	real_OS := darwin
endif
ifeq ($(detected_OS),Linux)
	detected_OS := linux
	real_OS := linux
endif


## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
TILT ?= $(LOCALBIN)/tilt
POETRY ?= $(LOCALBIN)/poetry

## Tool Versions
TILT_VERSION := 0.33.22

.PHONY: poetry
.PHONY: $(POETRY)
poetry: $(POETRY) ## Download poetry locally if necessary.
$(POETRY): $(LOCALBIN)
	test -s $(LOCALBIN)/poetry || curl -sSL https://install.python-poetry.org | POETRY_HOME=$(shell pwd) python3 -

.PHONY: tilt
.PHONY: $(TILT)
tilt: $(TILT) ## Download tilt locally if necessary. Architecture is locked at x86_64.
$(TILT): $(LOCALBIN)
	test -s $(LOCALBIN)/tilt || curl -fsSL https://github.com/tilt-dev/tilt/releases/download/v$(TILT_VERSION)/tilt.$(TILT_VERSION).$(detected_OS).$(arch).tar.gz | tar -xz -C $(LOCALBIN) tilt

tilt-up: poetry tilt
	$(LOCALBIN)/tilt up

tilt-down: tilt
	$(LOCALBIN)/tilt down
