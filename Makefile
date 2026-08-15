APP_NAME ?= infinite-calc
OUT_DIR ?= ./output
MACOS_ARCH ?= $(shell go env GOHOSTARCH)
WINDOWS_ARCH ?= amd64
NPM ?= npm
AIR ?= $(shell go env GOPATH)/bin/air

.PHONY: help dev build clean

help:
	@printf '%s\n' 'Usage:'
	@printf '%s\n' '  make dev      开发模式：前端 HMR + 后端 air 热重载'
	@printf '%s\n' '  make build    构建发布产物（前端 + macOS/Windows 二进制到 output/）'
	@printf '%s\n' '  make clean    清理产物与临时文件'
	@printf '%s\n' ''
	@printf '%s\n' 'Options:'
	@printf '%s\n' '  APP_NAME       二进制名前缀，默认: infinite-calc'
	@printf '%s\n' '  OUT_DIR        产物目录，默认: ./output'
	@printf '%s\n' '  MACOS_ARCH     macOS 目标架构，默认: 本机架构'
	@printf '%s\n' '  WINDOWS_ARCH   Windows 目标架构，默认: amd64 (Intel/AMD)'

# 开发：前端 vite HMR (:5173) + 后端 air 热重载 (:8080)，API 经 vite 代理。
dev:
	@echo "开发模式：前端 http://localhost:5173 (HMR)，后端 air 热重载 :8080"
	@trap 'kill 0' EXIT INT TERM; \
	$(NPM) run dev --prefix frontend & \
	$(AIR)

# 发布：构建前端 + macOS/Windows 二进制到 output/。
# Windows 默认编 amd64（Intel/AMD 可直接运行，覆盖默认的 arm64 本机架构）；
# macOS 默认取本机架构。两者均可用 make build MACOS_ARCH=.. WINDOWS_ARCH=.. 覆盖。
build:
	@$(NPM) install --prefix frontend
	@$(NPM) run build --prefix frontend
	@mkdir -p '$(OUT_DIR)'
	@GOOS=darwin  GOARCH='$(MACOS_ARCH)'   CGO_ENABLED=0 go build -ldflags="-s -w" -o '$(OUT_DIR)/$(APP_NAME)_macOS_$(MACOS_ARCH)' .
	@GOOS=windows GOARCH='$(WINDOWS_ARCH)' CGO_ENABLED=0 go build -ldflags="-s -w" -o '$(OUT_DIR)/$(APP_NAME)_windows_$(WINDOWS_ARCH).exe' .

clean:
	@rm -rf '$(OUT_DIR)' .air
