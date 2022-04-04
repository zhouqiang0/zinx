package ziface

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动服务器
	Start()

	// Stop 停止服务器
	Stop()

	// Serve 运行服务器
	Serve()

	// AddRouter 路由功能，给当前的服务注册一个路由方法，供客户端的连接处理
	AddRouter(msgID uint32, router IRouter)

	// GetConnMgr 获取当前的连接管理器
	GetConnMgr() IConnManager

	// SetAfterConnStart 该Server销毁连接之前自动调用的Hook函数
	SetAfterConnStart(func(connection IConnection))

	// SetBeforeConnStop 注册BeforeConnStop()钩子函数
	SetBeforeConnStop(func(connection IConnection))

	// CallAfterConnStart 调用AfterConnStart()钩子函数的方法
	CallAfterConnStart(connection IConnection)

	// CallBeforeConnStop 调用BeforeConnStop()钩子函数的方法
	CallBeforeConnStop(connection IConnection)
}
