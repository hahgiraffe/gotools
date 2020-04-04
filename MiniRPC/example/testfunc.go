/*
 * @Author: haha_giraffe
 * @Date: 2020-01-04 16:58:41
 * @Description: 测试函数调用
 */
package example

import (
	"fmt"
	"gotools/MiniRPC/rpcclient"
	"gotools/MiniRPC/rpcserver"
	"time"
)

func AddSum(a, b int) (res int, err error) {
	return a + b, nil
}

func FuncTest() {
	ser := rpcserver.CreateServer("localhost:9999")
	ser.Register("AddSum", AddSum)
	go ser.Start()

	time.Sleep(1 * time.Second)
	cli := rpcclient.CreateClient("localhost:9999")
	var Add func(int, int) (int, error)
	cli.Callrpc("AddSum", &Add)
	u, err := Add(1, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	u2, err := Add(24, 564)
	if err != nil {
		panic(err)
	}
	fmt.Println(u2)
}
