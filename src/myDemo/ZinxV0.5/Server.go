package main

import (
	"fmt"
	"zinx/src/zinx/ziface"
	"zinx/src/zinx/znet"
)

//基于Zinx框架来开发的服务端应用程序

// PingRRouter ping test 自定义路由
type PingRRouter struct {
	znet.BaseRouter
}

//// PreHandle Test PreHandle
//func (p *PingRRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
//	if err != nil {
//		fmt.Println("Call back before ping err ", err)
//	}
//}

// Handle Test Handle
func (p *PingRRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	//先读取客户端数据，再写回ping...ping...ping...
	fmt.Println("recv from client: msgID: ", request.GetMsgId(),
		", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

//// PostHandle Test PreHandle
//func (p *PingRRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PostHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
//	if err != nil {
//		fmt.Println("Call back after ping err ", err)
//	}
//}

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.5]")

	//3 给当前zinx框架添加自定义router
	s.AddRouter(&PingRRouter{})

	//2 启动server
	s.Serve()
}
