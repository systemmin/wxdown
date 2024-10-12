package main

import (
	"fmt"
	"go-wx-download/config"
	"go-wx-download/internal/controller"
	"go-wx-download/pkg/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// runMode 运行模式,binary 二进制启动 source 源码启动
var runMode = "source"

// version 版本号
var version = "1.0.7"

// LoggingMiddleware 记录每个请求的日志
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()
		//log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		//log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
func ContainsStr(slice []string, item string) bool {
	for _, element := range slice {
		contains := strings.HasPrefix(item, element)
		if contains {
			return contains
		}
	}
	return false
}

// AuthMiddleware 检查每个请求的身份验证
func AuthMiddleware(next http.Handler, config config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			fmt.Println("options")
			// options 直接返回 200
			w.WriteHeader(http.StatusOK)
			return
		}
		path := r.URL.Path
		paths := []string{"/ats", "/gather", "/collect", "/open"}
		if config.Auth.Enable {
			if ContainsStr(paths, path) {
				username, password, ok := r.BasicAuth()
				if ok {
					users := config.Auth.Users
					flag := false
					for _, user := range users {
						split := strings.Split(user, ":")
						if split[0] == username && split[1] == password {
							flag = true
							break
						}
					}
					if !flag {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					} else {
						next.ServeHTTP(w, r)
					}
				} else {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			} else {
				// 在上下文中存储用户信息或其他身份验证数据
				next.ServeHTTP(w, r)

			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func main() {

	// 获取可执行文件路径
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	exPath, err := filepath.EvalSymlinks(ex)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 获取当前工作路径
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前工作路径失败: %s\n", err)
		return
	}
	fmt.Println("cwd：", cwd)

	// 找到当前程序执行根路径
	if runMode == "source" {
		exPath = cwd
	} else {
		exPath = filepath.Join(exPath, "..")
	}

	// 加载配置文件
	cfg := config.LoadConfig(exPath)
	defaultPort := utils.Iif(cfg.Port != "", cfg.Port, "81")
	defaultDataPath := filepath.Join(exPath, utils.Iif(cfg.Path != "" && cfg.Path != "/" && cfg.Path != "./", cfg.Path, "data"))

	// 创建深层次目录，类似 mark -p
	cssPath := filepath.Join(defaultDataPath, "css")
	utils.IsNotExistCreate(cssPath)
	// 创建 task 目录
	utils.IsNotExistCreate(filepath.Join(defaultDataPath, "task"))

	// 将微信样式文件拷贝到 data 目录下
	join := filepath.Join(exPath, "web", "static", "css")
	dir, err := os.ReadDir(join)
	if err != nil {
		log.Println("读取文件失败", err)
	}
	for _, entry := range dir {
		src := filepath.Join(join, entry.Name())
		dst := filepath.Join(cssPath, entry.Name())
		utils.CopyFile(dst, src)
	}

	// 创建一个新的 mux 路由器
	mux := http.NewServeMux()

	// 应用中间件
	handler := LoggingMiddleware(AuthMiddleware(mux, *cfg))

	// 设置根路由
	fs := http.FileServer(http.Dir(filepath.Join(exPath, "web")))
	//http.Handle("/", fs)

	// 静态文件处理器
	mux.Handle("/", fs)

	// 设置公众号文件目录
	wxFs := http.FileServer(http.Dir(defaultDataPath))
	// 使用 wx 前缀
	mux.Handle("/wx/", http.StripPrefix("/wx/", wxFs))

	// 文件操作
	mux.HandleFunc("/ats/", func(w http.ResponseWriter, r *http.Request) {
		controller.Folder(w, r, defaultDataPath)
	})
	mux.HandleFunc("/ats/{folder}/{type}", func(w http.ResponseWriter, r *http.Request) {
		controller.Folder(w, r, defaultDataPath)
	})
	// 打开文件夹
	mux.HandleFunc("/open/", func(writer http.ResponseWriter, request *http.Request) {
		controller.Open(writer, request, defaultDataPath)
	})

	// 单个采集
	mux.HandleFunc("/gather/", func(writer http.ResponseWriter, request *http.Request) {
		controller.Gather(writer, request, cfg, defaultDataPath)
	})
	// 合计采集
	mux.HandleFunc("/collect/", func(writer http.ResponseWriter, request *http.Request) {
		controller.Collect(writer, request, cfg, defaultDataPath)
	})

	// wx 无实际意义加快响应
	mux.HandleFunc("/mp/", func(writer http.ResponseWriter, request *http.Request) {
		_, err2 := fmt.Fprint(writer, "{\"ok\":true}")
		if err2 != nil {
			return
		}
	})
	mux.HandleFunc("/report/", func(writer http.ResponseWriter, request *http.Request) {
		_, err2 := fmt.Fprint(writer, "{\"ok\":true}")
		if err2 != nil {
			return
		}
	})
	mux.HandleFunc("/voice/", func(writer http.ResponseWriter, request *http.Request) {
		_, err2 := fmt.Fprint(writer, "{\"ok\":true}")
		if err2 != nil {
			return
		}
	})

	utils.InitPrint(defaultPort, version, runMode, exPath, defaultDataPath)
	if cfg.Browser {
		utils.OpenBrowser("http://127.0.0.1:" + defaultPort)
	}
	// 指定监听的地址和端口
	addr := ":" + defaultPort

	// 启动服务器
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Printf("无法启动服务器: %s\n", err)
	}
}
