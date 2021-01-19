package main

import (
	"errors"
	"fmt"
	"strconv"
)

type actionFriendAdd struct{}

func init() {
	registerAction("friend.add", new(actionFriendAdd))
}

func (a actionFriendAdd) Help() (short, log string) {
	return "friend.add", "搜索用户添加好友"
}
func (a actionFriendAdd) Exec(c *Client, args []string) error {
	users := c.postman.SearchUsers("")
	for idx, u := range users {
		line := fmt.Sprintf("%d   %s\r\n", idx, u.Name)
		c.writeBack(line)
	}

	c.Term.SetPrompt("请选择需要添加好友的用户ID:")
	input, err := c.Term.ReadLine()
	if err != nil {
		return err
	}
	user, err := selectUser(users, input)
	if err != nil {
		return err
	}
	c.setSessionPrompt()
	return c.postman.FriendMake(c.User, user)
}

func (a actionFriendAdd) Hint(args []string) string {
	return ""
}

func selectUser(users []User, input string) (*User, error) {
	idx, err := strconv.Atoi(input)
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= len(users) {
		return nil, errors.New("无效的序号")
	}
	v := users[idx]
	return &v, nil
}
