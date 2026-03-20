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
	@if ! test -s $(LOCALBIN)/poetry; then \
		python3 -m venv $(shell pwd)/venv; \
		PYTHON_LIBDIR=$$(python3 -c "import sysconfig; print(sysconfig.get_config_var('LIBDIR') or '')"); \
		if [ -n "$$PYTHON_LIBDIR" ]; then \
			find "$$PYTHON_LIBDIR" -maxdepth 1 -name 'libpython*.so*' -exec ln -sf {} $(shell pwd)/venv/lib/ \; 2>/dev/null || true; \
		fi; \
		$(shell pwd)/venv/bin/pip install --quiet poetry; \
		ln -sf $(shell pwd)/venv/bin/poetry $(LOCALBIN)/poetry; \
	fi

.PHONY: tilt
.PHONY: $(TILT)
tilt: $(TILT) ## Download tilt locally if necessary. Architecture is locked at x86_64.
$(TILT): $(LOCALBIN)
	test -s $(LOCALBIN)/tilt || curl -fsSL https://github.com/tilt-dev/tilt/releases/download/v$(TILT_VERSION)/tilt.$(TILT_VERSION).$(detected_OS).$(arch).tar.gz | tar -xz -C $(LOCALBIN) tilt

tilt-up: poetry tilt
	$(LOCALBIN)/tilt up

tilt-down: tilt
	$(LOCALBIN)/tilt down
