package znet

import (
	"fmt"
	"net"
	"zinx/zinx/utils"
	ziface2 "zinx/zinx/ziface"
)

// Server 是IServer 的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int

	////当前的Server添加一个router，server注册的连接对应的处理业务
	//Router ziface.IRouter

	//当前Server的消息管理模块，用于绑定MsgID和对应的handler
	MsgHandler ziface2.IMsgHandle

	//该Server的连接管理器
	ConnMgr ziface2.IConnManager

	//该Server创建连接之后自动调用的Hook函数
	AfterConnStart func(conn ziface2.IConnection)

	//该Server销毁连接之前自动调用的Hook函数
	BeforeConnStop func(conn ziface2.IConnection)
}

func (s *Server) AddRouter(msgID uint32, router ziface2.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success!")
}

//// CallBackToClient 定义当前客户端的所绑定的handleAPI
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	//回显业务
//	fmt.Println("[Conn Handle] CallBackToClient ... ")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back buf err ", err)
//		return errors.New("CallBackToClient error")
//	}
//	return nil
//}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name : %s, listening at IP : %s, Port : %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn : %d, MaxPackageSize : %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	fmt.Printf("[Start] Server Listenner at IP : %s, Port : %d, is starting\n", s.IP, s.Port)

	go func() {
		//0 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
		}

		//2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}

		fmt.Println("start Zinx server success, ", s.Name, " success, Listening...")

		var cid uint32
		cid = 0
		//3 阻塞的等待客户端链接，处理客户端链接业务（读写）
		for {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//设置最大连接个数的判断，超过则关闭此次连接
			if s.ConnMgr.ConnNum() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("----------->Too Many Connections!!!<-----------MaxConn= ", utils.GlobalObject.MaxConn)
				_ = conn.Close()
				continue
			}

			//将处理新连接的业务方法与conn进行绑定 得到dealConn
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//启动当前的连接业务处理
			go dealConn.Start()

		}
	}()

}

func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止/回收
	fmt.Println("[STOP] Zinx server name = ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

func (s *Server) GetConnMgr() ziface2.IConnManager {
	return s.ConnMgr
}

// NewServer 初始化Server的方法
func NewServer() ziface2.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// SetAfterConnStart 该Server销毁连接之前自动调用的Hook函数
func (s *Server) SetAfterConnStart(hookFunc func(conn ziface2.IConnection)) {
	s.AfterConnStart = hookFunc
}

// SetBeforeConnStop 注册BeforeConnStop()钩子函数
func (s *Server) SetBeforeConnStop(hookFunc func(conn ziface2.IConnection)) {
	s.BeforeConnStop = hookFunc
}

// CallAfterConnStart 调用AfterConnStart()钩子函数的方法
func (s *Server) CallAfterConnStart(conn ziface2.IConnection) {
	if s.AfterConnStart != nil {
		fmt.Println("--------->call AfterConnStart()...")
		s.AfterConnStart(conn)
	}
}

// CallBeforeConnStop 调用BeforeConnStop()钩子函数的方法
func (s *Server) CallBeforeConnStop(conn ziface2.IConnection) {
	if s.BeforeConnStop != nil {
		fmt.Println("--------->call BeforeConnStop()...")
		s.BeforeConnStop(conn)
	}
}
