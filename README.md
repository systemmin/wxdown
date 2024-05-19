# wxdown 微信公众号离线文章保存

## 简介

本来一开始用 `nodejs` 写的，考虑大小、易操作、高性能、跨平台以及环境等问题，我就想能不能搞个不需依赖开发语言环境就能运行的。所以我就选择 `go`并且它本身就具备以上优点。作者本身是`java`开发，第一次使用 `go`，所以过程也是比较艰难，好在 GPT 在学习一门新的开发语言方面还是相当给力！💪💪💪

这是一个用于保存公众号文章到本地离线查看的软件，支持将 HTML 文章保存至本地，并提供 HTML 转 PDF 的功能。此外，软件还支持图片素材管理，可在 Windows、Mac 和 Linux 系统上运行，使用 Go 语言开发，具备轻量级、小体积、高性能和并发支持的特点。**不支持批量直接获取文章列表**

## 功能特点

- 保存公众号文章至本地
- 支持将 HTML 文章转换为 PDF 格式（需安装 [wkhtmltopdf](https://wkhtmltopdf.org/downloads.html)）
- 图片素材管理
- 保存原始地址
- 跨平台支持：Windows、Mac 和 Linux
- 使用 Go 语言开发，轻量级、高性能、高并发
- 提供简单易用的 Web 界面管理

## 下载

| 操作系统               | 版本/架构                                                    | 大小   |
| ---------------------- | ------------------------------------------------------------ | ------ |
| Windows                | [wxdown-1.0.0-windows-amd64.exe](https://864000.lanzouj.com/i5JaN1z84hlc) | 9.23MB |
| ARM Linux              | [wxdown-1.0.0-linux-arm64](https://864000.lanzouj.com/iueMH1z84e9c) | 8.75MB |
| Linux                  | [wxdown-1.0.0-linux-amd64](https://864000.lanzouj.com/iueMH1z84e9c) | 9.04MB |
| macOS                  | [wxdown-1.0.0-darwin-amd64](https://864000.lanzouj.com/ikGp81z84i7e) | 9.09MB |
| macOS（Apple Silicon） | [wxdown-1.0.0-darwin-arm64](https://864000.lanzouj.com/izIaY1z84iuh) | 8.75MB |

## 安装和运行

### Windows

包含了 `wkhtmltopdf`

1. 解压缩包
2. 打开目录
3. 双击 `wxdown-1.0.0-windows-amd64.exe` 启动
4. 浏览器访问 http://127.0.0.1:81

如下所示启动成功：

```bash
cwd： E:\code\go\go-wx-download
----------------------------------------
        欢迎使用 wxdown 工具！
----------------------------------------
运行模式 : binary
软件版本 : 1.0.0
操作系统 : windows
系统架构 : amd64
启动时间 : 2024-05-19 00:00:00
----------------------------------------
服务信息
----------------------------------------
服务地址：
        http://192.168.31.209:81        (浏览器访问)
        http://192.168.202.1:81 (浏览器访问)
        http://192.168.11.1:81  (浏览器访问)
        http://172.26.192.1:81  (浏览器访问)
        http://127.0.0.1:81     (浏览器访问)
采集接口：
        http://192.168.31.209:81/gather/        (GET|POST|HEAD)
        http://192.168.202.1:81/gather/ (GET|POST|HEAD)
        http://192.168.11.1:81/gather/  (GET|POST|HEAD)
        http://172.26.192.1:81/gather/  (GET|POST|HEAD)
        http://127.0.0.1:81/gather/     (GET|POST|HEAD)
----------------------------------------
配置信息
----------------------------------------
运行路径 : E:\code\go\go-wx-download
资源路径 : E:\code\go\go-wx-download\data
```

### Linux

添加权限

```bash
chmod +x wxdown-1.0.0-linux-amd64
```

启动程序

```
root@mac-max:/home/wx# ./wxdown-1.0.0-linux-amd64 
cwd： /home/wx
----------------------------------------
        欢迎使用 wxdown 工具！
----------------------------------------
运行模式 : binary
软件版本 : 1.0.0
操作系统 : linux
系统架构 : amd64
启动时间 : 2024-05-19 00:00:00
----------------------------------------
服务信息
----------------------------------------
服务地址：
        http://192.168.31.156:81        (浏览器访问)
        http://172.17.0.1:81    (浏览器访问)
        http://172.18.0.1:81    (浏览器访问)
        http://127.0.0.1:81     (浏览器访问)
采集接口：
        http://192.168.31.156:81/gather/        (GET|POST|HEAD)
        http://172.17.0.1:81/gather/    (GET|POST|HEAD)
        http://172.18.0.1:81/gather/    (GET|POST|HEAD)
        http://127.0.0.1:81/gather/     (GET|POST|HEAD)
----------------------------------------
配置信息
----------------------------------------
运行路径 : /home/wx
资源路径 : /home/wx/data
```

### Mac 

出现 `permission denied`  表示没有权限

```
(base) mac@macdeMacBook-Pro-3 ~ % /Users/mac/Desktop/wxdown-1.0.0-darwin-amd64/wxdown-1.0.0-darwin-amd64  
zsh: permission denied: /Users/mac/Desktop/wxdown-1.0.0-darwin-amd64/wxdown-1.0.0-darwin-amd64
```

添加权限

```bash
(base) mac@macdeMacBook-Pro-3 ~ % chmod +x /Users/mac/Desktop/wxdown-1.0.0-darwin-amd64/wxdown-1.0.0-darwin-amd64
```

双击 `wxdown-1.0.0-darwin-amd64` 启动或命令启动

```bash
(base) mac@macdeMacBook-Pro-3 ~ %  /Users/mac/Desktop/wxdown-1.0.0-darwin-amd64/wxdown-1.0.0-darwin-amd64
```

执行结果同上

**简单使用会下载和安装就可以了，后面都基本没啥用了**😄😄，不用再看了

## 目录结构

- `web`：HTML 页面，很简单也可以自己修改
  - `index.html` 主页面
  - `images.html` 图片预览页面
- `config.yaml`：系统全局配置文件
- `wxdown-1.0.0` 可执行文件，程序入口

 `config.yaml`

```yaml
# 服务端口
port: 81

# 本地数据文件存储路径
path: ./data

# HTML 转 PDF 配置
# 下载 wkhtmltopdf 路径 https://wkhtmltopdf.org/downloads.html
# window 建议下载后将 wkhtmltopdf目录下载所有内容拷贝到项目根目录下
wkhtmltopdf:
  # true 开启 false 关闭 默认关闭
  enable: false
  # linux 例如：/usr/local/wkhtmltopdf/bin/
  # window 例如：E:\Program Files\wkhtmltopdf\bin
  path:

# 采集线程配置
thread:
  # 同时下载 HTML 线程数量
  html: 5
  # 同时下载图片线程数量
  image: 10
```


## 接口

### 采集接口

- `GET`仅支持单次下载，`POST` 支持批量提交，请求头类型 `JSON` 格式
- http://127.0.0.1:81/gather/+需采集地址。就可以直接把地址发给采集软件

| 地址                        | 请求方式 | 请求参数                           | 请求体                                                       |
| --------------------------- | -------- | ---------------------------------- | ------------------------------------------------------------ |
| http://127.0.0.1:81/gather/ | GET      | /gather/https://mp.weixin.qq.com/1 | 无                                                           |
|                             | HEAD     | /gather/https://mp.weixin.qq.com/1 | 无                                                           |
|                             | POST     | /gather/                           | ["https://mp.weixin.qq.com/1","https://mp.weixin.qq.com/2",...] |

#### 书签脚本

注意⚠️：如果启动软件的机器和浏览文章的机器不是一台机器，使用局域网 `IP`（192.168.0.xxx）替换 `127.0.0.1`

```js
javascript:fetch("http://127.0.0.1:81/gather/" + window.location.href,{mode:"no-cors"});
```

使用方法：

1. 浏览器书签栏➡️右键➡️添加网页...➡️名称：随便你能记住就行➡️网址：输入下面`js`脚本
2. 打开浏览器公众号文章
3. 点击上面添加的**书签脚本**软件会自动采集

### 资源接口

| 地址                         | 请求方式 | 请求参数 | 请求体 |
| ---------------------------- | -------- | -------- | ------ |
| http://127.0.0.1:81/articles | GET      | 无       | 无     |

### 打开文件夹接口

| 地址                      | 请求方式 | 请求参数         | 请求体 |
| ------------------------- | -------- | ---------------- | ------ |
| http://127.0.0.1:81/open/ | GET      | /open/公众号名称 | 无     |

## 使用示例

### 主页面

![主界面](https://raw.githubusercontent.com/systemmin/wxdown/master/doc/2.jpg)

### 文章列表

![文章列表](https://raw.githubusercontent.com/systemmin/wxdown/master/doc/1.jpg)

### 图片库

![图库](https://raw.githubusercontent.com/systemmin/wxdown/master/doc/3.jpg)

## 开发和贡献

本软件使用 Go 语言开发，欢迎开发者贡献代码或提出改进建议。请在 [GitHub](https://github.com/systemmin/wxdown/issues) 上提交 issue 或 pull request。

## 注意事项

- 请勿传播未经授权的文章或图片。
- 在保存、转换文章或管理图片时，请注意版权和法律规定。
