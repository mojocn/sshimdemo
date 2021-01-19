package main

import "gorm.io/gorm"

type Msg struct {
	gorm.Model
	DeviceID string
	GroupID  uint
	FromID   uint
	ToID     uint
	Content  string
	Status   string
}
