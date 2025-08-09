# 脚本站

[脚本站](https://scriptcat.org)后端

## 启动程序

项目依赖很多中间件，推荐使用`docker-compose`启动，更多的一些配置可以参考[cago框架](https://github.com/cago-frame/cago)

```bash
docker-compose up -d
```

### 配置文件

你可以直接将`configs/config.yaml.example`复制一份为`configs/config.yaml`，然后运行`go run ./cmd/app/main.go`

## 调试

### 登录用户

你可以在`debug模式`&`dev环境`下访问： `http://127.0.0.1:8080/api/v2/login/debug` 登录uid为1的用户

因为项目与[油猴中文网](https://bbs.tampermonkey.net.cn)
是强关联的，所以需要你手动新建一个uid为1的用户在`pre_common_member`表

```sql
INSERT INTO `pre_common_member` (`uid`, `email`, `username`, `password`, `status`, `emailstatus`, `avatarstatus`,
                                 `videophotostatus`, `adminid`, `groupid`, `groupexpiry`, `extgroupids`, `regdate`,
                                 `credits`, `notifysound`, `timeoffset`, `newpm`, `newprompt`, `accessmasks`,
                                 `allowadmincp`, `onlyacceptfriendpm`, `conisbind`, `freeze`)
VALUES (1, 'admin@scriptcat.org', 'admin', '-----', 0, 1, 1, 0, 1, 1, 0, '', 1625882335, 20, 0, '', 0, 0, 0, 1, 0, 0,
        0);
```
