package main

import (
	"flag"
	"fmt"
	"nc/tcp"
	"strconv"
	"strings"
)

/*
程序入口
实现功能：
1.接受命令行参数并解析
2.调用相应的功能
*/

// ports 处理端口参数并调用扫描函数
func ports(port string) (star_port, end_port int) {
	//判断port参数是否包含-
	var a, b int
	if find := strings.Contains(port, "-"); find {
		arr := strings.Split(port, "-")
		n := len(arr)
		star, err := strconv.Atoi(arr[0])
		if err != nil {
			fmt.Println("star端口转换出错,err:\n", err)
		}
		end, err := strconv.Atoi(arr[n-1])
		if err != nil {
			fmt.Println("end端口转换出错,err:\n", err)
		}
		a, b = star, end

	} else {
		p, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("端口转换出错,err:\n", err)
		}
		a, b = p, p
		
	}
	return a, b
}

func main() {
	var mode, ip, port string
	var thread int

	//注册参数
	flag.StringVar(&mode, "m", "l", "模式: / l 监听模式 / s 扫描端口 /")
	flag.StringVar(&ip, "ip", "127.0.0.1", "IP地址")
	flag.StringVar(&port, "p", "80", "端口:支持格式：80，1-655535")
	flag.IntVar(&thread, "t", 20, "线程数:如果执行扫描端口模式，可使用-t指定线程数，默认20")

	//解析命令行参数
	flag.Parse()

	//判断用户参数值，调用对应的函数
	switch mode {
	case "l":
		addr := ip + ":" + port
		fmt.Println("正在启动TCP监听...")
		tcp.Tcplisten(addr)
	case "s":
		star_port, end_port := ports(port)
		fmt.Println("正在启动TCP扫描端口...")
		tcp.ScanTCP(thread, star_port, end_port, ip)
	default:
		fmt.Println("模式参数错误!请输入“-h”查看帮助信息")
	}
}
