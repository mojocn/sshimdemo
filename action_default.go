package main

import (
	"log"
	"strings"
)

var _ ActionDoer = new(ActionDefault)

type ActionDefault struct{}

func (a ActionDefault) Help() (alias, log string) {
	return "sys", "sys"
}
func (a ActionDefault) Exec(c *Client, args []string) error {
	content := strings.Join(args, "")
	if c.selectedFriend != nil {
		err := c.postman.SendMsgUser(c.User, c.selectedFriend, content)
		return err

	} else if group := c.selectedGroup; group != nil {
		for _, u := range group.Members {
			err := c.postman.SendMsgUser(c.User, &u, content)
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	} else {
		return c.postman.SendMsgBroadCast(c.User, content)
	}
}

func (a ActionDefault) Hint(args []string) string {
	return ""
}
