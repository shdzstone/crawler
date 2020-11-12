package rpcsupport

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

/*
RPC公用模块
* 生成注册某服务的RPC服务器
* 创建连接某RPC服务器的RPC客户端
*/

//RPC服务器
func ServeRpc(host string, service interface{}) error {
	//首先注册服务
	rpc.Register(service)

	//开启server:监听TCP的host端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", host))
	if err != nil {
		return err
	}

	for {
		//接收传入的连接
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error:%v", err.Error())
			continue
		}

		//使用rpc当场处理掉conn，开启goroutine是为了不让该处理过程阻塞当前goroutine
		go jsonrpc.ServeConn(conn)
	}
	return nil
}

//RPC客户端
func NewClient(host string) (*rpc.Client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%s", host))
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewClient(conn), nil
}
