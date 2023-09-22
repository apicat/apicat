<div align="center">
    <img alt="ApiCat" width="350px" src="https://cdn.apicat.net/uploads/2d02ff2f6b19d3d6d3f134c1872484aa.png"/>
</div>

# ApiCat

[English](https://github.com/apicat/apicat/blob/master/README.md) | 简体中文

ApiCat 是一款基于 AI 技术的 API 开发工具，它旨在通过自动化和智能化的方式，帮助开发人员更快速、更高效地开发 API。通过使用 ApiCat，开发人员可以在 API 的开发和管理上节省大量的时间。

访问我们的 [在线 Demo](https://apicat.zeabur.app) 进行试用。

## 功能特性

- **API 文档:** 你可以非常方便的创建和管理你的 API 文档
- **AI 支持:** 你可以通过 AI 帮助你快速生成 API 的文档、模型、响应等内容
- **Mock:** 简单好用的 Mock 功能让 API 开发更加迅速
- **迭代:** 清晰的迭代规划能让团队可以明确 API 的变动范围，保证每次开发任务的效率和质量
- **数据导入导出:** 可以将 API 数据全量导入到任何支持 OpenAPI 或 Swagger 的软件中，也可以反向导入给 ApiCat

## 安装部署

### 五种安装部署方式

#### 1. 下载可执行文件部署

##### 第一步：下载打包好的可执行文件

下载已经打包好的可执行文件 [下载地址](https://github.com/apicat/apicat/releases)

##### 第二步：启动

```
# 默认配置或加载环境变量配置启动服务
./apicat
# 通过配置文件启动服务
./apicat -c setting.example.yaml
```

#### 2. DockerHub 安装部署

##### 第一步：下载镜像

```
docker pull natuo/apicat:latest
```

##### 第二步：启动

```
docker run --name apicat-server -p 8000:8000 -d --link mysql natuo/apicat:latest -c /app/setting.default.yaml
```

#### 3. Zeabur 一键安装部署

注册 [Zeabur](https://zeabur.com/)，在 Marketplace 找到 ApiCat 一键部署。

#### 4. 本地 Docker 安装部署

##### 第一步：获取代码

```
git clone https://github.com/apicat/apicat.git
```

##### 第二步：构建本地镜像

```
docker build -t apicat:latest .
```

##### 第三步：启动

```
docker run --name apicat-server -p 8000:8000 -d --link mysql apicat:latest -c /app/setting.default.yaml
```

#### 5. 源代码安装部署

##### 第一步：获取代码

```
git clone https://github.com/apicat/apicat.git
```

##### 第二步：前端安装和编译

```
cd frontend
pnpm i
pnpm build
```

##### 第三步：后端安装和编译

```
go mod tidy
go build
```

##### 第四步：启动

```
# 默认配置或加载环境变量配置启动服务
./apicat
# 通过配置文件启动服务
./apicat -c setting.default.yaml
```

### 配置项说明

你可以通过两种方式设置自定义配置来启动 ApiCat

#### 1. 读取配置文件

参见 [backend/config/setting.example.yaml](https://github.com/apicat/apicat/blob/main/backend/config/setting.example.yaml)

#### 2. 读取环境变量

| 变量名称 | 描述 | 示例 |
| ------- | --- | ---- |
| APICAT_APP_NAME | 应用名称 | ApiCat |
| APICAT_APP_HOST | 绑定的 IP 地址，默认 0.0.0.0 | 0.0.0.0 |
| APICAT_APP_PORT | 绑定的端口，默认 8000 | 8000 |
| APICAT_LOG_PATH | 日志文件地址，为空输出到 stdout | logs/ |
| APICAT_LOG_LEVEL | 日志等级 | debug |
| APICAT_DB_HOST | MySQL IP 地址，必须 | 127.0.0.1 |
| APICAT_DB_PORT | MySQL 端口，必须 | 3306 |
| APICAT_DB_USER | MySQL 用户名，必须 | root |
| APICAT_DB_PASSWORD | MySQL 密码，必须 | 123456 |
| APICAT_DB_NAME | MySQL 数据库名称，必须 | apicat |
| APICAT_OPENAI_SOURCE | OpenAI 调用途径(openai, azure) | openai |
| APICAT_OPENAI_KEY | OpenAI Key | sk-xxxxxx |
| APICAT_OPENAI_ENDPOINT | OpenAI 调用终端地址，当 APICAT_OPENAI_SOURCE 为 azure 时有效 | https://xxxxxx.openai.azure.com/ |

## 交流

如果你有任何想和我们交流讨论的内容，欢迎加入我们的微信讨论群。

![Wechat Group](https://cdn.apicat.net/uploads/01bfb23802cdfad49f0d560ee80fc5e3.png)

## 功能截图

![AI-generate-schema](https://cdn.apicat.net/uploads/0c3518c1bfc421fc4f3f86c085f353d2.gif)

![AI-generate-api-by-schema](https://cdn.apicat.net/uploads/bbcae83511d797d22077d05d17c262cc.gif)

![AI-generate-api](https://cdn.apicat.net/uploads/cf617b56fa186960c228c79487cf6c5e.gif)

## 授权许可

[MIT](https://github.com/apicat/apicat/blob/main/LICENSE)