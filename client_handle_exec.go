package main

import (
	"golang.org/x/crypto/ssh"
	"log"
)

type execMsg struct {
	Command string
}

type exitStatusMsg struct {
	Status uint32
}

//https://stackoverflow.com/questions/33846959/golang-ssh-server-how-to-handle-file-transfer-with-scp
func (c *Client) HandleExec(msg string, ch ssh.Channel) {

	// ch can be used as a ReadWriteCloser if there should be interactivity
	runCommand(msg, ch)
	ex := exitStatusMsg{
		Status: 0,
	}

	// return the status code
	if _, err := ch.SendRequest("exit-status", false, ssh.Marshal(&ex)); err != nil {
		log.Printf("unable to send status: %v", err)
	}
	ch.Close()
}

func runCommand(cmd string, ch ssh.Channel) {
	//os.Cmd os.Exec 来实现 或者 write 其他的内容
	ch.Write([]byte("你输入的命令:" + cmd + " 功能自己实现TODO"))
}

func (c Client) ParseCommandLine(req *ssh.Request) (string, error) {
	var msg execMsg
	if err := ssh.Unmarshal(req.Payload, &msg); err != nil {
		return "", err
	} else {
		return msg.Command, err
	}
}
