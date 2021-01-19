package main

import (
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"log"
	"net"
)

var privateBytes = []byte(`
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----
`) // sshd-server 的私钥正式

// startSshSvrListen 启动ssh-server服务 db来管理 用户登录 和 用户登录日志 完成一些聊天的功能
// 这个方法将被 main.go 的 main方法调用,ssh-server的启动入口
func startSshSvrListen(addr string, db *gorm.DB) {
	//初始话 ssh-server 客户端的配置, 用户认证
	config := &ssh.ServerConfig{
		NoClientAuth: false, //如果是true ssh-sever不需要用户认证
		MaxAuthTries: 6,     //用户认证重试次数
		//PasswordCallback:            authUserMfa(db),       //ssh用户名密码认证 可以扩展成  LDAP ... google MFA
		PublicKeyCallback:           authPublicKeysOfGithub(db), //github用户ssh公钥登录 https://github.com/${githubUserName}.keys  https://github.com/mojocn.keys
		KeyboardInteractiveCallback: authKeyboard(db),           //键盘问答输入用户认证 可以扩展成  LDAP ... google MFA
		AuthLogCallback:             nil,                        //记录用户登录认证的callback
		ServerVersion:               "",                         //"ssh可以扩展的更多功能的聊天服务",   //自定义服务端的版本信息
	}
	// 解析需要给服务端设置ssh 私钥
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("无效私钥证书:", err)
	}
	config.AddHostKey(private)

	//监听socket
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("启动socket失败::", err)
	}
	for {
		// 处理连接
		conn, err := listener.Accept()
		if err != nil {
			// handle error
			log.Println(err)
			continue
		}
		// 用户认证协议SSH User Authentication Protocol
		// 开始工作
		// 开始 handshake 用户登录之前这里用户身份认证, ssh.NewServerConn 会调用上面 PasswordCallback  PublicKeyCallback KeyboardInteractiveCallback ...的callback
		// 用户登录成之后需要向后传递的参数可以 从 sConn.Permissions 中获取
		sConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			// handle error
			log.Print(err)
			continue
		}
		//处理 连接协议SSH Connection Protocol
		// 用户handshake 认证成功
		// 强制必须 丢弃服务的request,防止被攻击
		go ssh.DiscardRequests(reqs)
		// 核心/VIP/MVP 处理连接协议SSH Connection Protocol
		go handleChannels(chans, sConn)
	}
}
