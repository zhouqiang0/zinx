package znet

import (
	"fmt"
	"net"
	"zinx/src/zinx/utils"
	"zinx/src/zinx/ziface"
)

type Connection struct {
	//socket TCP套接字
	Conn *net.TCPConn

	//连接ID
	ConnID uint32

	//当前连接状态
	isClosed bool

	////当前连接所绑定的业务方法
	//handleAPI ziface.HandleFunc

	//告知当前连接已经退出/停止的channel
	ExitChan chan bool

	//该链接处理的方法
	Router ziface.IRouter
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running... ")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.Conn.RemoteAddr())
	defer c.Stop()

	for {
		//读取客户端数据到buf中，最大512字节
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			continue
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		//从路由中找到注册绑定的Conn对应的router调用
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

// Start 启动连接 让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID= ", c.ConnID)
	//启动当前连接的读数据业务
	go c.StartReader()
	//TODO 启动当前连接写数据的任务
}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID= ", c.ConnID)

	//如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭socket连接
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)
}

// GetTCPConnection 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的 TCP状态、IP、port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据、将数据发送给远程客户端
func (c *Connection) Send(data []byte) error {
	return nil
}

// NewConnection 初始化方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}
