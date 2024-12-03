### 打包镜像

1. 将 `fonts` 目录上传到 `Dockerfile` 同级目录下
2. 将 `wkhtmltox-0.12.6.1-2.almalinux9.x86_64.rpm` 目录上传到 `Dockerfile` 同级目录下

执行打包命令：

```shell
docker build -t wxdown .
```

### 挂载目录

- `data` 数据目录
- `config` 配置文件目录
- `certs` 证书目录

```shell
mkdir -p /home/wxdown/
```

### 启动容器

```shell
 docker run -p 81:81 --name wxdown \
 -v /home/wxdown/data/:/wx/ \
 -d wxdown
 ```

### 访问 81

http://localhost:81 或 http://127.0.0.1:81

`localhost` 可以修改为具体的 `IP` 地址
