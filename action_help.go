package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func init() {
	registerAction("help", new(actionHelp))
}

type actionHelp struct{}

func (a actionHelp) Help() (short, log string) {
	return "h", "操作指南"
}

func (a actionHelp) Exec(c *Client, args []string) error {

	//colorReset := string([]byte{byte(27), '[', '0', 'm'})

	paddingLen, pad2 := 0, 0
	for cmd, v := range ActionMap {
		al, _ := v.Help()
		if paddingLen < len(cmd) {
			paddingLen = len(cmd)
		}
		if pad2 < len(al) {
			pad2 = len(al)
		}
	}

	msg := "操作指南\r\n" //\r\n" + colorReset
	for cmd, v := range ActionMap {
		alias, commandDescribe := v.Help()

		spacedCommand := []rune(strings.Repeat(" ", paddingLen))
		for idx, ss := range cmd {
			spacedCommand[idx] = ss
		}
		spacedAlias := []rune(strings.Repeat(" ", pad2))
		for idx, ss := range alias {
			spacedAlias[idx] = ss
		}
		redCommand := color.New(color.FgRed, color.Bold).Sprintf("/%s", string(spacedCommand))
		cyanAlias := color.New(color.Italic, color.FgCyan).Sprintf("/%s", string(spacedAlias))
		commandDescribe = color.New(color.Faint).Sprintf("%s", commandDescribe)
		msg += fmt.Sprintf("%s %s %s\r\n", redCommand, cyanAlias, commandDescribe)
	}
	c.writeBack(msg)
	return nil
}

func (a actionHelp) Hint(args []string) string {
	return ""
}
