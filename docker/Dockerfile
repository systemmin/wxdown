# 使用一个基础镜像，例如 Ubuntu
FROM centos:7

# 设置 LANG 环境变量为 en_US.UTF-8
ENV LANG en_US.UTF-8

# 设置工作目录
WORKDIR /wx

# 更新包列表并安装依赖项
RUN yum -y update && yum -y install \
    epel-release \
    wget \
    xorg-x11-server-Xvfb \
    xorg-x11-fonts-75dpi \
    xorg-x11-fonts-Type1 \
    libX11 \
    libXext \
    libXrender \
    libjpeg \
    openssl \
    wqy-zenhei-fonts \
    wqy-microhei-fonts \
    wqy-zenhei-fonts \
    && yum clean all

# 将本地下载的 wkhtmltopdf 包复制到容器中
COPY wkhtmltox-0.12.6-1.centos7.x86_64.rpm /tmp/wkhtmltox.rpm

# 安装 wkhtmltopdf 包
RUN rpm -Uvh /tmp/wkhtmltox.rpm

# 删除复制的 rpm 包
RUN rm /tmp/wkhtmltox.rpm

# 将 Go 程序的二进制文件复制到容器中
COPY ./wx /wx

# 给二进制文件执行权限
RUN chmod +x /wx/wxdown

# 暴露端口
EXPOSE 81

# 设置启动命令
CMD ["./wxdown"]
