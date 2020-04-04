/*
 * @Author: haha_giraffe
 * @Date: 2020-01-04 17:04:13
 * @Description: file content
 */
package example

import (
	"encoding/gob"
	"fmt"
	"gotools/MiniRPC/rpcclient"
	"gotools/MiniRPC/rpcserver"
	"log"
	"time"
)

type User struct {
	Name string
	Age  int
}

var userDB = map[int]User{
	1:  User{"chs", 43},
	3:  User{"goodman", 433},
	7:  User{"monkey", 34},
	10: User{"monster", 64},
}

func QueryUser(id int) (User, error) {
	if u, ok := userDB[id]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("id %d not in user db", id)
}

var GetDB = map[int]string{
	1:     "1",
	10:    "10",
	100:   "100",
	1000:  "1000",
	10000: "10000",
}

func GetFunc(id int) (string, error) {
	if str, ok := GetDB[id]; ok {
		return str, nil
	}
	return "", fmt.Errorf("id %d not found ", id)
}

func Teststruct() {
	fmt.Println("Hello RPC")
	gob.Register(User{})
	//start server
	ser := rpcserver.CreateServer("localhost:3212")
	ser.Register("QueryUser", QueryUser)
	ser.Register("Get", GetFunc)
	go ser.Start()

	//time wait
	time.Sleep(1 * time.Second)

	//start client
	cli := rpcclient.CreateClient("localhost:3212")
	var Query func(int) (User, error)
	cli.Callrpc("QueryUser", &Query)
	u, err := Query(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	u2, err := Query(7)
	if err != nil {
		panic(err)
	}
	fmt.Println(u2)
	fmt.Println("-----------")
	var Get func(int) (string, error)
	cli.Callrpc("Get", &Get)
	str, err := Get(10)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("str is ", str)
	str, err = Get(1000)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("str is ", str)
	str, err = Get(888)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("str is ", str)

}
