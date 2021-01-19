package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"log"
)

func handleChannels(channels <-chan ssh.NewChannel, sshConn *ssh.ServerConn) {
	user, err := getPermissionUser(sshConn) // 获取  SSH User Authentication Protocol 传递过来的认证用户信息
	if err != nil {
		log.Println(err)
		return
	}
	//创建session pty client
	c := NewClient(user, sshConn, postManVar)
	promptString := fmt.Sprintf("[%s] ", user.Name)

	hasShell := false

	for ch := range channels {
		//交互Session
		//一个Session就是一个远程的程序的执行。这个程序或许是shell、应用程序、系统调用或者内建的子系统。它可能没有绑定到虚拟终端上，又或者有或没有涉及到X11转发。同时间，可以有多个Session正在被运行。
		if t := ch.ChannelType(); t != "session" {
			ch.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
			continue
		}

		channel, requests, err := ch.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}
		defer channel.Close()

		c.Term = term.NewTerminal(channel, promptString)

		for req := range requests {
			var width, height int
			var ok bool
			switch req.Type {
			case "shell":
				// 开启交互 tty shell
				//一旦一个Session被设置完毕，在远端就会有一个程序被启动。这个程序可以是一个Shell，也可以时一个应用程序或者是一个有着独立域名的子系统。
				if c.Term != nil && !hasShell {
					go c.HandleShell(channel)
					ok = true
					hasShell = true
				}
			case "pty-req": //通过如下消息可以让服务器为Session分配一个虚拟终端
				//当客户端的终端窗口大小被改变时，或许需要发送这个消息给服务器。
				width, height, ok = parsePtyRequest(req.Payload)
				if ok {
					err := c.Resize(width, height)
					ok = err == nil
				}
			case "window-change":
				width, height, ok = parseWinchRequest(req.Payload)
				if ok {
					err := c.Resize(width, height)
					ok = err == nil
				}
			case "exec":
				// ssh root@mojotv.cn whoami
				//一旦一个Session被设置完毕，在远端就会有一个程序被启动。这个程序可以是一个Shell，也可以时一个应用程序或者是一个有着独立域名的子系统。
				command, err := c.ParseCommandLine(req) // 协议 req.Payload 里面的用户命令输出
				if err != nil {
					log.Printf("error parsing ssh execMsg: %s\n", err)
					return
				} else {
					ok = true
				}
				//开始执行从 whoami 远程shell 命令
				// 执行完成 结果直接返回
				go c.HandleExec(command, channel)
			case "env":
				//在shell或command被开始时之后，或许有环境变量需要被传递过去。然而在特权程序里不受控制的设置环境变量是一个很有风险的事情，
				//所以规范推荐实现维护一个允许被设置的环境变量列表或者只有当sshd丢弃权限后设置环境变量。
				//todo set language i18n
				log.Print(string(req.Payload))
			case "subsystem":
				//一旦一个Session被设置完毕，在远端就会有一个程序被启动。这个程序可以是一个Shell，也可以时一个应用程序或者是一个有着独立域名的子系统。
				// 实现一下功能可以实现 sftp功能
				//fmt.Fprintf(debugStream, "Subsystem: %s\n", req.Payload[4:])
				if string(req.Payload[4:]) == "sftp" {
					ok = true
					go c.HandleSftp(channel)
				}

			default:
				log.Println(req.Type, string(req.Payload))
			}
			if req.WantReply {
				req.Reply(ok, nil)
			}
		}
	}
}
