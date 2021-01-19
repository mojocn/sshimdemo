package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	registerAction("translate", new(ActionTranslate))
}

type ActionTranslate struct{}

func (a ActionTranslate) Help() (alias, log string) {
	return "translate", "英文翻译成中文:eg /translate this app is awesome and I love it."
}
func (a ActionTranslate) Exec(c *Client, args []string) error {
	line := strings.Join(args[1:], " ")
	ch, err := translateEn2Ch(line)
	if err != nil {
		return err
	}
	c.Primary(ch)
	return nil
}

func (a ActionTranslate) Hint(args []string) string {
	if len(args) < 2 {
		return "必须包含需要翻译的内容"
	}
	return ""
}

func translateEn2Ch(text string) (string, error) {
	url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=ene&tl=zh-cn&dt=t&q=%s", url.QueryEscape(text))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("google translate api status code not 200")
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ss := string(bs)
	ss = strings.ReplaceAll(ss, "[", "")
	ss = strings.ReplaceAll(ss, "]", "")
	ss = strings.ReplaceAll(ss, "null,", "")
	ss = strings.Trim(ss, `"`)
	ps := strings.Split(ss, `","`)
	return ps[0], nil
}
