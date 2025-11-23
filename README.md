# 敢学 - Golang交互式学习平台

## 项目概述

敢学是一个完整的Golang交互式学习平台，采用前后端分离架构，包含多个协同工作的组件。该平台允许用户通过浏览器直接编写、运行Go代码，并提供实时反馈，是学习Golang的理想工具。

## 项目架构

敢学项目由三个核心组件构成，通过Docker容器化部署：

- **[敢学 Web](./ganxue-web)**：前端界面，提供用户交互体验
- **[敢学 Server](./ganxue-server)**：后端服务，处理业务逻辑和数据管理
- **[敢学 RunCode](./ganxue-runcode)**：代码执行服务，安全地运行用户提交的代码

### 基础设施组件

- **MySQL**：存储用户数据、学习进度等结构化信息
- **Redis**：缓存、会话管理和消息队列
- **MongoDB**：存储非结构化数据，如学习文档

## 组件说明

### 敢学 Web

前端界面采用HTML、CSS和JavaScript构建，提供以下功能：

- 沉浸式代码编辑器
- 章节式渐进学习路径
- 实时代码执行反馈
- 多设备进度同步
- JWT双令牌认证机制

**详细信息请查看**：[敢学 Web README](./ganxue-web/README.md)

### 敢学 Server

后端服务采用Go语言开发，使用Hertz框架，提供以下功能：

- 用户认证与管理
- 学习内容管理
- 代码执行协调
- 数据存储与检索
- 安全的API接口

**详细信息请查看**：[敢学 Server README](./ganxue-server/README.md)

### 敢学 RunCode

代码执行服务，负责安全地运行用户提交的代码，具有以下特点：

- 隔离的执行环境
- 支持多语言代码执行（Go、C、C++等）
- 系统调用过滤
- 资源限制控制
- 与Redis集成进行任务队列管理

## 快速开始

### 环境要求

- Docker 和 Docker Compose
- Git

### 部署步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/gangantongxue/ganxue.git
   cd ganxue
   ```

2. **配置环境变量**
   - 复制 `.env-example` 文件为 `.env`
   - 编辑 `.env` 文件，设置数据库密码等配置项

3. **启动服务**
   ```bash
   docker-compose up -d
   ```

4. **访问平台**
   - Web界面：http://localhost
   - API服务：http://localhost:8080

## 开发说明

### 目录结构

```
ganxue/
├── ganxue-web/        # 前端项目
├── ganxue-server/     # 后端服务
├── ganxue-runcode/    # 代码执行服务
├── docker-compose.yaml # Docker部署配置
├── .env-example       # 环境变量示例
└── README.md          # 项目说明（当前文件）
```

### 技术栈

- **前端**：HTML + CSS + JavaScript
- **后端**：Go 语言 + Hertz 框架
- **数据库**：MySQL、Redis、MongoDB
- **容器化**：Docker + Docker Compose

## 贡献指南

欢迎对项目进行贡献！请遵循以下步骤：

1. Fork 项目仓库
2. 创建特性分支
3. 提交更改
4. 推送到你的分支
5. 发起 Pull Request

## 许可证

本项目采用 [MIT 许可证](LICENSE)

## 联系方式

如有任何问题或建议，请通过以下方式联系：

- 邮箱：gangantongxue@outlook.com
- 网站：[ganxue.top](https://ganxue.top)