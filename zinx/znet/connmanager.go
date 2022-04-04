package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/zinx/ziface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁

}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护map, 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//添加conn
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID= ", conn.GetConnID(), " added to ConnManager successfully: conn num = ", connMgr.ConnNum())
}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护map, 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID= ", conn.GetConnID(), " removed from ConnManager successfully: conn num = ", connMgr.ConnNum())
}

func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护map, 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connID NOT FOUND！")
	}

}

func (connMgr *ConnManager) ConnNum() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	//保护map, 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除并停止conn的工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()

		//删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections success! conn num= ", connMgr.ConnNum())
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}
