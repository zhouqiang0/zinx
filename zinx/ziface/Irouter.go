package ziface

//路由抽象接口， 路由的数据都是IRequest

type IRouter interface {
	// PreHandle 在处理conn业务之前的方法Hook
	PreHandle(request IRequest)

	// Handle 在处理conn业务的主方法Hook
	Handle(request IRequest)

	// PostHandle 在处理conn业务之后的方法Hook
	PostHandle(request IRequest)
}
