/*
 * @Author: haha_giraffe
 * @Date: 2019-12-27 20:42:54
 * @Description: 匹配寻找
 */
package search

import (
	"fmt"
	"log"
)

type Result struct {
	Field   string
	Content string
}

//定义一个接口
type Matcher interface {
	Search(feed *Feed, searchTerm string) ([]*Result, error)
}

func Match(matcher Matcher, feed *Feed, searchTerm string, results chan<- *Result) {
	searchResults, err := matcher.Search(feed, searchTerm)
	if err != nil {
		log.Println(err)
		return
	}
	for _, result := range searchResults {
		results <- result
	}
}

//输出显示函数，这个是在主goroutine中阻塞
func Display(results chan *Result) {
	//等待channel传数据过来，如果没有则阻塞，如果通道关闭则返回
	for result := range results {
		fmt.Printf("%s:\n%s\n\n", result.Field, result.Content)
	}
}
