/*
 * @Author: haha_giraffe
 * @Date: 2019-12-27 17:15:06
 * @Description: 实现一个在RSS数据源中进行查找的功能(大小写区分)
 */
package main

import (
	_ "checkProject/matchers" //下划线就可以在不显示调用包中函数情况下，调用包内部init函数和包域变量
	"checkProject/search"
	"fmt"
	"io"
	"log"
	"os"
)

//init函数是优先于main函数执行，且不能被其他函数调用
func init() {
	file, err := os.Create("result")
	if err != nil {
		log.Fatalln("create file error")
	}
	output := io.MultiWriter(os.Stdout, file)
	log.SetOutput(output)
}

func main() {
	var keyword string
	fmt.Println("请输入关键字")
	fmt.Scanln(&keyword)
	search.Run(keyword)
}
