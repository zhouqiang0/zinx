package main

import (
	"fmt"
	"io"
	"net"
	"time"
	znet2 "zinx/zinx/znet"
)

//模拟客户端
func main() {
	fmt.Println("client start...")
	time.Sleep(2 * time.Second)

	//1 直接连接远程服务器，得到一个conn
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		//发送封包的的message消息
		dp := znet2.NewDataPack()
		binaryMsg, err := dp.Pack(znet2.NewMsgPackage(0, []byte("zinxV0.8 client0 Test Message")))
		if err != nil {
			fmt.Println("Pack error: ", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("client write error: ", err)
			return
		}

		//服务器应该回复message id:1, ping...ping...ping...
		//读取head
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error: ", err)
			break
		}
		//将head解包为message结构体
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil { //msgHead: msgId|msgLen
			fmt.Println("client unpack head error: ", err)
			break
		}

		//读取data(ping...ping...ping...)
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet2.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error: ", err)
				return
			}

			fmt.Println("----------->Recv Server Msg : ID= ", msg.Id, ", len= ", msg.DataLen,
				", data= ", string(msg.Data), "<------------")
		}

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}

}
