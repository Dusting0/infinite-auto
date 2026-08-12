# 无限规则计算器

无限流 TRPG 规则通用计算器（血量 / 防御 / 掷骰 / 三豁免 / 自定义伤害结算）。

## 运行

```bash
go run .
# 浏览器打开 http://localhost:8080
```

需要 Go 1.22+。仓库已提交前端构建产物，默认运行不需要额外 npm 命令。

## 前端开发

UI 使用 Vite + React + TypeScript + Tailwind CSS。开发 UI 时建议同时启动 Go API 和 Vite：

```bash
# 终端 1：Go API / 生产静态服务
go run .

# 终端 2：Vite 开发服务，/api 会代理到 :8080
npm install --prefix frontend
npm run dev --prefix frontend
```

更新前端源码后，构建产物会输出到 `web/static/dist`，供 Go embed 使用：

```bash
npm run build --prefix frontend
```

## 构建

```bash
make build OS=macos ARCH=amd64
make build OS=windows ARCH=amd64
make all
```

`make build` 会先执行前端构建，再编译 Go 二进制。

## 功能概要

- 血量：冲击 / 严重 / 恶性、长短休
- 防御预设：高速 / 破甲 / 破魔、完美防御、接触攻击、DR / ER / 吸收
- 掷骰：DP、加骰 8/9/10、附加成功、机运骰；结果汇总到右侧掷骰面板
- 三豁免：强韧 / 反射 / 意志快速检定
- 自定义伤害：击破结算，仅展示伤害不自动扣血
- Session 隔离：多浏览器互不干扰

## 产物

预编译二进制见项目协作目录中的 `产物/`（不纳入本仓库）。

## 版本

v0.1 · Powered by Dusting
