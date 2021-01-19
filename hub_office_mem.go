package main

import (
	"gorm.io/gorm"
	"log"
	"sync"
)

type PostMam struct {
	sync.Mutex
	UserMap map[string]*User    // ssh session id => User
	hub     map[string]chan Msg // ssh session id => User
	db      *gorm.DB
}

func (pm *PostMam) SearchUsers(q string) []User {
	var list []User
	err := pm.db.Find(&list).Error
	if err != nil {
		log.Println(err)
		return nil
	}
	return list
}

func (pm *PostMam) MyFriends(me *User, q string) []User {
	var uids []int
	var list []User
	err := pm.db.Model(new(User)).
		Joins("JOIN user_users AS aa ON aa.doe_id = ? AND aa.status = ? AND aa.foo_id = users.id", me.ID, UserUserStatusAccept).
		Joins("JOIN user_users AS bb ON bb.foo_id = ? AND bb.status = ? AND bb.doe_id = users.id", me.ID, UserUserStatusAccept).
		Distinct("users.id").Pluck("id", &uids).Error
	if err != nil {
		log.Println(err)
		return nil
	}
	err = pm.db.Where("id <> ?", me.ID).Find(&list, uids).Error
	if err != nil {
		log.Println(err)
		return nil
	}
	return list
}

func (pm *PostMam) UserOnline(me *User) error {
	panic("implement me")
}

func (pm *PostMam) UserOffline(me *User) error {
	panic("implement me")
}

func (pm *PostMam) FriendMake(me, other *User) error {

	ins := UserUser{
		Model:  gorm.Model{},
		FooID:  me.ID,
		DoeID:  other.ID,
		Role:   "",
		Status: UserUserStatusAccept,
	}
	return pm.db.Save(&ins).Error

}

func (pm *PostMam) FriendAccept(me *User, uu *UserUser) error {
	panic("implement me")
}

func (pm *PostMam) GroupJoin(me *User, group *Group) error {
	panic("implement me")
}

func (pm *PostMam) GroupApprove(me *User, gu *GroupUser) error {
	panic("implement me")
}

func (pm *PostMam) GetFriends(me *User) ([]User, error) {
	panic("implement me")
}

func (pm *PostMam) GetMembers(group *Group) ([]User, error) {
	panic("implement me")
}
