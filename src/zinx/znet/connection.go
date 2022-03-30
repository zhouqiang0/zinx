package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
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

	////该链接处理的方法
	//Router ziface.IRouter

	//消息的管理模块
	MsgHandler ziface.IMsgHandle
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running... ")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.Conn.RemoteAddr())
	defer c.Stop()

	for {
		//读取客户端数据到buf中，最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err ", err)
		//	continue
		//}
		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的msg Head(8字节)
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg Head error ", err)
			break
		}

		//拆包，得到MsgID, MsgDataLen放在一个msg对象中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		//根据dataLen, 读取Data, 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 { //从head中获取后续要读的长度
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//从路由中找到注册绑定的Conn对应的router调用
		//go func(request ziface.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)
		go c.MsgHandler.DoMsGHandler(&req)

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

// SendMsg 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg ")
	}
	//将data进行封包，得到二进制流binaryMsg
	dp := DataPack{}

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgId)
		return errors.New("Pack error msg ")
	}

	//将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg error id : ", msgId, " error : ", err)
		return errors.New("conn write error ")
	}

	return nil
}

// NewConnection 初始化方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandle,
		ExitChan:   make(chan bool, 1),
	}
	return c
}
