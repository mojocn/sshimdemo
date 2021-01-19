package main

import (
	"golang.org/x/crypto/ssh"
)

func (c *Client) HandleSftp(sshCh ssh.Channel) {
	//defer sshCh.Close() // SSH_MSG_CHANNEL_CLOSE
	//sftpServer, err := sftp.NewServer(sshCh, sftp.ReadOnly())
	//if err != nil {
	//	return
	//}
	//_ = sftpServer.Serve()
}
