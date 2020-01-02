/*
 * @Author: haha_giraffe
 * @Date: 2019-12-27 17:44:22
 * @Description: 主要的逻辑和业务实现
 */
package search

import (
	"log"
	"sync"
)

var matchers = make(map[string]Matcher)

func Run(s string) {
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}
	results := make(chan *Result)
	var waitGroup sync.WaitGroup
	// fmt.Printf("feeds len is %d, name is %s, URI is %s, Type is %s\n", len(feeds), feeds[0].Name, feeds[0].URI, feeds[0].Type)
	waitGroup.Add(len(feeds))
	for _, feed := range feeds {
		//开始解析数据源
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}
		//为每个数据流创建一个goroutine执行搜索功能，当goroutine结束的时候就减少waitGroup
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, s, results)
			waitGroup.Done()
		}(matcher, feed) //TODO 现在传入这个matcher，调用Search返回空
	}
	//再开启一个goroutine进行等待waitGroup，当其为0的时候再关闭channel，使得Display中range循环结束并返回
	go func() {
		waitGroup.Wait() //goroutine在这里阻塞直到waitGroup计数到达0
		close(results)
	}()
	Display(results)
}

//注册一个匹配器给后面的程序使用，在init函数中调用，这里在rss中注册了一个，default中注册了一个
func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Match already registered")
	}
	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}
