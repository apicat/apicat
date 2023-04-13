# ApiCat

[English](https://github.com/apicat/apicat/blob/master/README.md) | 简体中文

ApiCat 是一款基于 AI 技术的 API 开发工具，它旨在通过自动化和智能化的方式，帮助开发人员更快速、更高效地开发 API。ApiCat 支持 OpenAPI 和 Swagger 的数据文件导入和导出，并可以对用户输入的 API 需求进行分析和识别，自动生成相应的 API 文档和代码等内容。

您可以访问我们的 [在线 Demo](http://demo.apicat.net) 进行试用。

ApiCat 目前还在早期阶段，欢迎 Star 和 Watch 来关注项目的最新动态。

## 功能特性

- 支持 OpenAPI 和 Swagger 的数据文件导入和导出，方便开发人员进行 API 规范描述和管理。
- 通过 AI 技术，可以自动识别 API 的需求和结构，生成相应的 API 文档和代码等内容，提高开发效率和质量。

## 安装部署

### 获取代码

```
git clone https://github.com/apicat/apicat.git
```

### 编译和启动服务

```
# 进入项目
cd apicat

# 编译前端项目
cd fronted # 到前端项目下
pnpm install
pnpm build

# 更新整理依赖
go mod tidy

# 编译项目
go build

# 修改配置文件
vim ./setting.example.yaml

# 启动服务
./apicat
```

## 交流

ApiCat 的成长离不开它的每一位用户，如果你有任何想和我们交流讨论的内容，欢迎和我们联系，通过下方二维码加入我们的微信讨论群。

![Wechat Group](https://cdn.apicat.net/uploads/01bfb23802cdfad49f0d560ee80fc5e3.png)

## 授权许可

[MIT](https://github.com/apicat/apicat/blob/main/LICENSE)