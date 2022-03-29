package ziface

//IRequest接口 把客户端请求的连接信息，和请求的数据包装到一个Request中

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection

	// GetData 得到请求的消息数据
	GetData() []byte

	// GetMsgId 得到请求的消息ID
	GetMsgId() uint32
}
