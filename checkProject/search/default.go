/*
 * @Author: haha_giraffe
 * @Date: 2019-12-27 20:44:49
 * @Description:
 */
package search

type defaultMatcher struct{}

func init() {
	var matcher defaultMatcher
	Register("default", matcher)
}

//defaultMatcher 实现的Search方法，也就实现了Matcher接口
func (m defaultMatcher) Search(feed *Feed, searchTerm string) ([]*Result, error) {
	return nil, nil
}
