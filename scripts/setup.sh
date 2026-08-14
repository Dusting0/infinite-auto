#!/usr/bin/env bash
# 一键配置 macOS 开发环境：Homebrew → go/node → air → 前端依赖。
# 用法：
#   bash scripts/setup.sh             安装/补齐依赖
#   bash scripts/setup.sh --dry-run   只检查不安装，报告缺失项（别名：--check）
set -euo pipefail

cd "$(dirname "$0")/.."

DRY_RUN=0
for arg in "$@"; do
  case "$arg" in
    --dry-run|--check) DRY_RUN=1 ;;
    -h|--help) echo "用法: bash scripts/setup.sh [--dry-run]"; exit 0 ;;
    *) echo "未知参数: $arg" >&2; exit 2 ;;
  esac
done

if [[ $DRY_RUN -eq 1 ]]; then
  echo "=== dry-run 模式：只检查，不安装任何内容 ==="
fi

# 把已存在但不在 PATH 里的 brew 注入当前 shell（仅设置环境变量，安全无副作用）
inject_brew() {
  if [[ -x /opt/homebrew/bin/brew ]]; then
    eval "$(/opt/homebrew/bin/brew shellenv)"
  elif [[ -x /usr/local/bin/brew ]]; then
    eval "$(/usr/local/bin/brew shellenv)"
  fi
}

# 报告某命令是否存在及版本（仅用于 dry-run 展示）
report_cmd() { # name  version-cmd...
  local name="$1"; shift
  if command -v "$name" >/dev/null 2>&1; then
    echo "  ✓ $name 已安装（$("$@" 2>/dev/null | head -1)）"
  else
    echo "  ✗ $name 未安装"
  fi
}

# ---------- 1. Homebrew ----------
echo "==> 1/4 Homebrew"
inject_brew
if command -v brew >/dev/null 2>&1; then
  echo "  ✓ brew 已安装（$(brew --version | head -1)）"
else
  echo "  ✗ brew 未安装"
  if [[ $DRY_RUN -eq 0 ]]; then
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    inject_brew
  else
    echo "    （dry-run：跳过安装）"
  fi
fi

# ---------- 2. go / node ----------
echo "==> 2/4 Brewfile 依赖（go、node）"
if [[ $DRY_RUN -eq 1 ]]; then
  report_cmd go go version
  report_cmd node node -v
  report_cmd npm npm -v
else
  brew bundle --file=scripts/Brewfile
fi

# ---------- 3. air ----------
echo "==> 3/4 air（Go 热重载）"
air_bin=""
if command -v go >/dev/null 2>&1; then
  air_bin="$(go env GOPATH)/bin/air"
fi
if [[ -n "$air_bin" && -x "$air_bin" ]]; then
  echo "  ✓ air 已安装（$("$air_bin" -v 2>/dev/null | tail -1)）"
elif command -v air >/dev/null 2>&1; then
  echo "  ✓ air 已安装（PATH）"
else
  echo "  ✗ air 未安装"
  if [[ $DRY_RUN -eq 0 ]]; then
    go install github.com/air-verse/air@latest
  else
    echo "    （dry-run：跳过安装）"
  fi
fi

# ---------- 4. 前端依赖 ----------
echo "==> 4/4 前端依赖（frontend/node_modules）"
if [[ -d frontend/node_modules ]]; then
  echo "  ✓ 已安装（$(ls frontend/node_modules | wc -l | tr -d ' ') 个包）"
else
  echo "  ✗ 未安装"
  if [[ $DRY_RUN -eq 0 ]]; then
    npm install --prefix frontend
  else
    echo "    （dry-run：跳过安装）"
  fi
fi

if [[ $DRY_RUN -eq 1 ]]; then
  echo
  echo "=== dry-run 检查完成。补齐缺失项请运行: bash scripts/setup.sh ==="
else
  cat <<EOF

✅ 开发环境就绪！
   启动开发：make dev      （前端 HMR + 后端 air 热重载）
   构建产物：make build    （前端 + macOS/Windows 二进制）
EOF
fi
