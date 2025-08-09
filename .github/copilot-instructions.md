# Copilot Instructions for ScriptList

## 项目概览

ScriptList是一个用户脚本分享平台的后端服务，基于Go语言和[CaGo框架](https://github.com/cago-frame/cago)
构建。项目采用清洁架构模式，集成了多种中间件和微服务组件。

## 核心架构模式

### 分层架构 (Repository-Service-Controller)

- **Controller层**: `internal/controller/*_ctr/` - 处理HTTP请求，路由绑定
- **Service层**: `internal/service/*_svc/` - 业务逻辑处理，事务管理
- **Repository层**: `internal/repository/*_repo/` - 数据访问层，数据库操作
- **Entity层**: `internal/model/entity/` - 数据模型定义

### 依赖注册模式

在`cmd/app/main.go`中，所有repository必须通过`RegisterXxx(NewXxxRepo())`方式注册才能使用。这是CaGo框架的约定：

```go
script_repo.RegisterScript(script_repo.NewScriptRepo())
user_repo.RegisterUser(user_repo.NewUserRepo())
```

## 关键开发工作流

### 本地开发环境设置

```bash
# 1. 复制配置文件
cp configs/config.yaml.example configs/config.yaml

# 2. 启动依赖服务 (推荐)
docker-compose up -d

# 3. 运行应用
go run ./cmd/app/main.go

# 4. 调试登录 (dev环境)
# 访问 http://127.0.0.1:8080/api/v2/login/debug 登录uid=1的用户
```

### 构建和部署

- 使用`Makefile`进行构建: `make build`
- CI/CD通过`.gitea/workflows/deploy.yaml`实现自动部署
- 支持多环境: `prod`(生产), `pre`(预发布), `test`(测试)

### 代码生成和工具

```bash
make generate    # 生成mock和CaGo代码
make lint-fix    # 自动修复lint问题
make test        # 运行测试
```

## 项目特有模式

### 权限控制系统

基于角色的访问控制，支持用户和组级别权限：

- **角色**: `owner` > `manager` > `guest`
- **权限检查**: 使用`script_svc.Access().CheckHandler("resource", "action")`中间件
- **上下文传递**: 通过`CtxScript(ctx)`、`CtxAccess(ctx)`获取当前脚本/权限信息

### 异步消息处理

使用Producer-Consumer模式处理异步任务：

- **Producer**: `internal/task/producer/` - 发布消息
- **Consumer**: `internal/task/consumer/subscribe/` - 订阅处理
- **主要消息**: 脚本创建/更新、统计数据、ES同步

### 多数据存储集成

- **MySQL**: 主要业务数据
- **Elasticsearch**: 脚本搜索和全文检索
- **ClickHouse**: 统计数据存储
- **Redis**: 缓存和会话管理

### API定义规范 (CaGo框架)

使用声明式结构体定义API，在`internal/api/`下按模块组织：

```go
// 请求结构体 - 包含路由元信息
type CreateIssueRequest struct {
	mux.Meta `path:"/scripts/:id/issues" method:"POST"`
	ScriptID int64    `uri:"id" binding:"required"`
	Title    string   `json:"title" binding:"required,max=128" label:"标题"`
	Content  string   `json:"content" binding:"max=10485760" label:"反馈内容"`
	Labels   []string `json:"labels" binding:"max=128" label:"标签"`
}

// 响应结构体
type CreateIssueResponse struct {
	ID int64 `json:"id"`
}
```

**绑定说明**: `path`(请求路径)、`method`(请求方法)、`uri`(路径参数)、`json`(JSON体)、`form`(表单/查询参数)、`binding`(
验证规则)

定义好API后，通过`make generate`生成对应的处理代码。

### 脚本元数据解析

脚本代码中的`@xxx`元数据会被自动解析为JSON存储：

```go
// 在script_entity.Code.UpdateCode()中处理
// 提取@match、@include等元数据用于分类和域名匹配
```

## 关键集成点

### OAuth认证集成

与油猴中文网(bbs.tampermonkey.net.cn)强关联，共享用户体系。

### 统计数据收集

支持脚本下载量、更新量、访问量等多维度统计，数据通过消息队列异步处理。

## 数据库迁移

- 迁移文件位于`migrations/`目录
- 使用gormigrate管理版本
- 支持分布式锁防止并发迁移
- 包含从Redis到关系型数据库的数据迁移逻辑

## 调试技巧

- **调试登录**: dev环境可访问`/api/v2/login/debug`
- **API文档**: Swagger文档通过`docs/`目录生成
- **日志**: 使用结构化日志，支持链路追踪
- **统计ES迁移**: POST `/api/v2/scripts/migrate/es`进行全量数据迁移

## 性能优化要点

- Redis用作缓存层，减少数据库查询
- 异步处理统计数据，避免阻塞主流程
- ES搜索数据批量同步，降低写入频率
- ClickHouse处理大量统计数据的OLAP查询

添加新功能时，确保遵循现有的分层架构模式，注册必要的repository依赖，考虑是否需要异步处理，以及相应的权限控制。
