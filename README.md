# EasyGo框架介绍 

## 框架实现依赖类库
- 访问地址路由: github.com/julienschmidt/httprouter
- redis连接池: github.com/gomodule/redigo/redis
- 唯一id生成器: github.com/rs/xid
- mysql数据库orm:  github.com/astaxie/beego/orm
- 数据库驱动(连接池): github.com/go-sql-driver/mysql
- 日志库:            github.com/sirupsen/logrus


## 框架结构设计
./easygo/
```shell
├── api # 接口控制器
│   ├── response.go   （定义全局响应状态码与实现统计输出json结构方法）
│   └── smsActions.go  (定义短信控制器方法)
|── component # 组件层  (实现独立业务组件)
├── conf # 配置文件
│   ├── config.go      (公共配置: 注册环境变量及业务配置信息)
│   ├── mysqlConf.go   (mysql数据库配置信息: 支持多环境配置)
│   └── redisConf.go   (redis数据库配置信息: 支持多环境配置)
├── lib # 公共类库
│   ├── redisPool.go   (实现redis连接池)
│   └── safeVerify.go  (安全类库: ip, 地址签名等校验)
├── log # 日志类库
│   ├── logger.go      
│   └── logrus.go      (logrus日志类库统一封装)
├── main.go
├── models # 数据库model
│   └── verifyCode.go  （验证码model文件，对应verify_code表，该文件可有beego框架bee命令创建，自动生成数据库"增删改"数据库操作方法）
└── route # 路由
​    └── router.go      (定义控制器路由表，配置路由信息，详情查看文件)``
​    
​    
## 数据处理流程

1. main.go入口文件，注册mysql驱动、数据库信息，实例化路由表，监听指定端口。
2. 路由转发请求到控制器方法。
3. 控制器接收请求数据，开始参数校验、安全校验、创建日志请求跟踪id，最终通过包名方式调用models发送验证码方法(例子:Smsmodels.SendVerifyCode(phone, businessType))
4. models接受参数数据，并记录日志，开启数据处理流程
   4.1 调用短信渠道路由方法，获取该请求，短信服务商。
   4.2 根据路由信息，通过包名方式调用该渠道发送短信sdk。
   4.3 记录发送短信验证码信息，通过orm方法入库。
   4.4 发送短信统计数据缓存到  redis 。

 5. 控制器根据  model 数据处理结果，返回响应json数据信息。
