package ziface

// IConnManager 连接管理模块抽象层
type IConnManager interface {
	// Add 添加连接
	Add(conn IConnection)

	// Remove 删除连接
	Remove(conn IConnection)

	// Get 根据connID获取连接
	Get(connID uint32) (IConnection, error)

	// ConnNum 得到当前连接总数
	ConnNum() int

	// ClearConn 清楚所有连接
	ClearConn()
}
