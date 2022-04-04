package znet

import (
	"fmt"
	"strconv"
	"zinx/zinx/utils"
	ziface2 "zinx/zinx/ziface"
)

//消息处理模块的实现

type MsgHandle struct {
	// msgID对应处理方法的映射
	Apis map[uint32]ziface2.IRouter

	//负责Worker读取任务的消息队列
	TaskQueue []chan ziface2.IRequest

	//业务工作Worker池的worker数量 == 消息队列的数量
	WorkerPoolSize uint32
}

// NewMsgHandle 初始化
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:      make(map[uint32]ziface2.IRouter),
		TaskQueue: make([]chan ziface2.IRequest, utils.GlobalObject.WorkerPoolSize),

		//从全局配置获取
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

// DoMsGHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsGHandler(request ziface2.IRequest) {
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
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface2.IRouter) {
	// 判断当前ID是否已经绑定
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID= " + strconv.Itoa(int(msgID)))
	}

	//往map添加
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " success!")
}

// StartWorkerPool 启动一个Worker工作池, 开启工作池的行为只有一次，一个zinx框架全局只有一个worker池
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize 分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个Worker被启动
		//1 给当前的worker对应的channel消息队列 开辟空间。第i个worker就要第i个channel
		mh.TaskQueue[i] = make(chan ziface2.IRequest, utils.GlobalObject.MaxWorkTaskLen)

		//2 启动当前Worker，阻塞等待消息从channel中传来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface2.IRequest) {
	fmt.Println("WorkerID= ", workerID, ", is started...")

	//阻塞等待对应channel的消息
	for {
		select {
		//如果有消息进来，出列一个客户端Request，执行该request绑定的handler
		case request := <-taskQueue:
			mh.DoMsGHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由worker来处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface2.IRequest) {
	//1 将消息平均分配给Worker
	//根据客户端建立的ConnID来进行分配，轮询
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID= ", request.GetConnection().GetConnID(),
		", request MsgID= ", request.GetMsgId(), " to WorkerID= ", workerID)

	//2 将消息发送给workerID对应的TaskQueue
	mh.TaskQueue[workerID] <- request
}
