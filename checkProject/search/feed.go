/*
 * @Author: haha_giraffe
 * @Date: 2019-12-27 20:16:46
 * @Description: 用于读取json数据
 */
package search

import (
	"encoding/json"
	"os"
)

const dataFile = "data/data.json"

//Feed结构体用于存储json数据
type Feed struct {
	Name string `json:"site"`
	URI  string `json:"link"`
	Type string `json:"type"`
}

//将数据从json文件中解析出来并返回有一个切片
func RetrieveFeeds() ([]*Feed, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//将文件解码到一个切片中
	var feeds []*Feed
	err = json.NewDecoder(file).Decode(&feeds)
	//这里如果有错误也可以直接返回，因为调用者需要自己检查返回值
	return feeds, err
}
