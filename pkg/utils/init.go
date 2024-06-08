package utils

import (
	"fmt"
	"runtime"
	"time"
)

func InitPrint(port string, version string, runMode string, runPath string, resourcePath string) {
	// 定义格式化布局
	currentTime := time.Now()
	// 格式化时间
	formattedTime := currentTime.Format(time.DateTime)
	// 获取 ip 地址
	ips := Ips()
	gang := "----------------------------------------"

	fmt.Println(gang)
	fmt.Println("\t\t欢迎使用 wxdown 工具！")
	fmt.Println(gang)

	fmt.Printf("运行模式 : %s\n", runMode)
	fmt.Printf("软件版本 : %s\n", version)
	fmt.Printf("操作系统 : %s\n", runtime.GOOS)
	fmt.Printf("系统架构 : %s\n", runtime.GOARCH)
	fmt.Printf("启动时间 : %s\n", formattedTime)
	fmt.Println("检测更新 : https://github.com/systemmin/wxdown")

	fmt.Println(gang)
	fmt.Println("\t\t服务信息")
	fmt.Println(gang)
	fmt.Println("服务地址：")

	for _, ip := range ips {
		fmt.Printf("\thttp://%s:%s\t(浏览器访问)\n", ip, port)
	}
	fmt.Println("采集接口：")
	for _, ip := range ips {
		fmt.Printf("\thttp://%s:%s/gather/\t(GET|POST|HEAD)\n", ip, port)
	}

	fmt.Println(gang)
	fmt.Println("\t\t配置信息")
	fmt.Println(gang)
	fmt.Println("运行路径 : " + runPath)
	fmt.Println("资源路径 : " + resourcePath)
}
