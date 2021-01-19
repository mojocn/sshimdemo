package main

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name    string
	Desc    string
	Status  string
	Members []User `gorm:"-" json:"members"`
}

type GroupUser struct {
	gorm.Model
	GroupID uint
	UserID  uint
	Role    string
	Status  string
}
