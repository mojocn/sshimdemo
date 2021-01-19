package main

type Mailer interface {
	officer
	messenger
}

type MsgWriter interface {
	Warning(msg string)
	Danger(msg string)
	Primary(msg string)
	Success(msg string)
	WritePigeonMsg(msg Msg)
}

type officer interface {
	MyFriends(me *User, q string) []User
	SearchUsers(q string) []User
	UserOnline(me *User) error
	UserOffline(me *User) error
	FriendMake(me, other *User) error
	FriendAccept(me *User, uu *UserUser) error
	GroupJoin(me *User, group *Group) error
	GroupApprove(me *User, gu *GroupUser) error
	GetFriends(me *User) ([]User, error)
	GetMembers(group *Group) ([]User, error)
	RegisterClientDevice(user *User, deviceID string) error
}

type messenger interface {
	SendMsgUser(from, to *User, content string) error
	SendMsgBroadCast(me *User, content string) error
	ReceiveMsgLoop(deviceSessionID string, mw MsgWriter, exit chan bool)
}
