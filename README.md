<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [微信公众号开发](#%E5%BE%AE%E4%BF%A1%E5%85%AC%E4%BC%97%E5%8F%B7%E5%BC%80%E5%8F%91)
- [TODO](#todo)

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

### TODO

- [x] 支持路由组 usage: /music/play/:sid
- [ ] remind
- [ ] ins

