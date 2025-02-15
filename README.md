<p align="center">
<img src="web/logo.png">
<p>

# wxdown

这是一个用于保存公众号文章到本地离线查看的软件，支持将 HTML 文章保存至本地，并提供 HTML 转 PDF 的功能。此外，软件还支持图片、音频、视频下载，可在 Windows、Mac 和 Linux 系统上运行，使用 Go
语言开发，具备轻量级、小体积、高性能和并发支持的特点。

**不支持批量直接获取文章列表**

## 功能特点

- 保存公众号文章至本地
- 支持将 HTML 文章转换为 PDF 格式（需安装 [wkhtmltopdf](https://wkhtmltopdf.org/downloads.html))
- 图片、音视频管理
- 支持首页合集、标签合集保存
- 保存原始地址
- 跨平台支持：Windows、Mac 和 Linux
- 使用 Go 语言开发，轻量级、高性能、高并发
- 提供简单易用的 Web 界面管理

## 使用文档

**移步** [wiki](https://github.com/systemmin/wxdown/wiki)

## 软件版本

### v1.0.9 (当前版本)

替换 `wxdown`
- 解决视频解析导致系统中断 bug [issues#4](https://github.com/systemmin/wxdown/issues/15)

[下载地址](https://864000.lanzouu.com/b0hcr5bad) 密码：wxdown

### [历史版本](https://github.com/systemmin/wxdown/blob/master/CHANGELOG.md)

## docker 运行

创建挂载目录：

```shell
mkdir -p /home/wxdown/data
```

启动容器

```shell
docker run -p 81:81 --name wxdown -d registry.cn-hangzhou.aliyuncs.com/wxdown/wxd:latest
```

拷贝配置文件

```shell
docker cp wxdown:/wx/config.yaml /home/wxdown/
```

拷贝证书（可选）

```shell
docker cp wxdown:/wx/certs/ /home/wxdown/
```

停掉&删除容器

```shell
docker stop wxdown && docker rm -f wxdown
```

挂载数据目录

```shell
docker run -p 81:81 --name wxdown \
	-v /home/wxdown/data/:/wx/data/ \
  	-v /home/wxdown/config.yaml:/wx/config.yaml \
  	-d registry.cn-hangzhou.aliyuncs.com/wxdown/wxd:latest
```

## 构建镜像

[docker build](https://github.com/systemmin/wxdown/blob/master/docker/CMD.md)

## 目录结构

```text
project/
│
├── cmd/
│   ├── yourapp/
│   │   └── main.go
│
├── internal/
│   ├── pkg1/
│   │   ├── file1.go
│   │   └── file2.go
│   │
│   └── pkg2/
│       ├── file1.go
│       └── file2.go
│
├── pkg/
│   └── sharedpkg/
│       ├── file1.go
│       └── file2.go
│
├── config/
│   └── config.go
│
├── web/
│   ├── static/
│   │   └── ...
│   │
│   ├── templates/
│   │   └── ...
│   │
│   └── main.go
│
├── tests/
│   └── ...
│
├── build.bat
│
├── config.yaml
│
├── Dockerfile
│
├── README.md
│
├── main.go
│
└── go.mod
```

- `cmd/`: 该目录用于存放项目的可执行文件，每个子目录代表一个可执行文件，命名一般为项目名或者服务名。
- `internal/`: 该目录用于存放项目内部的私有代码，这些代码只能被当前项目使用，不能被外部包导入。通常会按照功能或模块划分子目录，并在其中编写相应的代码文件。
- `pkg/`: 该目录用于存放可被其他项目导入和使用的代码包。这些包应该是稳定的、通用的、可复用的功能模块。比如一些工具函数、常用的数据结构等。
- `config/`: 该目录用于存放项目的配置文件，比如 JSON、YAML 或者 properties 文件等。
- `web/`: 该目录用于存放 Web 应用相关的代码，如前端静态文件、模板文件以及 Web 服务器的代码。
- `tests/`: 该目录用于存放项目的测试代码，通常按照包的结构组织测试文件。
- `build.bat`: window 下批处理打包脚本
- `config.yaml`: 系统配置文件
- `docker`: docker 镜像构建相关内容。
- `README.md`: 项目的说明文档，包括项目的简介、安装、使用方法等信息。
- `go.mod`: Go modules 的配置文件，用于管理项目的依赖关系。

## 安装教程

1. 确保您的计算机上已安装Go语言环境。
   ```shell
   go version 1.22.3
   ```
2. 克隆项目仓库：
   ```shell
   git clone https://github.com/systemmin/wxdown.git
   ```
3. 进入项目目录：
   ```shell
   cd wxdown
   ```
4. 安装项目依赖：
   ```shell
   go mod download
   ```
5. 编译项目：
   ```shell
   go build -ldflags "-X main.runMode=binary -X main.version=1.0.0" -o wxdown.exe  main.go
   ```
6. 运行项目：

   ```shell
   ./wxdown.exe
   ```

## 本地启动

1. 安装项目依赖：

```shell
go mod download
```

2. 根目录执行

```shell
go run main.go
```

## 许可证

本项目遵循 [Apache-2.0](https://spdx.org/licenses/Apache-2.0.html) 许可证。

## 开发和贡献

本软件使用 Go 语言开发，欢迎开发者贡献代码或提出改进建议。请在 [GitHub](https://github.com/systemmin/wxdown/issues) 上提交 issue 或 pull request。

## 注意事项

**注意：** 请确保遵循网站的使用条款和政策。自行承担风险。

- 本软件完全免费。
- 请勿传播未经授权的文章或图片。
- 使用软件造成的影响由使用者承担。
- 仅供学习交流，严禁用于商业用途，请勿传播下载的数据。
- 本脚本所获取的资源完全合法，与浏览器能直接获得的资源一致。
- 在保存、转换文章或管理图片时，请注意版权和法律规定。

## 鼓励作者

**您**的鼓励💪就是我前进的动力🫶鼓励方式有多种，选择你喜欢的😍，当然不鼓励也没关系哦！😄

- [GitHub Start](https://github.com/systemmin/wxdown) ⭐
- 打赏💰（*量力而行*）

<img src="https://dtking.cn/pay.png" alt="赞赏" style="zoom: 67%;" />