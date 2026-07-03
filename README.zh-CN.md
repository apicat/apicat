# apicat

[English](README.md) | 简体中文

**Agent 时代的 API 文档基础设施。**

AI Agent 正在成为代码的主要编写者——同时也是主要消费者。当 Agent 对接一个 API 时，它需要结构化、机器可读的文档。而现有 API 文档工具全部为人设计（漂亮的 UI、"Try it" 按钮），spec 反倒成了附属品。

apicat 把这个关系倒过来：

> **OpenAPI spec 是一等公民，HTML 渲染只是它的一个视图。**

Agent 读原始 spec，人读同一份源渲染出的干净 HTML——因为源只有一份，两边永远一致。

## 特性

- **一条命令，零配置** — 指向包含 OpenAPI spec 的目录，即刻获得可浏览的文档
- **多文件 spec** — 跨文件 `$ref` 开箱即用（OpenAPI 3.0 / 3.1）
- **单个二进制** — Go 编译，模板与静态资源全部内嵌；无需 Node、Python 或任何运行时依赖
- **服务端渲染** — 快速、干净、低阅读疲劳的 HTML；渲染页面同时是 spec 的质量检查器（缺失的 description、不全的 example 一眼可见）

## 安装

```bash
go install github.com/apicat/apicat/v3/cmd/apicat-cli@latest
```

或从源码构建：

```bash
git clone https://github.com/apicat/apicat.git
cd apicat
go build -o apicat-cli ./cmd/apicat-cli
```

## 使用

```bash
apicat-cli path/to/spec-dir
```

目录中需包含 `openapi.yaml`、`openapi.yml` 或 `openapi.json` 入口文件，被引用的文件相对它解析。然后打开 http://127.0.0.1:8080。

参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-port` | `8080` | 监听端口 |
| `-host` | `127.0.0.1` | 绑定地址 |

## 路线图

- **本地** — spec 质量检查与 lint 提示、搜索
- **云端** — `apicat-cli publish` 发布到可自部署的服务端：文档共享、团队协作、spec 版本管理（`latest` + 固定 tag）
- **Agent 原生访问** — 原始 spec 端点、API 发现与搜索、MCP Server（原生对接 Claude Code / Cursor 等）

## 关于 v3 重写

v3 是彻底的重写。此前的 2.x 代码保留在 [`old`](https://github.com/apicat/apicat/tree/old) 分支，2.x 的各个发布版本仍可在 [tags](https://github.com/apicat/apicat/tags) 中获取。

## 许可证

[AGPL-3.0](LICENSE)
