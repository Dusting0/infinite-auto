# 无限规则计算器

无限流 TRPG 规则通用计算器（血量 / 防御 / 掷骰 / 三豁免 / 自定义伤害结算）。

## 环境准备（macOS 一键配置）

```bash
bash scripts/setup.sh
```

脚本会自动安装 Homebrew、Go、Node、air 及前端依赖。详见 `scripts/Brewfile`。

## 开发

```bash
make dev
# 浏览器打开 http://localhost:5173 （Vite HMR，/api 代理到 :8080）
```

前端改动即时热更新；后端 Go 改动由 air 自动重新编译重启。

## 运行（生产模式，直接跑嵌入产物）

```bash
go run .
# 浏览器打开 http://localhost:8080
```

需要 Go 1.22+。仓库已提交前端构建产物，默认运行不需要额外 npm 命令。

## 构建

```bash
make build    # 构建前端 + macOS/Windows 二进制到 output/
```

可选 `ARCH=amd64` 指定架构。

## 功能概要

- 血量：冲击 / 严重 / 恶性、长短休
- 防御预设：高速 / 破甲 / 破魔、完美防御、接触攻击、DR / ER / 吸收
- 掷骰：DP、加骰 8/9/10、附加成功、机运骰；结果汇总到右侧掷骰面板
- 三豁免：强韧 / 反射 / 意志快速检定
- 自定义伤害：击破结算，命中后自动回填到「受伤」供一键扣血
- Session 隔离：多浏览器互不干扰

## 产物

预编译二进制见项目协作目录中的 `产物/`（不纳入本仓库）。

## 版本

v0.1 · Powered by Dusting
