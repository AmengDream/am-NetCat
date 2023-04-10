package tcp

import (
	"fmt"
	"net"
	"sort"
)

// goroutine池来控制goroutine数量
func Works(ports, res chan int, ip string) {

	for port := range ports {
		address := fmt.Sprintf("%s:%d", ip, port)
		//调试代码	fmt.Println(address)

		conn, err := net.Dial("tcp", address)
		//这里还需要加个计时器，如果长时间无响应则err=nil
		if err != nil {
			//连接建立失败，端口关闭或被过滤
			res <- 0
			continue
		}
		conn.Close()
		res <- port
	}
}

//TCP扫描端口：

func ScanTCP(thread, star, end int, ip string) {
	ports := make(chan int, thread) //创建一个缓冲通道
	results := make(chan int)       //创建非缓冲通过跟踪完成情况
	var r []int                     //创建切片从results通道中接收结果

	//创建线程池，控制工作线程数量
	for i := 0; i < cap(ports); i++ {
		go Works(ports, results, ip)
	}
	// 创建任务,通过通道发布任务
	go func() {
		for i := star; i <= end; i++ {
			ports <- i
		}
	}()

	//处理返回结果
	for i := 0; i <= (end - star); i++ {
		p := <-results

		//调试代码 fmt.Println(p)
		if p != 0 {
			r = append(r, p)
		}
	}

	//关闭通道
	close(ports)
	close(results)

	//对结果进行排序处理
	sort.Ints(r)

	//打印结果
	for _, v := range r {
		fmt.Printf("%d open\n", v)
	}
}
