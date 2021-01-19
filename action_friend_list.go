package main

import (
	"fmt"
)

const (
	inputExit = "/q"
	inputHome = "/home"
)

type actionFriendList struct{}

func init() {
	registerAction("friend", new(actionFriendList))
}

func (a actionFriendList) Help() (short, log string) {
	return "friends", "显示当前用户的好友列表"
}
func (a actionFriendList) Exec(c *Client, args []string) error {

	users := c.postman.MyFriends(c.User, "")
	for idx, u := range users {
		line := fmt.Sprintf("%d  %s(%s)\r\n", idx, u.Name, u.Status)
		c.writeBack(line)
	}
	c.Term.SetPrompt("请选择需要对话好友ID,输入/q 推出选择:")
	input, err := c.Term.ReadLine()
	if err != nil {
		return err
	}
	if input == inputHome || input == inputExit {
		c.PromptHome()
		return nil
	}

	user, err := selectUser(users, input)
	if err != nil {
		c.PromptHome()
		return err
	}
	c.selectedFriend = user
	c.setSessionPrompt()
	return nil
}

func (a actionFriendList) Hint(args []string) string {
	return ""
}
