### v1.0.8

替换所有文件（不包含 `data`）

- 新增 docker 镜像
- 新增浏览器剪切板监听，Chrome 浏览器获取剪切板内容需 https 协议，详见配置文件 [config.yaml#L15](https://github.com/systemmin/wxdown/blob/master/config.yaml#L15)
- 新增 `https` 证书配置，自签证书浏览器会提示不安全连接，可 [手动安装证书](https://dtking.cn/blog/custom-ssl-v3/#4-dao3-ru4-zheng4-shu1-dao4-liu2-lan3-qi4) 文件 `certs/certificate.crt`
- 新增下载重复文章跳过功能，详见配置文件 [config.yaml#L19](https://github.com/systemmin/wxdown/blob/master/config.yaml#L19)
- 优化公网部署转PDF加载慢，转PDF服务地址统一改 `127.0.0.1/wx/....`
- 优化 html 转 PDF 问题

[下载地址](https://864000.lanzouj.com/b0hchwr3e) 密码：1gn6

### v1.0.7

替换根目录 `wxdown`、`web` 、`config.yaml`文件重新启动完成更新

- 增加图片转 `base64` 详见 `config.yaml`
- 优化合集下载缺少部分文章问题
- 优化图片集部分无法下载图片问题

| 操作系统        | 链接                                                                                    | 文件大小   |
|-------------|---------------------------------------------------------------------------------------|--------|
| Windows     | [wxdown-1.0.7-windows-amd64.zip](https://864000.lanzouj.com/icw432ce30if "Windows")   | 18.6 M |
| Linux       | [wxdown-1.0.7-linux-amd64.zip](https://864000.lanzouj.com/ivusg2ce337c "Linux")       | 6.8 M  |
| ARM Linux   | [wxdown-1.0.7-linux-arm64.zip](https://864000.lanzouj.com/iHwXY2ce2ygb "ARM Linux")   | 6.4 M  |
| macOS       | [wxdown-1.0.7-darwin-amd64.zip](https://864000.lanzouj.com/i7xyx2ce31ve "macOS")      | 6.5 M  |
| macOS（M 系列） | [wxdown-1.0.7-darwin-arm64.zip](https://864000.lanzouj.com/iXT7l2ce32mb "macOS M 系列") | 6.8 M  |


### v1.0.6

替换根目录 `wxdown` 文件重新启动完成更新

- 向下兼容 [视频合集](https://mp.weixin.qq.com/s?__biz=MjM5NzQ4OTkwNg==&mid=2247534412&idx=6&sn=edeb6c07b40583a0e5f37253660137a4&chksm=a6db3ad191acb3c78e3af57dca5432d211f488a7c8d7db3d58c5a69b9152840d1944bf2afd92&scene=178)
  无法下载问题（2021年前后文章）
- 只做视频处理，页面无任何处理。
- 图片无法预览问题，点击下载 [web](https://864000.lanzouj.com/i3LbK2ab74sd "Pages")，替换根目录 `web` 文件夹

注意事项：

1. 没有该视频需求的无需更新
2. 出现以下错误时，删除根目录 `wkhtmltopdf.exe`（单独安装了 wkhtmltopdf 会出现以下问题）。
   ```shell
    exec: "wkhtmltopdf": cannot run executable found relative to current directory
    ```

| 操作系统        | 链接                                                                                    | 文件大小   |
|-------------|---------------------------------------------------------------------------------------|--------|
| Windows     | [wxdown-1.0.6-windows-amd64.zip](https://864000.lanzouj.com/idWga29zu25e "Windows")   | 18.6 M |
| Linux       | [wxdown-1.0.6-linux-amd64.zip](https://864000.lanzouj.com/imv9V29zu00h "Linux")       | 6.8 M  |
| ARM Linux   | [wxdown-1.0.6-linux-arm64.zip](https://864000.lanzouj.com/i8Ef529zu0mj "ARM Linux")   | 6.4 M  |
| macOS       | [wxdown-1.0.6-darwin-amd64.zip](https://864000.lanzouj.com/i52k629zu2kj "macOS")      | 6.5 M  |
| macOS（M 系列） | [wxdown-1.0.6-darwin-arm64.zip](https://864000.lanzouj.com/igIiu29ztzgh "macOS M 系列") | 6.8 M  |

### v1.0.5

替换根目录 `web`、`wxdown` 文件重新启动完成更新

- 解决有空格的文件夹路径问题
    - 生成 PDF 失败 [issues#4](https://github.com/systemmin/wxdown/issues/5)
    - 无法通过浏览器打开文件夹
- 解决首页合集无法下载 [52破解](https://www.52pojie.cn/forum.php?mod=viewthread&tid=1960591&page=8#pid51251321)
- 新增图片预览功能

注意事项：

1. 不要安装在带有空格的文件路径中，避免不必要的麻烦（该问题已解决）
2. window 下无法转 PDF 据部分[使用者反馈](https://www.52pojie.cn/forum.php?mod=redirect&goto=findpost&ptid=1960591&pid=51247865)
   可以通过双击 `wkhtmltopdf.exe`

| 操作系统        | 链接                                                                                    | 文件大小   |
|-------------|---------------------------------------------------------------------------------------|--------|
| Windows     | [wxdown-1.0.5-windows-amd64.zip](https://864000.lanzouj.com/ifaG229hlsad "Windows")   | 18.6 M |
| Linux       | [wxdown-1.0.5-linux-amd64.zip](https://864000.lanzouj.com/i0utt29hlqpg "Linux")       | 6.8 M  |
| ARM Linux   | [wxdown-1.0.5-linux-arm64.zip](https://864000.lanzouj.com/iaDN929hlr2j "ARM Linux")   | 6.4 M  |
| macOS       | [wxdown-1.0.5-darwin-amd64.zip](https://864000.lanzouj.com/iQG4029hlsqj "macOS")      | 6.5 M  |
| macOS（M 系列） | [wxdown-1.0.5-darwin-arm64.zip](https://864000.lanzouj.com/iC0XT29hlt9i "macOS M 系列") | 6.8 M  |

### v1.0.4

替换根目录 `web`、`wxdown` 文件重新启动完成更新

- 解决 [webp](https://mp.weixin.qq.com/s/_eeCF9JLOKF-YeojxsTmog) 图片格式无法转PDF问题
- 解决标签合集不满足分页条件错误问题 [issues#4](https://github.com/systemmin/wxdown/issues/4)
- 新增图片集文章类型下载 [示例](https://mp.weixin.qq.com/s/2E5aiMre-NO0Vw9rnGTkCQ)

| 操作系统                 | 文件名                            | 链接                                      | 文件大小   |
|----------------------|--------------------------------|-----------------------------------------|--------|
| Windows              | wxdown-1.0.4-windows-amd64.zip | https://864000.lanzouj.com/iWeh328yffqb | 18.2 M |
| ARM Linux            | wxdown-1.0.4-darwin-arm64.zip  | https://864000.lanzouj.com/iyiSZ28yfi1e | 6.6 M  |
| Linux                | wxdown-1.0.4-linux-amd64.zip   | https://864000.lanzouj.com/idCnP28yfbef | 7.0 M  |
| macOS                | wxdown-1.0.4-linux-arm64.zip   | https://864000.lanzouj.com/igmYJ28yfcpc | 6.6 M  |
| macOS（Apple Silicon） | wxdown-1.0.4-darwin-amd64.zip  | https://864000.lanzouj.com/iom6a28yfgwd | 6.9 M  |

### v1.0.3

替换根目录 `web`、`wxdown`、`config.yaml` 文件重新启动完成更新

- 增加自定义目录名称
- 增加启动时默认在浏览器打开管理端
- 优化公众号最近更新样式问题
- 优化html、图片异步下载（可能会出现页面下载完了，图片还没有过会就好）
- 优化管理页面
- 重构代码
-

| 操作系统                 | 文件名                            | 链接                                      | 文件大小   |
|----------------------|--------------------------------|-----------------------------------------|--------|
| Windows              | wxdown-1.0.3-windows-amd64.zip | https://864000.lanzouj.com/iR2ZM24ur81i | 18.1 M |
| ARM Linux            | wxdown-1.0.3-linux-arm64.zip   | https://864000.lanzouj.com/iq5LZ24ur59i | 6.0 M  |
| Linux                | wxdown-1.0.3-linux-amd64.zip   | https://864000.lanzouj.com/iDzcY24ur4be | 6.3 M  |
| macOS                | wxdown-1.0.3-darwin-arm64.zip  | https://864000.lanzouj.com/i34G224ur3cj | 6.0 M  |
| macOS（Apple Silicon） | wxdown-1.0.3-darwin-amd64.zip  | https://864000.lanzouj.com/iwuPf24ur2ef | 6.3 M  |

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
