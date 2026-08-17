#!/usr/bin/env bash
# 发布脚本：交叉编译 + 打包运行时资源
set -euo pipefail

# 检查参数
if [[ "$#" -ne 1 ]]; then
    echo "usage: $0 <version>"
    echo "  e.g. $0 1.1.3   (tag v1.1.3 也可，会自动去掉前缀 v)"
    exit 1
fi

# 检查main.go文件
if [[ ! -f "main.go" ]]; then
    echo "file main.go does not exists"
    exit 2
fi

# 整理依赖
go mod tidy

# 版本号（去掉 tag 前缀 v）
VERSION="${1#v}"
# 可执行文件名
BIN="wxdown"
# 输出目录
WORKDIR="release_${VERSION}"

rm -rf "$WORKDIR"
mkdir -p "$WORKDIR"

# 目标平台
OS=(  "linux"  "linux"  "darwin"  "darwin"  "windows"  "windows")
ARCH=("amd64"  "arm64"  "arm64"   "amd64"   "arm64"    "amd64")
# 归档后缀名（同时作为包内目录名）
LABEL=("Linux_x64" "Linux_arm64" "macOS_arm64" "macOS_x64" "Windows_arm64" "Windows_x64")

MAX_INDEX=$((${#OS[@]} - 1))

# 打包目录为 zip：优先用 zip，没有则回退到 .NET ZipFile（Windows git bash）
# 用法：zip_dir <workdir> <archive_name> <source_basename>
zip_dir() {
    local workdir="$1" archive="$2" src="$3"
    if command -v zip >/dev/null 2>&1; then
        (cd "$workdir" && zip -qr "$archive" "$src")
        return
    fi
    local win_src win_dst
    win_src=$(cygpath -w "$workdir/$src" 2>/dev/null || echo "$workdir/$src")
    win_dst=$(cygpath -w "$workdir/$archive" 2>/dev/null || echo "$workdir/$archive")
    # 带重试：Windows Defender 等可能短暂占用刚拷入的大文件
    ZIP_SRC="$win_src" ZIP_DST="$win_dst" powershell -NoProfile -Command '
        $ErrorActionPreference = "Stop"
        Add-Type -AssemblyName System.IO.Compression.FileSystem
        $src = $env:ZIP_SRC; $dst = $env:ZIP_DST
        if (Test-Path $dst) { Remove-Item $dst -Force }
        for ($i = 0; $i -lt 5; $i++) {
            try {
                [System.IO.Compression.ZipFile]::CreateFromDirectory($src, $dst, [System.IO.Compression.CompressionLevel]::Optimal, $true)
                exit 0
            } catch {
                if ($i -eq 4) { Write-Error $_; exit 1 }
                Start-Sleep -Milliseconds 500
            }
        }
    '
}

for i in $(seq 0 "$MAX_INDEX"); do
    goos="${OS[$i]}"
    goarch="${ARCH[$i]}"
    label="${LABEL[$i]}"
    stagedir="${WORKDIR}/${BIN}_${VERSION}_${label}"

    echo "==> building ${goos}/${goarch} -> ${stagedir}"
    mkdir -p "$stagedir"

    binname="${BIN}"
    if [[ "$goos" == "windows" ]]; then
        binname="${BIN}.exe"
    fi

    # 编译（CGO_ENABLED=0 保证静态链接）
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
        go build -trimpath \
        -ldflags "-s -w -X main.runMode=binary -X main.version=${VERSION}" \
        -o "${stagedir}/${binname}" main.go

    # 拷贝运行时资源
    cp -R web "$stagedir/"
    cp config.yaml "$stagedir/"
    cp -R certs "$stagedir/"

    # Windows 自带 wkhtmltopdf
    if [[ "$goos" == "windows" ]]; then
        cp docker/wkhtmltopdf.exe "$stagedir/"
    fi

    # 打包
    if [[ "$goos" == "windows" ]]; then
        archive="${BIN}_${VERSION}_${label}.zip"
        zip_dir "$WORKDIR" "$archive" "$(basename "$stagedir")"
    else
        archive="${BIN}_${VERSION}_${label}.tar.gz"
        tar -czf "${WORKDIR}/${archive}" -C "$WORKDIR" "$(basename "$stagedir")"
    fi

    # 清理暂存目录，只留归档
    rm -rf "$stagedir"
    echo "    -> ${WORKDIR}/${archive}"
done

# 生成 sha256 校验和
(
    cd "$WORKDIR"
    : > checksums.txt
    for f in *.zip *.tar.gz; do
        [[ -e "$f" ]] || continue
        if command -v sha256sum >/dev/null 2>&1; then
            sha256sum "$f" >> checksums.txt
        else
            shasum -a 256 "$f" >> checksums.txt
        fi
    done
)

echo ""
echo "Build completed. Artifacts in ${WORKDIR}/"
ls -lh "$WORKDIR"
