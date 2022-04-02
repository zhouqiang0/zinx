package ziface

//消息管理抽象层

type IMsgHandle interface {
	// DoMsGHandler 调度/执行对应的Router消息处理方法
	DoMsGHandler(request IRequest)

	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)

	// StartWorkerPool 启动一个Worker工作池, 开启工作池的行为只有一次，一个zinx框架全局只有一个worker池
	StartWorkerPool()

	// SendMsgToTaskQueue 将消息交给TaskQueue，由worker来处理
	SendMsgToTaskQueue(request IRequest)
}
