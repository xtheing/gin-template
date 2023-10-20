## 项目说明
这是一个使用`gin`框架搭建的基础服务端系统，包括了验证`jwt`的验证，数据库的连接范例，`model`数据库定义范例，接口格式的定义，认证中间件等基础功能。

可以基于此框架的基础上快速的进行开发任务。

### 克隆项目

```bash
$ git clone https://gitee.com/theing/gin_base.git
```

### 运行

```bash
$ go run main.go
```

### 编译运行

```bash
$ go build && ./gin_study
```

项目默认使用8080端口下运行

TODO
---

- [ ] 用户反馈功能
- [ ] jupyter 空间的启动
- [ ] 各个厂商短信发送的实现，验证码
- [ ] 文件上传的实现
- [ ] 升级整个项目，使用最新的golang环境
- [x] 更改配置文件，适配环境变量获取配置的方式。
- [x] 配置dockerfile 和镜像

更新日志
---
2023-10-20
- 增加dockerfile
- 环境变量配置
