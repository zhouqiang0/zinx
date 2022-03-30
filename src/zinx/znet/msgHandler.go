package znet

import (
	"fmt"
	"strconv"
	"zinx/src/zinx/ziface"
)

//消息处理模块的实现

type MsgHandle struct {
	// msgID对应处理方法的映射
	Apis map[uint32]ziface.IRouter
}

// NewMsgHandle 初始化
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// DoMsGHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsGHandler(request ziface.IRequest) {
	//1 从request中找到MsgID
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api MsgID= ", request.GetMsgId(), " NOT FOUND! please register!")
	}

	//2 调度业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)

}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前ID是否已经绑定
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID= " + strconv.Itoa(int(msgID)))
	}

	//往map添加
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " success!")
}
