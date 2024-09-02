<p align="center">
<img src="web/logo.png">
<p>

# wxdown

这是一个用于保存公众号文章到本地离线查看的软件，支持将 HTML 文章保存至本地，并提供 HTML 转 PDF 的功能。此外，软件还支持图片、音频、视频下载，可在
Windows、Mac 和 Linux 系统上运行，使用 Go 语言开发，具备轻量级、小体积、高性能和并发支持的特点。

**不支持批量直接获取文章列表**

## 功能特点

- 保存公众号文章至本地
- 支持将 HTML 文章转换为 PDF 格式（需安装 [wkhtmltopdf](https://wkhtmltopdf.org/downloads.html)）
- 图片、音视频管理
- 支持首页合集、标签合集保存
- 保存原始地址
- 跨平台支持：Windows、Mac 和 Linux
- 使用 Go 语言开发，轻量级、高性能、高并发
- 提供简单易用的 Web 界面管理

### 历史版本

使用文档移步 [wiki](https://github.com/systemmin/wxdown/wiki)


### v1.0.4(当前版本)

替换根目录 `web`、`wxdown` 文件重新启动完成更新

- 解决 [webp](https://mp.weixin.qq.com/s/_eeCF9JLOKF-YeojxsTmog) 图片格式无法转PDF问题
- 解决标签合集不满足分页条件错误问题 [issues#4](https://github.com/systemmin/wxdown/issues/4) 
- 新增图片集文章类型下载 [示例](https://mp.weixin.qq.com/s/2E5aiMre-NO0Vw9rnGTkCQ)

| 操作系统               | 文件名                         | 链接                                    | 文件大小 |
| ---------------------- | ------------------------------ | --------------------------------------- | -------- |
| Windows                | wxdown-1.0.4-windows-amd64.zip | https://864000.lanzouj.com/iWeh328yffqb | 18.2 M   |
| ARM Linux              | wxdown-1.0.4-darwin-arm64.zip  | https://864000.lanzouj.com/iyiSZ28yfi1e | 6.6 M    |
| Linux                  | wxdown-1.0.4-linux-amd64.zip   | https://864000.lanzouj.com/idCnP28yfbef | 7.0 M    |
| macOS                  | wxdown-1.0.4-linux-arm64.zip   | https://864000.lanzouj.com/igmYJ28yfcpc | 6.6 M    |
| macOS（Apple Silicon） | wxdown-1.0.4-darwin-amd64.zip  | https://864000.lanzouj.com/iom6a28yfgwd | 6.9 M    |

### v1.0.3

替换根目录 `web`、`wxdown`、`config.yaml` 文件重新启动完成更新

- 增加自定义目录名称
- 增加启动时默认在浏览器打开管理端
- 优化公众号最近更新样式问题
- 优化html、图片异步下载（可能会出现页面下载完了，图片还没有过会就好）
- 优化管理页面
- 重构代码
- 

| 操作系统               | 文件名                         | 链接                                    | 文件大小 |
| ---------------------- | ------------------------------ | --------------------------------------- | -------- |
| Windows                | wxdown-1.0.3-windows-amd64.zip | https://864000.lanzouj.com/iR2ZM24ur81i | 18.1 M   |
| ARM Linux              | wxdown-1.0.3-linux-arm64.zip   | https://864000.lanzouj.com/iq5LZ24ur59i | 6.0 M    |
| Linux                  | wxdown-1.0.3-linux-amd64.zip   | https://864000.lanzouj.com/iDzcY24ur4be | 6.3 M    |
| macOS                  | wxdown-1.0.3-darwin-arm64.zip  | https://864000.lanzouj.com/i34G224ur3cj | 6.0 M    |
| macOS（Apple Silicon） | wxdown-1.0.3-darwin-amd64.zip  | https://864000.lanzouj.com/iwuPf24ur2ef | 6.3 M    |


### v1.0.2

替换根目录 `web`、`wxdown`、`config.yaml` 文件重新启动完成更新

- 增加自定义目录名称（合集）
- 增加 macOS 打开目录
- 增加 svg 内嵌图片下载
- 增加 http 基础认证（详情配置文件）
- 修改 svg 文件下载 bug
- 移除自动检测缺少目录上个版本

| 操作系统                 | 版本/架构                                                                     | 大小     |
|----------------------|---------------------------------------------------------------------------|--------|
| Windows              | [wxdown-1.0.2-windows-amd64.exe](https://864000.lanzouj.com/i6flE21a198h) | 9.23MB |
| ARM Linux            | [wxdown-1.0.2-linux-arm64](https://864000.lanzouj.com/irG1o21a17bi)       | 8.75MB |
| Linux                | [wxdown-1.0.2-linux-amd64](https://864000.lanzouj.com/iQlpa21a170h)       | 9.04MB |
| macOS                | [wxdown-1.0.2-darwin-amd64](https://864000.lanzouj.com/itt0z21a19sh)      | 9.09MB |
| macOS（Apple Silicon） | [wxdown-1.0.2-darwin-arm64](https://864000.lanzouj.com/imOXU21a1a9e)      | 8.75MB |

### v1.0.1

- 增加音频、视频下载
- 增加首页合集、标签合集下载
- 优化页面样式、自适应移动端
- 部分图片解析异常bug优化
- 替换根目录的 `web` 目录和 `wxdown` 开头可执行文件，重新启动完成更新

| 操作系统                 | 版本/架构                                                                     | 大小     |
|----------------------|---------------------------------------------------------------------------|--------|
| Windows              | [wxdown-1.0.1-windows-amd64.exe](https://864000.lanzouj.com/ihbIY1zvivwj) | 9.23MB |
| ARM Linux            | [wxdown-1.0.1-linux-arm64](https://864000.lanzouj.com/iq50m1zvivjg)       | 8.75MB |
| Linux                | [wxdown-1.0.1-linux-amd64](https://864000.lanzouj.com/iaOY71zviveb)       | 9.04MB |
| macOS                | [wxdown-1.0.1-darwin-amd64](https://864000.lanzouj.com/ittpb1zviv0h)      | 9.09MB |
| macOS（Apple Silicon） | [wxdown-1.0.1-darwin-arm64](https://864000.lanzouj.com/itnUE1zviv6d)      | 8.75MB |

### v1.0.0

| 操作系统                 | 版本/架构                                                                     | 大小     |
|----------------------|---------------------------------------------------------------------------|--------|
| Windows              | [wxdown-1.0.0-windows-amd64.exe](https://864000.lanzouj.com/i5JaN1z84hlc) | 9.23MB |
| ARM Linux            | [wxdown-1.0.0-linux-arm64](https://864000.lanzouj.com/iueMH1z84e9c)       | 8.75MB |
| Linux                | [wxdown-1.0.0-linux-amd64](https://864000.lanzouj.com/iueMH1z84e9c)       | 9.04MB |
| macOS                | [wxdown-1.0.0-darwin-amd64](https://864000.lanzouj.com/ikGp81z84i7e)      | 9.09MB |
| macOS（Apple Silicon） | [wxdown-1.0.0-darwin-arm64](https://864000.lanzouj.com/izIaY1z84iuh)      | 8.75MB |

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
- `Dockerfile`: 如果你的项目需要使用 Docker 镜像进行部署，可以在此处编写 Dockerfile。
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

本软件使用 Go 语言开发，欢迎开发者贡献代码或提出改进建议。请在 [GitHub](https://github.com/systemmin/wxdown/issues) 上提交
issue 或 pull request。

## 注意事项

**注意：** 请确保遵循网站的使用条款和政策。自行承担风险。

- 本软件完全免费。
- 请勿传播未经授权的文章或图片。
- 使用软件造成的影响由使用者承担。
- 仅供学习交流，严禁用于商业用途，请勿传播下载的数据。
- 本脚本所获取的资源完全合法，与浏览器能直接获得的资源一致。
- 在保存、转换文章或管理图片时，请注意版权和法律规定。

