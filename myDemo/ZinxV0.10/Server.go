package main

import (
	"fmt"
	ziface2 "zinx/zinx/ziface"
	znet2 "zinx/zinx/znet"
)

//基于Zinx框架来开发的服务端应用程序

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet2.BaseRouter
}

// Handle Test Handle
func (p *PingRouter) Handle(request ziface2.IRequest) {
	fmt.Println("Call Router Handle")
	//先读取客户端数据，再写回ping...ping...ping...
	fmt.Println("recv from client: msgID: ", request.GetMsgId(),
		", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

// HelloRouter  test 自定义路由
type HelloRouter struct {
	znet2.BaseRouter
}

// Handle Test Handle
func (h *HelloRouter) Handle(request ziface2.IRequest) {
	fmt.Println("Call Hello Handle")
	//先读取客户端数据，再写回hello...
	fmt.Println("recv from client: msgID: ", request.GetMsgId(),
		", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello welcome to zinx-v0.9"))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnAfterStart 创建建立连接后的钩子函数
func DoConnAfterStart(conn ziface2.IConnection) {
	fmt.Println("------------>建立连接后的钩子函数开始执行。。。")
	if err := conn.SendMsg(202, []byte("DoConnAfterStart BEGIN")); err != nil {
		fmt.Println(err)
	}

	//给连接设置一些属性
	fmt.Println(" Set Connection Name, Home, Hoe ...")
	conn.SetProperty("Name", "ZhouQ - Gopher")
	conn.SetProperty("GitHub", "https://github.com/zhouqiang0")
}

// DoConnBeforeStop 创建关闭连接前的钩子函数
func DoConnBeforeStop(conn ziface2.IConnection) {
	fmt.Println("------------>关闭连接前的钩子函数开始执行。。。")
	fmt.Println("conn ID = ", conn.GetConnID(), " is Lost...")

	//获取连接属性
	name, _ := conn.GetProperty("Name")
	fmt.Println(name)

	github, _ := conn.GetProperty("GitHub")
	fmt.Println(github)
}

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet2.NewServer()

	//2 注册链接Hook函数
	s.SetAfterConnStart(DoConnAfterStart)
	s.SetBeforeConnStop(DoConnBeforeStop)

	//3 给当前zinx框架添加自定义router
	s.AddRouter(0, &PingRouter{})

	s.AddRouter(1, &HelloRouter{})

	//2 启动server
	s.Serve()
}
