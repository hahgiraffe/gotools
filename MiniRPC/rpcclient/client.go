/*
 * @Author: haha_giraffe
 * @Date: 2019-12-29 21:17:32
 * @Description:客户端的实现
 */
package rpcclient

import (
	"errors"
	"gotools/MiniRPC/transfer"
	"log"
	"net"
	"reflect"
)

//每个Client都有一个
type Client struct {
	conn net.Conn
}

func CreateClient(addr string) *Client {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	return &Client{conn}
}

func (c *Client) Callrpc(rpcname string, fptr interface{}) {
	container := reflect.ValueOf(fptr).Elem()
	//这里是一个内置函数，好处在于可以直接用到container
	f := func(req []reflect.Value) []reflect.Value {
		cReqTransPort := transfer.NewTransport(c.conn)
		//定义一个内置函数用于处理error
		errorHandle := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}

		inArgs := make([]interface{}, 0, len(req))
		for _, arg := range req {
			inArgs = append(inArgs, arg.Interface())
		}

		//组装数据并序列化
		reqRPC := transfer.RPCdata{Name: rpcname, Args: inArgs}
		b, err := transfer.Encode(reqRPC)
		if err != nil {
			panic(err)
		}

		//发送数据
		err = cReqTransPort.Send(b)
		if err != nil {
			return errorHandle(err)
		}

		//读取数据
		rsp, err := cReqTransPort.Read()
		if err != nil {
			return errorHandle(err)
		}

		//反序列化
		respDecode, _ := transfer.Decode(rsp)
		if respDecode.Err != "" {
			return errorHandle(errors.New(respDecode.Err))
		}

		if len(respDecode.Args) == 0 {
			respDecode.Args = make([]interface{}, container.Type().NumOut())
		}

		numOut := container.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 { // unpack arguments (except error)
				if respDecode.Args[i] == nil { // if argument is nil (gob will ignore "Zero" in transmission), set "Zero" value
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(respDecode.Args[i])
				}
			} else { // unpack error argument
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}
		return outArgs
	}

	container.Set(reflect.MakeFunc(container.Type(), f)) //MakeFunc？？？

}
