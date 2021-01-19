package main

import (
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
)

func (c *Client) HandleShell(channel ssh.Channel) {
	defer channel.Close()
	exitChan := make(chan bool, 1)
	// FIXME: This shouldn't live here, need to restructure the call chaining.
	//c.Server.Add(c)
	// 用户进入聊天界面 开始注册用户在线状态, 聊天消息的队列
	err := c.postman.RegisterClientDevice(c.User, c.DeviceSessionID)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		// Block until done, then remove.
		c.Conn.Wait()
		c.closed = true
		//c.Server.Remove(c)
		//close(c.Messages)
		//c.postman.UserOffline(c.User)  // 用户离线
	}()

	go func() {
		//todo:: send history msg
		//for msg := range c.Messages {
		//	c.Write(msg)
		//}
		// 接受其他用户发送给你的消息 或 广播消息
		c.postman.ReceiveMsgLoop(c.DeviceSessionID, c, exitChan)
	}()
	new(actionHelp).Exec(c, nil) //输出帮助信息
	for {
		line, err := c.Term.ReadLine()
		if err != nil {
			break
		}
		// 使用 默认的 hook action 来处理 交互shell的键盘输入(聊天 或者 指令)
		var doer ActionDoer = new(ActionDefault)
		// choose action
		isCmd, action, args := parseInputLine(line) // 解析用户输入 return 是否是指令 or 是发送聊天消息
		if isCmd {
			v, ok := ActionMap[action]
			if ok {
				doer = v // 匹配指令的 hook
			} else {
				c.Danger("未知动作指令: " + line)
				continue
			}
		}
		//
		if hint := doer.Hint(args); hint != "" { //check arg
			c.Warning("Invalid command: " + line)
			continue
		}
		err = doer.Exec(c, args) // 执行自定义的action hook
		if err != nil {
			c.Warning(err.Error())
			log.Println(err)
			//c.TermWrite(err.Error())
		}
	}

}

//parseInputLine 解析input
func parseInputLine(line string) (isCmd bool, action string, args []string) {
	parts := strings.Split(line, " ")
	if len(parts) > 0 && strings.HasPrefix(parts[0], "/") {
		args = []string{}
		for _, p := range parts {
			if t := strings.TrimSpace(p); t != "" {
				args = append(args, t)
			}
		}
		return true, strings.TrimPrefix(args[0], "/"), args
	}
	return false, "", []string{line}
}
