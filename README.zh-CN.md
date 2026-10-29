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
- **单个二进制** — Go 编译，React 构建产物全部内嵌；运行时无需 Node 或 Python
- **React 工作台** — React + TypeScript + Vite + shadcn/ui + Tailwind CSS，后端使用 Echo v5；包含概览、接口文档、公共组件、环境管理和 HTTP 调试
- **公共定义** — 从侧栏或 `/components` 独立查看 Schema、参数、请求体、响应及其他 `components` 定义
- **搜索** — 按 method、path、summary、operationId、description 过滤接口，也可按分类、名称、说明查找公共定义
- **可折叠导航** — 接口 tag 和公共定义分类默认收起；搜索时展开匹配分类，清空搜索后恢复原状态
- **浅色 / 深色主题** — 默认跟随系统，可手动切换
- **热重载** — `-watch` 在 spec 文件变更时自动重载；改坏了会保留上一个可用版本
- **Lint 提示** — 缺失的 summary、description、operationId、错误响应会在总览页汇总计数，并在出现处以内联标记提示

## 安装

```bash
go install github.com/apicat/apicat/v3/cmd/apicat-cli@latest
```

或从源码构建（Go 1.25+、Node 24+ / npm）：

```bash
git clone https://github.com/apicat/apicat.git
cd apicat
npm --prefix web ci
npm --prefix web run build
go build -o apicat-cli ./cmd/apicat-cli
```

## 使用

apicat 是项目制的：默认情况下，spec 集中存放在 `~/.apicat/projects/` 下，每个项目一个目录，目录名由项目绝对路径编码而来——`[A-Za-z0-9_-]` 之外的字符全部替换为 `-`（如 `/Users/alice/dev/shop` → `-Users-alice-dev-shop`）。在项目根目录直接运行即可自动定位：

```bash
cd path/to/your-project
apicat-cli
```

要配置项目的 spec 目录，在项目根目录创建 `.apicat/setting.json`：

```json
{
  "spec_dir": "docs/api"
}
```

支持相对路径和绝对路径，相对路径以项目根目录为基准。在该根目录运行 `apicat-cli`、`apicat-cli -watch` 或 `apicat-cli check` 即可使用配置。优先读取 `<当前项目>/.apicat/setting.json`，只有该文件不存在时，才读取 `~/.apicat/projects/<项目路径编码>/setting.json`。两处配置中的相对路径都以当前项目根目录为基准，不会向上查找。路径按字面值解析，不展开 `~` 或环境变量。

两处配置文件都不存在，或选中配置中的 `spec_dir` 未设置、为空字符串、仅含空白时，仍使用原来的 `~/.apicat/projects/<编码路径>/` 默认目录。本地配置一旦存在就完整覆盖兜底配置。JSON 格式错误、未知字段、字段类型错误或无法读取配置时会报错，不会继续尝试其他配置。配置只选择目录，不会创建目录或迁移已有文档。修改配置后需重启 CLI；`-watch` 只重载 spec。

显式指定目录的优先级最高，会跳过配置读取：

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

## 本地工作台

- **环境管理**（`/environments`）：创建开发、测试、生产环境，配置基础地址及变量。公共变量适用于所有环境，同名且启用的环境变量覆盖公共变量。在地址、参数值、请求头或请求体中使用 `{{变量名}}`。变量按字面值替换，不执行脚本。
- **HTTP 调试**：选择接口并打开“接口调试”，修改路径参数、查询参数、请求头、认证和原始请求体后发送。请求由 Go 发出，目标 API 无需为浏览器配置 CORS。支持 HTTP(S)、Bearer 请求头、JSON/文本/URL 编码表单、取消请求、超时设置以及响应头和正文查看，二进制响应以 Base64 展示。不自动重定向或主动重试，不共享 Cookie jar；响应上限为 2 MiB，超限时明确标记截断。
- **参数处理**：继承路径级参数及接口/路径级服务器地址。常用查询数组和对象支持 form、spaceDelimited、pipeDelimited 和浅层 deepObject。复杂路径编码和 multipart 文件上传暂不自动处理，遇到不支持的情况会提示；可将已编码参数直接填写到地址中。目前提供交互调试，不是带断言的自动化测试运行器。
- **数据保存**：环境数据保存在 `<spec目录>/.apicat/workspace.json`，使用仅当前用户可读写的文件权限、原子替换，并检测同一个运行中工作台的保存版本冲突。“隐藏值”仅遮罩界面，不加密文件；含密钥时应排除在版本控制之外。请求草稿和响应仅保留在页面内存中，不持久化为历史记录。
- **本地访问**：API 校验 Host/Origin，配置写入和请求执行要求会话令牌。默认仍监听 `127.0.0.1`；当前是本地工作台，不是带多用户认证的服务端。

### 前端开发

先运行 `go run ./cmd/apicat-cli -watch /path/to/spec-dir`，再运行 `npm --prefix web run dev`。Vite 提供开发界面，并将 `/api` 代理到 `127.0.0.1:8080`；后端使用其他端口时修改 `web/vite.config.ts` 的代理配置。

`npm --prefix web run build` 执行 TypeScript 检查并生成 `web/dist`，随后重新编译 Go 才会嵌入新资源。构建产物提交到仓库，因此 `go install` 和仅使用 Go 的构建方式仍然可用；修改前端源码时必须同步生成构建产物。`make build` 会安装锁定的前端依赖、构建资源并编译 CLI。使用 `npm --prefix web test` 和 `go test ./...` 运行回归测试。Spec watch 会更新浏览器中的文档并保留请求草稿；加载失败时保留上一份正常文档。

## 检查文档

无需启动网页服务即可检查 spec：

```bash
apicat-cli check                         # 当前项目的 spec 目录
apicat-cli check /path/to/spec-dir       # 显式指定目录
apicat-cli check -json                   # 供 AI 和 CI 使用的结构化输出
apicat-cli check -json /path/to/spec-dir
```

参数放在目录之前。若目录恰好名为 `check`，仍可用 `apicat-cli ./check` 启动预览。

命令加载文档、解析引用，并执行与页面相同的六项文档质量检查。一个接口即使有多个 tag，也只检查一次。报告包含文档级和接口级提示；字段描述缺失按接口汇总数量。

退出码：**0** = 加载成功且无质量提示；**1** = 存在质量提示；**2** = 加载、参数或输出错误。`-json` 在 stdout 输出一个 JSON 对象，包含 `status`（`ok`、`hints` 或 `error`）、`spec_dir`、`errors`、`lint_total` 和 `hints`（`method`、`path`、`message`，文档级提示省略 method/path）。参数错误和解析器诊断写入 stderr。加载失败时 lint 未执行，不能将零计数视为通过，应先检查 `status` 和 `errors`。

这是加载和文档质量检查，不是完整的 OpenAPI 规范校验：不检查 operationId 唯一性、与实现的一致性、未使用的 schema 或超过渲染深度的字段。AI 创建或调整文档后应执行此命令，修复范围内的问题并再次检查，再报告完成情况。

## Agent Skill

[`skills/apicat-spec/`](skills/apicat-spec/SKILL.md) 是一个 [Agent Skill](https://code.claude.com/docs/en/skills)，教会 coding agent 按 apicat 的方式编写 OpenAPI spec——存放到 apicat 的项目目录、渲染出来零 lint 提示。在你的项目中使用：

```bash
mkdir -p .claude/skills
cp -r skills/apicat-spec .claude/skills/
```

之后让 agent 创建或更新 API spec，它会自动遵循这些约定：spec 的存放位置（从本地或兜底 `setting.json` 读取 `spec_dir`，未设置时使用 `~/.apicat/projects/<编码路径>/`）、入口文件命名、多文件 `$ref` 组织、必填的 description、错误响应、example 等。Skill 文件是纯 Markdown，其他支持技能/指令文件的工具同样可用。

## 路线图

- **云端** — `apicat-cli publish` 发布到可自部署的服务端：文档共享、团队协作、spec 版本管理（`latest` + 固定 tag）
- **Agent 原生访问** — 原始 spec 端点、API 发现与搜索、MCP Server（原生对接 Claude Code / Cursor 等）

## 关于 v3 重写

v3 是彻底的重写。此前的 2.x 代码保留在 [`old`](https://github.com/apicat/apicat/tree/old) 分支，2.x 的各个发布版本仍可在 [tags](https://github.com/apicat/apicat/tags) 中获取。

## 许可证

[AGPL-3.0](LICENSE)
