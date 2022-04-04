package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/zinx/ziface"
)

//定义存储有关Zinx框架的全局参数，一些参数可通过zinx.json由用户进行配置

type GlobalObj struct {
	//Server
	TcpServer ziface.IServer //当前Zinx全局的Server对象
	Host      string         //当前服务器主机监听IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器名称

	//Zinx
	Version        string //当前Zinx的版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前Zinx框架数据包的最大值

	WorkerPoolSize uint32 //当前业务工作Worker池的Goroutine数量，worker工作池的消息队列数量
	MaxWorkTaskLen uint32 //Zinx框架中每个worker对应的消息队列允许的任务最大数量
}

// GlobalObject 定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

// Reload 去zinx.json加载用户自定义参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//解析zinx/json
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//初始化GlobalObject
func init() {
	GlobalObject = &GlobalObj{
		TcpServer:      nil,
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.10",
		MaxConn:        100,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkTaskLen: 1024,
	}

	//初始化后，尝试从conf/zinx.json中加载一些用户自定义的参数
	GlobalObject.Reload()
}
