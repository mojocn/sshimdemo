package main

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"log"
	"time"
)

type Client struct {
	DeviceSessionID string //设备uuid
	Conn            *ssh.ServerConn
	postman         *PostMam //处理 关系和 mq
	User            *User    //当前用户
	selectedFriend  *User    //对话的好友
	selectedGroup   *Group   //对话的群组
	Color           string
	IsAdmin         bool
	ready           chan struct{}
	Term            *term.Terminal
	termWidth       int
	termHeight      int
	LastTX          time.Time
	beepMe          bool
	colorMe         bool
	closed          bool
}

// NewClient constructs a new client
// 1.记录client terminal的状态
// 2.当前用户的状态
// 3.消息发送接收
// 4.好友群组关系管理
// 5.读取客户段terminal的输入
func NewClient(user *User, conn *ssh.ServerConn, pm *PostMam) *Client {

	return &Client{
		DeviceSessionID: string(conn.SessionID()),
		Conn:            conn,
		postman:         pm,
		User:            user,
		selectedFriend:  nil,
		selectedGroup:   nil,
		Term:            nil, //pty
		termWidth:       0,
		termHeight:      0,
	}
}

// TermWrite 写入消息到当强用户ssh 客户端
func (c *Client) writeBack(msg string) {
	c.Term.Write([]byte(msg))
}

func (c *Client) PromptHome() {
	c.Term.SetPrompt(fmt.Sprintf("[%s]", "🌏"))
}
func (c *Client) SetPrompt(s string) {
	c.Term.SetPrompt(fmt.Sprintf("[%s]", s))
}

func (c *Client) Danger(msg string) {
	content := color.RedString("🔴  %s\r\n", msg)
	c.writeBack(content)
}
func (c *Client) Warning(msg string) {
	content := color.YellowString("🟠  %s\r\n", msg)
	c.writeBack(content)
}

func (c *Client) Success(msg string) {
	content := color.GreenString("🟢  %s\n", msg)
	c.writeBack(content)
}

func (c *Client) Primary(msg string) {
	content := color.BlueString("🔵  %s\r\n", msg)
	c.writeBack(content)
}

func (c *Client) MsgPrivate(msg string) {
	content := color.HiCyanString("💬  %s\r\n", msg)
	c.writeBack(content)
}

func (c *Client) MsgGroup(msg string) {
	content := color.HiYellowString("📻 %s\r\n", msg)
	c.writeBack(content)
}
func (c *Client) WritePigeonMsg(msg Msg) {
	if msg.GroupID > 0 {
		c.MsgGroup(msg.Content + "\r\n")
	} else {
		c.MsgPrivate(msg.Content + "\r\n")
		return
	}
}

// Resize resizes the client to the given width and height
func (c *Client) Resize(width, height int) error {
	width = 1000000 // TODO: Remove this dirty workaround for text overflow once ssh/terminal is fixed
	err := c.Term.SetSize(width, height)
	if err != nil {
		log.Printf("Resize failed: %dx%d", width, height)
		return err
	}
	c.termWidth, c.termHeight = width, height
	return nil
}

func (c *Client) setSessionPrompt() {
	prompt := fmt.Sprintf("[%s]", c.User.Name)
	if c.selectedFriend != nil {
		prompt = fmt.Sprintf("[%s -> %s]", c.User.Name, c.selectedFriend.Name)
	}
	if c.selectedGroup != nil {
		prompt = fmt.Sprintf("[%s IN %s]", c.User.Name, c.selectedGroup.Name)
	}
	c.Term.SetPrompt(prompt)
}
