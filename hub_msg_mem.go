package main

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
)

func NewPostMam(db *gorm.DB) *PostMam {
	return &PostMam{
		Mutex:   sync.Mutex{},
		UserMap: make(map[string]*User),
		hub:     make(map[string]chan Msg),
		db:      db,
	}
}

func (pm *PostMam) RegisterClientDevice(user *User, deviceSessionID string) error {
	pm.Lock()
	defer pm.Unlock()
	pm.UserMap[deviceSessionID] = user
	pm.hub[deviceSessionID] = make(chan Msg)
	return nil
}

func (pm *PostMam) SendMsgUser(from, to *User, content string) error {
	pm.Lock()
	defer pm.Unlock()
	msg := Msg{
		FromID:  from.ID,
		ToID:    to.ID,
		Content: content,
		Status:  "",
	}
	//查找用户的session_id
	var targetDeviceIds []string
	for deviceSessionID, user := range pm.UserMap {
		if user.ID == to.ID {
			targetDeviceIds = append(targetDeviceIds, deviceSessionID)
		}
	}
	if len(targetDeviceIds) == 0 {
		return errors.New("用户是离线状态不能发生消息")
	}

	//向用户session发送消息
	for _, deviceId := range targetDeviceIds {
		channel, ok := pm.hub[deviceId]
		if ok {
			channel <- msg
		}
	}
	return nil
}

func (pm *PostMam) SendMsgBroadCast(me *User, msg string) error {
	//todo check user has right
	pm.Lock()
	defer pm.Unlock()
	for deviceSessionID, to := range pm.UserMap {
		msg := Msg{
			GroupID: 1,
			FromID:  me.ID,
			ToID:    to.ID,
			Content: msg,
			Status:  "",
		}
		pm.hub[deviceSessionID] <- msg
	}
	return nil
}

func (pm *PostMam) ReceiveMsgLoop(deviceSessionID string, writer MsgWriter, exit chan bool) {
	pm.Lock()
	msgChan, ok := pm.hub[deviceSessionID]
	pm.Unlock()
	if !ok {
		log.Println("user msg chan invalid")
		return
	}
	for {
		select {
		case m := <-msgChan:
			writer.WritePigeonMsg(m)
		case <-exit:
			return
		}
	}
}
