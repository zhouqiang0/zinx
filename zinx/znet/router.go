package znet

import (
	"zinx/zinx/ziface"
)

//实现router时，先嵌入这个BaseRouter基类，根据需求对这个基类的方法进行重写就好

type BaseRouter struct{}

//BaseRouter方法都为空
//是因为有写Router不希望有PreHandle 或PostHandle这两个业务
//BaseRouter仅仅为实现IRouter接口的过渡
//其他router继承BaseRouter, 仅需重写自己所需的方法

func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) Handle(request ziface.IRequest) {}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
