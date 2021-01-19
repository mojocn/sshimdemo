package main

import (
	"strings"
)

func init() {
	registerAction("square", new(ActionSquare))
}

type ActionSquare struct{}

func (a ActionSquare) Help() (alias, log string) {
	return "square", "议事广场,可以尽情的灌水"
}
func (a ActionSquare) Exec(c *Client, args []string) error {
	c.SetPrompt("⛲️")
	content := strings.Join(args, "")
	return c.postman.SendMsgBroadCast(c.User, content)
}

func (a ActionSquare) Hint(args []string) string {
	return ""
}
