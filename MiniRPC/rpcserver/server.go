/*
 * @Author: haha_giraffe
 * @Date: 2019-12-29 21:18:32
 * @Description: 服务器
 */
package rpcserver

import (
	"fmt"
	"gotools/MiniRPC/transfer"
	"io"
	"log"
	"net"
	"reflect"
)

type Server struct {
	addr  string
	funcs map[string]reflect.Value
}

func CreateServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

//将函数名注册到server的funcs字段中
func (s *Server) Register(str string, fFunc interface{}) {
	if _, ok := s.funcs[str]; ok {
		return
	}
	// fmt.Println(reflect.ValueOf(fFunc))
	s.funcs[str] = reflect.ValueOf(fFunc)
}

//查询数据并返回
func (s *Server) Execute(data transfer.RPCdata) transfer.RPCdata {
	f, ok := s.funcs[data.Name]
	if !ok {
		//没找到（之前没有注册）
		e := fmt.Sprintf("func %s not register", data.Name)
		log.Println(e)
		return transfer.RPCdata{Name: data.Name, Args: nil, Err: e}
	}
	log.Printf("func %s is called", data.Name)
	//把参数取出来
	inArgs := make([]reflect.Value, len(data.Args))
	for i := range data.Args {
		inArgs[i] = reflect.ValueOf(data.Args[i])
	}
	//执行注册的函数（或者调用 Value.Call([]reflect.Value) ([]reflect.Value)
	out := f.Call(inArgs)
	//这里-1是因为最后一个传回的参数是error，不放到resArgs中
	resArgs := make([]interface{}, len(out)-1)
	for i := 0; i < len(out)-1; i++ {
		resArgs[i] = out[i].Interface()
	}

	var er string
	//如果最后一个返回参数是error类型的，则返回对应的值（放到结构体中Err字段）
	if _, ok := out[len(out)-1].Interface().(error); ok {
		er = out[len(out)-1].Interface().(error).Error()
	}
	return transfer.RPCdata{Name: data.Name, Args: resArgs, Err: er}
}

func (s *Server) Start() {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Println("listen error")
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept err")
			return
		}
		//每次接收到一个连接，就新建一个goroutine来处理逻辑
		go func() {
			connTransport := transfer.NewTransport(conn)
			for {
				//读取数据
				req, err := connTransport.Read()
				if err != nil {
					//如果err == io.EOF则说明传输结束
					if err != io.EOF {
						log.Println("Read error ", err)
						return
					}
				}

				//反序列化
				decReq, err := transfer.Decode(req)
				if err != nil {
					log.Println("decode error")
					return
				}

				response := s.Execute(decReq)

				//序列化
				b, err := transfer.Encode(response)
				if err != nil {
					log.Println("encode error")
					return
				}

				//发送数据
				err = connTransport.Send(b)
				if err != nil {
					log.Println("Send error")
					return
				}
			}
		}()
	}
}
