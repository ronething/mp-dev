<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [微信公众号开发](#%E5%BE%AE%E4%BF%A1%E5%85%AC%E4%BC%97%E5%8F%B7%E5%BC%80%E5%8F%91)
- [实现思路](#%E5%AE%9E%E7%8E%B0%E6%80%9D%E8%B7%AF)
- [使用](#%E4%BD%BF%E7%94%A8)
- [TODO](#todo)
- [致谢](#%E8%87%B4%E8%B0%A2)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

### 微信公众号开发

- music 目前仅支持网抑云

```
usage:
- /help # 查看帮助
- /music/play/:sid # 播放音乐 sid 为歌曲 id
- /music/url/:sid  # 获取音乐下载链接
- /music/search/:keywords   # 通过关键字搜索歌曲
- /music/search/:keywords/:page # 分页搜索歌曲
- /music/:name # 搜索并播放歌曲 默认取第一首
```

### 实现思路

[Go 语言之微信公众号开发](https://www.jianshu.com/p/ab9d10a172a0)

### 使用

- 依赖环境

    * go 1.13+ (其他版本没有测试过)

- 测试

> 假设你已经正确填写好配置文件

```sh
git clone https://github.com/ronething/mp-dev.git
cd mp-dev && make build && cd bin
./wechat-mp -c ./example.yaml
```

测试期间可使用内网穿透,如 [ngrok](https://ngrok.com)

- 部署

> 假设你已经正确填写好配置文件
> 并且服务器是 linux amd64 架构

```sh
git clone https://github.com/ronething/mp-dev.git
cd mp-dev && make deploy
```

- 体验

![](./asserts/wechatsearch.png)

### TODO

- [x] 支持路由组 usage: /music/play/:sid
- [ ] remind/ins

### 致谢

- [NeteaseCloudMusicApi](https://github.com/Binaryify/NeteaseCloudMusicApi)
- [silenceper/wechat](https://github.com/silenceper/wechat)
- [echo](https://github.com/labstack/echo)
- [upx](https://github.com/upx/upx)
