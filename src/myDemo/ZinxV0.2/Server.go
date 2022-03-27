package main

import "zinx/src/zinx/znet"

//基于Zinx框架来开发的服务端应用程序

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.2]")
	//2 启动server
	s.Serve()
}
