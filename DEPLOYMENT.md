# 部署指南

## 🚀 概述

本文档详细说明了如何部署 gin-template 项目到不同的环境中，包括开发、测试和生产环境。

## 📋 目录

- [环境要求](#环境要求)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [部署方式](#部署方式)
- [监控配置](#监控配置)
- [安全配置](#安全配置)
- [故障排除](#故障排除)

## 🛠️ 环境要求

### 基础要求
- **Go**: 1.21.0 或更高版本
- **数据库**: PostgreSQL 12+ 或 MySQL 8.0+
- **缓存**: Redis 6.0+
- **操作系统**: Linux (推荐 Ubuntu 20.04+), macOS, Windows

### 可选组件
- **Docker**: 20.10+ (用于容器化部署)
- **Prometheus**: 2.30+ (用于监控)
- **Grafana**: 8.0+ (用于可视化)
- **Nginx**: 1.18+ (用于反向代理)

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/xtheing/gin-template.git
cd gin-template
```

### 2. 安装依赖
```bash
# 安装 Go 依赖
go mod download

# 安装开发工具
make install-tools
```

### 3. 配置环境
```bash
# 复制环境变量模板
cp .env_example .env

# 编辑配置文件
vim .env
```

### 4. 启动服务
```bash
# 开发模式
make dev

# 或者直接运行
go run main.go
```

### 5. 验证部署
```bash
# 健康检查
curl http://localhost:8080/api/health/

# API 测试
curl http://localhost:8080/api/auth/info
```

## ⚙️ 配置说明

### 环境变量配置

#### 必需配置
```bash
# 服务配置
TPL_HOST=0.0.0.0
TPL_PORT=8080
TPL_MODE=release

# 数据库配置
TPL_DB_HOST=localhost
TPL_DB_PORT=5432
TPL_DB_USER=postgres
TPL_DB_PASSWORD=password
TPL_DB_NAME=gin_template

# JWT 配置
TPL_JWT_SECRET=your-super-secret-jwt-key-here
```

#### 可选配置
```bash
# 缓存配置
TPL_CACHE_HOST=localhost
TPL_CACHE_PORT=6379
TPL_CACHE_PASSWORD=
TPL_CACHE_DB=0

# 日志配置
TPL_LOG_LEVEL=info
TPL_LOG_FILE=/var/log/gin-template/app.log

# 安全配置
TPL_TRUSTED_PROXIES=127.0.0.1,::1
TLS_CONFIG_PATH=/etc/ssl/certs/tls.json
```

### 数据库配置

#### PostgreSQL 配置
```sql
-- 创建数据库
CREATE DATABASE gin_template;

-- 创建用户
CREATE USER gin_template_user WITH PASSWORD 'password';

-- 授权
GRANT ALL PRIVILEGES ON DATABASE gin_template TO gin_template_user;

-- 连接数据库
\c gin_template;

-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

#### MySQL 配置
```sql
-- 创建数据库
CREATE DATABASE gin_template CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户
CREATE USER 'gin_template_user'@'%' IDENTIFIED BY 'password';

-- 授权
GRANT ALL PRIVILEGES ON gin_template.* TO 'gin_template_user'@'%';

-- 刷新权限
FLUSH PRIVILEGES;
```

## 🐳 部署方式

### 1. 直接部署

#### 编译应用
```bash
# 生产环境编译
make build-prod

# 或者手动编译
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gin-template main.go
```

#### 系统服务配置 (systemd)
```ini
# /etc/systemd/system/gin-template.service
[Unit]
Description=Gin Template API Service
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=gin-template
Group=gin-template
WorkingDirectory=/opt/gin-template
ExecStart=/opt/gin-template/bin/gin-template
Restart=always
RestartSec=5
Environment=GIN_MODE=release
EnvironmentFile=/opt/gin-template/.env

# 安全配置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/gin-template

[Install]
WantedBy=multi-user.target
```

#### 启动服务
```bash
# 创建用户
sudo useradd -r -s /bin/false gin-template

# 部署文件
sudo cp bin/gin-template-linux /opt/gin-template/bin/
sudo cp .env /opt/gin-template/
sudo chown -R gin-template:gin-template /opt/gin-template/

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable gin-template
sudo systemctl start gin-template

# 查看状态
sudo systemctl status gin-template
sudo journalctl -u gin-template -f
```

### 2. Docker 部署

#### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gin-template main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/gin-template .
COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./gin-template"]
```

#### docker-compose.yml
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - TPL_DB_HOST=postgres
      - TPL_CACHE_HOST=redis
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: gin_template
      POSTGRES_USER: gin_template_user
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    depends_on:
      - app
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

#### 部署命令
```bash
# 构建镜像
docker build -t gin-template:latest .

# 使用 docker-compose
docker-compose up -d

# 查看日志
docker-compose logs -f app
```

### 3. Kubernetes 部署

#### deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gin-template
  labels:
    app: gin-template
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gin-template
  template:
    metadata:
      labels:
        app: gin-template
    spec:
      containers:
      - name: gin-template
        image: gin-template:latest
        ports:
        - containerPort: 8080
        env:
        - name: TPL_DB_HOST
          value: "postgres-service"
        - name: TPL_CACHE_HOST
          value: "redis-service"
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /api/health/
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/health/
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: gin-template-service
spec:
  selector:
    app: gin-template
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
```

#### 部署命令
```bash
# 应用配置
kubectl apply -f deployment.yaml

# 查看状态
kubectl get pods -l app=gin-template
kubectl logs -f deployment/gin-template
```

## 📊 监控配置

### Prometheus 配置

#### prometheus.yml
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "rules/*.yml"

scrape_configs:
  - job_name: 'gin-template'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/api/metrics'
    scrape_interval: 15s

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
```

#### 告警规则
```yaml
# rules/gin-template.yml
groups:
  - name: gin-template
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status_code=~"5.."}[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} errors per second"

      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 0.5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "95th percentile response time is {{ $value }} seconds"

      - alert: DatabaseConnectionHigh
        expr: database_connections{state="open"} > 80
        for: 3m
        labels:
          severity: critical
        annotations:
          summary: "Database connection count is high"
          description: "Database has {{ $value }} open connections"

      - alert: ServiceDown
        expr: up{job="gin-template"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service is down"
          description: "Gin Template service has been down for more than 1 minute"
```

### Grafana 仪表板

#### 数据源配置
```json
{
  "name": "Gin Template",
  "type": "prometheus",
  "url": "http://prometheus:9090",
  "access": "proxy",
  "isDefault": true
}
```

#### 仪表板面板
- HTTP 请求量和响应时间
- 数据库连接和查询性能
- 缓存命中率和操作统计
- 系统资源使用情况
- 错误率和成功率

## 🔒 安全配置

### HTTPS 配置

#### TLS 配置文件
```json
{
  "cert_file": "/etc/ssl/certs/server.crt",
  "key_file": "/etc/ssl/certs/server.key",
  "ca_file": "/etc/ssl/certs/ca.crt",
  "server_name": "api.example.com",
  "client_auth": "require-and-verify"
}
```

#### Nginx 反向代理配置
```nginx
server {
    listen 80;
    server_name api.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate /etc/ssl/certs/api.example.com.crt;
    ssl_certificate_key /etc/ssl/certs/api.example.com.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 安全头
        proxy_set_header X-Frame-Options DENY;
        proxy_set_header X-Content-Type-Options nosniff;
        proxy_set_header X-XSS-Protection "1; mode=block";
    }
}
```

### 防火墙配置

#### UFW 配置
```bash
# 允许 SSH
sudo ufw allow 22/tcp

# 允许 HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 启用防火墙
sudo ufw enable

# 查看状态
sudo ufw status verbose
```

## 🔧 故障排除

### 常见问题

#### 1. 数据库连接失败
```bash
# 检查数据库服务
sudo systemctl status postgresql

# 测试连接
psql -h localhost -U gin_template_user -d gin_template

# 查看日志
sudo tail -f /var/log/postgresql/postgresql-15-main.log
```

#### 2. Redis 连接失败
```bash
# 检查 Redis 服务
sudo systemctl status redis

# 测试连接
redis-cli ping

# 查看日志
sudo tail -f /var/log/redis/redis-server.log
```

#### 3. 应用启动失败
```bash
# 检查配置
./gin-template -config config/application.yml -check

# 查看应用日志
sudo journalctl -u gin-template -f

# 检查端口占用
sudo netstat -tlnp | grep 8080
```

#### 4. 性能问题
```bash
# 查看系统资源
top
htop
iostat -x 1

# 查看应用指标
curl http://localhost:8080/api/metrics

# 数据库慢查询
psql -c "SELECT query, mean_time, calls FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"
```

### 日志分析

#### 应用日志位置
- **Systemd**: `journalctl -u gin-template`
- **文件**: `/var/log/gin-template/app.log`
- **Docker**: `docker logs gin-template`

#### 日志级别
- **error**: 错误信息，需要立即处理
- **warn**: 警告信息，可能的问题
- **info**: 一般信息，正常操作
- **debug**: 调试信息，详细执行过程

### 性能调优

#### 数据库优化
```sql
-- 创建索引
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
CREATE INDEX CONCURRENTLY idx_users_username ON users(username);

-- 分析表统计信息
ANALYZE users;

-- 查看慢查询
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;
```

#### 应用优化
```bash
# 调整 Go 运行时参数
export GOMAXPROCS=4
export GOGC=100

# 启用 pprof
./gin-template -pprof=:6060

# 负载测试
hey -n 10000 -c 100 http://localhost:8080/api/health/
```

## 📚 参考资源

- [Gin 框架文档](https://gin-gonic.com/docs/)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)
- [Redis 文档](https://redis.io/documentation)
- [Prometheus 文档](https://prometheus.io/docs/)
- [Docker 文档](https://docs.docker.com/)
- [Kubernetes 文档](https://kubernetes.io/docs/)

## 🤝 支持

如果在部署过程中遇到问题，请：

1. 查看本文档的故障排除部分
2. 检查项目的 GitHub Issues
3. 提交新的 Issue 并提供详细信息
4. 联系技术支持：support@theing.com

---

**注意**: 在生产环境部署前，请务必完成安全配置和性能调优，确保系统的稳定性和安全性。
