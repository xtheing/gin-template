## 项目说明
这是一个使用`gin`框架搭建的基础服务端系统，包括了验证`jwt`的验证，数据库的连接范例，`model`数据库定义范例，接口格式的定义，认证中间件等基础功能。

可以基于此框架的基础上快速的进行开发任务。

### 克隆项目

```bash
$ git clone https://gitee.com/theing/gin_study.git
```

### 运行

```bash
$ go run main.go routers.go
```

### 编译运行

```bash
$ go build && ./gin_study
```

项目默认使用8080端口下运行