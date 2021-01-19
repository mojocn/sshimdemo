package main

import "log"

var ActionMap = map[string]ActionDoer{}

//ActionDoer 编写插件hook 来扩展更多的功能
type ActionDoer interface {
	Help() (alias, long string)
	Exec(c *Client, args []string) error
	Hint(args []string) string
}

//registerAction 注册编写的action hook 扩展功能
func registerAction(name string, doer ActionDoer) {
	_, ok := ActionMap[name]
	if ok {
		log.Fatal("action has already existed: ", name)
	} else {
		ActionMap[name] = doer
	}
}
