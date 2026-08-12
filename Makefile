APP_NAME ?= infinite-calc
OUT_DIR ?= ./output
OS ?= macos
ARCH ?= $(shell go env GOHOSTARCH)

GOFLAGS ?= -ldflags="-s -w"
CGO_ENABLED ?= 0
NPM ?= npm

.PHONY: help frontend-install frontend-build build build-go macos windows all clean

help:
	@printf '%s\n' 'Usage:'
	@printf '%s\n' '  make frontend-build'
	@printf '%s\n' '  make build OS=macos [ARCH=amd64|arm64]'
	@printf '%s\n' '  make build OS=windows [ARCH=amd64|arm64]'
	@printf '%s\n' '  make macos'
	@printf '%s\n' '  make windows'
	@printf '%s\n' '  make all'
	@printf '%s\n' ''
	@printf '%s\n' 'Options:'
	@printf '%s\n' '  APP_NAME  Binary name prefix, default: infinite-calc'
	@printf '%s\n' '  OUT_DIR   Output directory, default: ./output'
	@printf '%s\n' '  OS        Target OS: macos or windows'
	@printf '%s\n' '  ARCH      Target architecture, default: host architecture'

frontend-install:
	@$(NPM) install --prefix frontend

frontend-build: frontend-install
	@$(NPM) run build --prefix frontend

build: frontend-build build-go

build-go:
	@mkdir -p '$(OUT_DIR)'
	@case '$(OS)' in \
		macos|darwin) \
			GOOS=darwin GOARCH='$(ARCH)' CGO_ENABLED='$(CGO_ENABLED)' go build $(GOFLAGS) -o '$(OUT_DIR)/$(APP_NAME)_macOS_$(ARCH)' . ;; \
		windows) \
			GOOS=windows GOARCH='$(ARCH)' CGO_ENABLED='$(CGO_ENABLED)' go build $(GOFLAGS) -o '$(OUT_DIR)/$(APP_NAME)_windows_$(ARCH).exe' . ;; \
		*) \
			printf '%s\n' 'Unsupported OS: $(OS). Use OS=macos or OS=windows.' >&2; \
			exit 2 ;; \
	esac

macos:
	@$(MAKE) build OS=macos ARCH='$(ARCH)'

windows:
	@$(MAKE) build OS=windows ARCH='$(ARCH)'

all: frontend-build
	@$(MAKE) build-go OS=macos ARCH='$(ARCH)'
	@$(MAKE) build-go OS=windows ARCH='$(ARCH)'

clean:
	@rm -f '$(OUT_DIR)/$(APP_NAME)_macOS_'* '$(OUT_DIR)/$(APP_NAME)_windows_'*.exe
