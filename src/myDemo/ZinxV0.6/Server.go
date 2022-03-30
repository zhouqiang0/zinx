package main

import (
	"fmt"
	"zinx/src/zinx/ziface"
	"zinx/src/zinx/znet"
)

//基于Zinx框架来开发的服务端应用程序

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
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
	znet.BaseRouter
}

// Handle Test Handle
func (h *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Hello Handle")
	//先读取客户端数据，再写回hello...
	fmt.Println("recv from client: msgID: ", request.GetMsgId(),
		", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello welcome to zinx-v0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.5]")

	//3 给当前zinx框架添加自定义router
	s.AddRouter(0, &PingRouter{})

	s.AddRouter(1, &HelloRouter{})

	//2 启动server
	s.Serve()
}
