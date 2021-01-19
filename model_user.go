package main

import (
	"fmt"
	"gorm.io/gorm"
)

const (
	UserStatusOnline   = "online"
	UserStatusOffline  = "offline"
	UserStatusDisabled = "disabled"
)

type User struct {
	gorm.Model
	Name   string
	Email  string
	Status string
}

const partitionCnt = 4

func (receiver User) MessengerTopic() (string, int) {
	return fmt.Sprintf("user:%d", receiver.ID), partitionCnt
}

const (
	UserUserStatusPending = "pending"
	UserUserStatusAccept  = "accepted"
	UserUserStatusBroken  = "broken"
)

type UserUser struct {
	gorm.Model
	FooID  uint
	DoeID  uint
	Role   string
	Status string
}
