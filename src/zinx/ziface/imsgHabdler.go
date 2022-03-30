package ziface

//消息管理抽象层

type IMsgHandle interface {
	// DoMsGHandler 调度/执行对应的Router消息处理方法
	DoMsGHandler(request IRequest)

	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}
