##  打包

```shell
docker build -t wxdown .
```

## 启动

- `/wx/data/` 数据目录文件
- `/wx/config/` 配置文件目录，将 `config.yaml` 拷贝至`/home/wxdown/config/`下

1. 创建挂载目录
    ```shell
    mkdir -p /home/wxdown/data/
    mkdir -p /home/wxdown/config/
    ```

2. 启动 docker 容器
    ```shell
    docker run -p 81:81 --name wxdown \
    -v /home/wxdown/data/:/wx/data/ \
    -v /home/wxdown/config/:/wx/config/ \
    -d wxdown
    ```

## 访问

浏览器访问 http://localhost:81 localhost 可以修改为具体的 `IP` 地址

## 扩展

### 进入 wxdown 容器

```shell
docker exec -it wxdown /bin/sh
```

### HTML 转 PDF

命令：`禁用 js` `html 地址` `输出文件路径`

```shell
wkhtmltopdf --debug-javascript http://127.0.0.1 /wx/output.pdf
```

