package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	//模拟的服务器
	//1 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}

	//创建一个携程，复杂从客户端处理业务
	go func() {
		//2 从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			go func(conn net.Conn) {
				//处理客户端请求
				//------->拆包过程<--------
				//定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					//第一次从conn读，head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err ", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						//第二次从conn读，根据head中的dataLen 再读取打他内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err ", err)
							return
						}

						//完整的消息已经读取完毕，打印
						fmt.Println("-------->Recv MsgID: ", msg.Id, ", dataLen: ", msg.DataLen, ", data :", string(msg.Data))

					}

				}

			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dail err ", err)
		return
	}

	//创建一个封包对象
	dp := NewDataPack()

	//模拟粘包过程
	//封装msg1
	//msg1 := &Message{
	//	Id:      1,
	//	DataLen: 4,
	//	Data:    []byte("zinx"),
	//}
	sendData1, err := dp.Pack(NewMsgPackage(2, []byte("zinx")))
	if err != nil {
		fmt.Println("client pack msg1 err ", err)
		return
	}

	//封装msg2
	//msg2 := &Message{
	//	Id:      2,
	//	DataLen: 10,
	//	Data:    []byte("hello zinx"),
	//}

	sendData2, err := dp.Pack(NewMsgPackage(2, []byte("hello zinx")))
	if err != nil {
		fmt.Println("client pack msg2 err ", err)
		return
	}
	//一起发送msg1,2
	sendData1 = append(sendData1, sendData2...)

	_, err = conn.Write(sendData1)
	if err != nil {
		fmt.Println("conn write err ", err)
	}

	//客户端阻塞
	select {}
}
