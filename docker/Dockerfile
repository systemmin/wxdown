# Centos 官方镜像
FROM quay.io/centos/centos:stream9

# 设置 LANG 和 LC_ALL 环境变量为 en_US.UTF-8
ENV LANG en_US.UTF-8
ENV LC_ALL en_US.UTF-8

# 设置工作目录
WORKDIR /wx

# 更新包列表并安装依赖项、X 环境和字体，同时清理缓存
RUN dnf -y update && \
    dnf -y install \
    epel-release \
    wget \
    xorg-x11-server-Xvfb \
    xorg-x11-fonts-75dpi \
    xorg-x11-fonts-Type1 \
    libX11 \
    libXext \
    libXrender \
    libjpeg \
    openssl && \
    dnf clean all && \
    # 安装并生成 locale
    dnf install -y glibc-locale-source glibc-langpack-en && \
    localedef -i en_US -f UTF-8 en_US.UTF-8 && \
    # 清理无用的临时文件
    rm -rf /var/cache/dnf

# 将本地下载的 wkhtmltopdf 包复制到容器中
COPY wkhtmltox-0.12.6.1-2.almalinux9.x86_64.rpm /tmp/wkhtmltox.rpm

# 安装 wkhtmltopdf 包
RUN rpm -Uvh /tmp/wkhtmltox.rpm && \
    rm /tmp/wkhtmltox.rpm

# 复制字体到容器中，解决中文乱码问题
COPY ./fonts/ /usr/share/fonts/

# 刷新字体缓存
RUN fc-cache -fv

# 将 Go 程序的二进制文件复制到容器中
COPY ./wx /wx

# 给 Go 程序二进制文件添加执行权限
RUN chmod +x /wx/wxdown

# 暴露端口
EXPOSE 81

# 设置启动命令
CMD ["./wxdown"]
