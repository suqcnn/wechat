#wechat

## 安装
```shell
go get -u github.com/xuebing1110/wechat
```

## 基础用法
### 加载配置文件
```go
conf_file := "/etc/wechat.yaml"
wechat_conf, c_err := wechat.LoadConfig(*conf_file)
if c_err != nil {
    log.Panicln(c_err)
}
```

### 初始化数据库
使用XOrm建立数据库表结构;
使用文档：https://lunny.gitbooks.io/xorm-manual-zh-cn/content/
开源项目：https://github.com/go-xorm/xorm

```go
err := wechat.InitDB(wechat_conf.DBTYPE, wechat_conf.DBString,
    &wechat.User{})
if err != nil {
    log.Panic(err)
}
```

> 数据库类型默认为sqlite3

### 开启微信服务
```go
err = wechat_conf.InitWechatAndListen()
if err != nil {
    log.Panic(err)
}
```

## 自定义扩展业务
使用场景包括：
### 指定文本消息的处理(MATCHTYPE_EQUAL)
```go
wechat.Register(wechat.MATCHTYPE_EQUAL, "公交卡余额", TransCardQueryHandler)
```

### 以指定文本为前缀的处理(MATCHTYPE_PREFIX)
```go
wechat.Register(wechat.MATCHTYPE_PREFIX, "姓名@", RegNameHandler)
```

### 上下文的处理(MATCHTYPE_PREFIX)
以用户签到场景为例：

1. 用户输入“绑定”
```go
wechat.Register(wechat.MATCHTYPE_EQUAL, "绑定", UserBindStart)
```
`UserBindStart`实现逻辑，可返回可选产品列表:
```go
func RegNameHandler(hb *wechat.HandleBean) (string, error) {
    return `请选择产品列表:\n1:AAA\n2:BBB\n回复相应数字即可",nil
}
```

2. 用户选择产品信息
```go
wechat.Register(wechat.MATCHTYPE_CONTEX, "绑定", UserChoseProduct, "登陆")
```
`UserChoseProduct`实现逻辑，接收参数并做相应处理:
```go
func UserChoseProduct(hb *wechat.HandleBean, new_status string) (string, error) {
    //选择产品成功
    if hb.Content == "AAA" || hb.Content == "BBB" {
        return wechat.SaveUserStatus(new_status,"请输入用户密码，例如"user/password"")
    } else {//选择产品错误
        
    }
}
```

3. 用户输入用户名及密码，例如"user/password"，服务器做`校验`处理
```go
wechat.Register(wechat.MATCHTYPE_CONTEX, "登陆", UserBindCheck)
```