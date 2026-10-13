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
- **搜索** — 按 method、path、summary、operationId、description 过滤接口
- **浅色 / 深色主题** — 默认跟随系统，可手动切换
- **热重载** — `-watch` 在 spec 文件变更时自动重载；改坏了会保留上一个可用版本
- **Lint 提示** — 缺失的 summary、description、operationId、错误响应会在总览页汇总计数，并在出现处以内联标记提示

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

apicat 是项目制的：spec 集中存放在 `~/.apicat/projects/` 下，每个项目一个目录，目录名由项目绝对路径编码而来——`[A-Za-z0-9_-]` 之外的字符全部替换为 `-`（如 `/Users/alice/dev/shop` → `-Users-alice-dev-shop`）。在项目根目录直接运行即可自动定位：

```bash
cd path/to/your-project
apicat-cli
```

也可以显式指定任意 spec 目录：

```bash
apicat-cli path/to/spec-dir
```

无论哪种方式，spec 目录中需包含 `openapi.yaml`、`openapi.yml` 或 `openapi.json` 入口文件，被引用的文件相对它解析。然后打开 http://127.0.0.1:8080。

参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-port` | `8080` | 监听端口 |
| `-host` | `127.0.0.1` | 绑定地址 |
| `-watch` | 关闭 | spec 文件变更时自动重载 |

## Agent Skill

[`skills/apicat-spec/`](skills/apicat-spec/SKILL.md) 是一个 [Agent Skill](https://code.claude.com/docs/en/skills)，教会 coding agent 按 apicat 的方式编写 OpenAPI spec——存放到 apicat 的项目目录、渲染出来零 lint 提示。在你的项目中使用：

```bash
mkdir -p .claude/skills
cp -r skills/apicat-spec .claude/skills/
```

之后让 agent 创建或更新 API spec，它会自动遵循这些约定：spec 的存放位置（`~/.apicat/projects/<编码路径>/`）、入口文件命名、多文件 `$ref` 组织、必填的 description、错误响应、example 等。Skill 文件是纯 Markdown，其他支持技能/指令文件的工具同样可用。

## 路线图

- **云端** — `apicat-cli publish` 发布到可自部署的服务端：文档共享、团队协作、spec 版本管理（`latest` + 固定 tag）
- **Agent 原生访问** — 原始 spec 端点、API 发现与搜索、MCP Server（原生对接 Claude Code / Cursor 等）

## 关于 v3 重写

v3 是彻底的重写。此前的 2.x 代码保留在 [`old`](https://github.com/apicat/apicat/tree/old) 分支，2.x 的各个发布版本仍可在 [tags](https://github.com/apicat/apicat/tags) 中获取。

## 许可证

[AGPL-3.0](LICENSE)
