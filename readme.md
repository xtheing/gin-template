## 项目说明
这是一个使用`gin`框架搭建的基础服务端系统，包括了验证`jwt`的验证，数据库的连接范例，`model`数据库定义范例，接口格式的定义，认证中间件等基础功能。

可以基于此框架的基础上快速的进行开发任务。

## 🔒 安全优化（最新更新）

### 已完成的安全增强：

1. **JWT 密钥安全化**
   - JWT 密钥从硬编码移至配置文件
   - 支持环境变量配置
   - 添加密钥生成工具

2. **密码策略加强**
   - 实现密码复杂度验证（长度、大小写、数字、特殊字符）
   - 密码强度评分系统（0-100分）
   - 常见弱密码检测
   - 重复字符检查

3. **敏感配置环境变量化**
   - 数据库密码支持环境变量
   - JWT 配置支持环境变量
   - 提供完整的 `.env_example` 示例

4. **数据库连接优化**
   - 添加数据库连接池配置
   - 优化连接管理
   - 改进错误处理

## 🚀 快速开始

### 克隆项目

```bash
$ git clone https://gitee.com/theing/gin_base.git
$ cd gin-template
```

### 环境配置

1. 复制环境变量示例文件：
```bash
$ cp .env_example .env
```

2. 生成安全密钥：
```bash
$ go run cmd/generate_key.go
```

3. 编辑 `.env` 文件，设置你的数据库和JWT配置

### 运行

```bash
$ go run main.go
```

### 编译运行

```bash
$ go build && ./gin-template
```

项目默认使用8080端口下运行

## 📝 环境变量配置

项目支持通过环境变量进行配置，主要配置项：

- `TPL_JWT_SECRET`: JWT 密钥（必须设置）
- `TPL_JWT_EXPIRE_HOURS`: JWT 过期时间（小时）
- `TPL_JWT_ISSUER`: JWT 签发者
- `TPL_POSTGRES_*`: PostgreSQL 数据库配置
- `TPL_MYSQL_*`: MySQL 数据库配置
- `TPL_SERVER_PORT`: 服务端口
- `TPL_IS_DEBUG`: 调试模式
- `TPL_IS_PRINT_SQL`: SQL 打印模式

## 🔧 密码安全策略

系统现在包含强密码验证：

- **最低要求**：至少6位，包含小写字母和数字
- **推荐要求**：8位以上，包含大小写字母、数字、特殊字符
- **评分系统**：0-100分，60分以上才允许注册
- **弱密码检测**：自动检测常见弱密码模式

## 🛠️ 开发工具

### 密钥生成工具
```bash
$ go run cmd/generate_key.go
```
生成安全的 JWT 密钥和数据库密码。

### Docker 支持
```bash
$ docker build -t gin-template .
$ docker run -p 8080:8080 --env-file .env gin-template
```

TODO
---
- [ ] 集成swagger
- [ ] 请求失败返回结构处理
- [X] 登陆验证请求体的格式
- [ ] 用户反馈功能
- [ ] jupyter 空间的启动
- [ ] 各个厂商短信发送的实现，验证码
- [ ] 文件上传的实现
- [x] 升级整个项目，使用最新的golang环境
- [x] 更改配置文件，适配环境变量获取配置的方式。
- [x] 配置dockerfile 和镜像

更新日志
---
2023-10-25
- 使用最新golang环境

2023-10-20
- 增加dockerfile
- 使用环境变量获取配置
