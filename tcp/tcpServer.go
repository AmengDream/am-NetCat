package tcp

/*
TCP服务端程序的处理流程：
1.监听端口
2.等待客户端请求建立连接
3.创建goroutine处理连接
4.关闭连接
*/

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os/exec"
)

// Exec 函数是一个处理函数，处理连接的命令执行
func Exec(conn net.Conn) {
	//defer conn.Close()

	//函数Command返回一个*Cmd，用于使用给出的参数执行name指定的程序。只设定了Path和Args两个参数。并返回一个指针类型的Cmd实例。
	shell := exec.Command("/bin/bash", "-i")

	//把Cmd实例的标准输入和标准输出重定向到conn,这样conn的输入输出会被当做程序的输入输出
	shell.Stdout = conn
	shell.Stdin = conn

	err := shell.Run() //Run启动指定的命令并等待它完成
	if err != nil {
		e := fmt.Sprintln("命令运行失败!err:", err)
		conn.Write([]byte(e))
	}
}

// Connection 是一个处理函数。负责处理连接的交互
func Connection(conn net.Conn) {
	//4.关闭连接
	defer conn.Close()

	addr := conn.RemoteAddr()
	fmt.Printf("成功连接!%s\n", addr)

	for {
		re := bufio.NewReader(conn)      //创建从连接中读取的对象
		str, rerr := re.ReadString('\n') //读到换行为止
		if rerr == io.EOF {
			fmt.Print("客户端")
			goto exit
		} else if rerr != nil {
			fmt.Printf("从客户端%s读取信息发生错误,err:%s\n", addr, rerr)
			goto exit
		} else if str == "" {
			continue
		}

		//这里的判断条件加'exit\n'是发现content[:n]转换成字符串会包含最后的换行符
		if str == "exit\n" {
			goto exit
		}

		fmt.Printf("接收到%s信息:\n%s\n", addr, str)

		//发送数据
		wr := bufio.NewWriter(conn)
		if str != "" {
			_, werr := wr.WriteString("ok,服务器已收到您的信息\n")
			if werr != nil {
				fmt.Printf("向客户端%s发送信息失败,err:%s", addr, werr)
			}
		}

		//判断是否要执行命令
		if str == "shell\n" {
			// _, werr := wr.WriteString("切换到shell模式\n")

			// if werr != nil {
			// 	fmt.Printf("向客户端%s发送信息失败,err:%s", addr, werr)
			// }
			conn.Write([]byte("进入命令执行模式\n"))
			Exec(conn)
			goto exit
		}

		wr.Flush()

	}
exit:
	fmt.Printf("%s断开连接,服务正在等待新的连接...\n", addr)

}

// Tcplisten 函数建立一个TCP的监听并接受客户端的连接
func Tcplisten(address string) {

	//1.创建一个监听对象listen
	listen, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("创建监听失败,err:", err)
		return
	} else {
		fmt.Printf("创建监听成功!正在监听%s\n", address)
	}

	//2.等待连接
	for {
		conn, err := listen.Accept()

		if err != nil {
			fmt.Println("建立连接错误,err:", err)
			continue
		} else {
			conn.Write([]byte("连接成功!输入shell可进入命令执行模式\n"))
			go Connection(conn) //3.启动一个goroutine处理后面的交互
		}

	}

}
