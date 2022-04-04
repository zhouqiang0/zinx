package ziface

//将请求的消息封装到一个Message中，定义抽象接口

type IMessage interface {
	// GetMsgId 获取消息的ID
	GetMsgId() uint32

	// GetMsgLen 获取消息长度
	GetMsgLen() uint32

	// GetData 获取消息内容
	GetData() []byte

	// SetMsgId 设置消息ID
	SetMsgId(uint32)

	// SetDataLen 设置消息长度
	SetDataLen(uint32)

	// SetData 设置消息内容
	SetData([]byte)
}
