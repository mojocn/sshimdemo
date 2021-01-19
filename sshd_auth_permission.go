package main

import (
	"encoding/json"
	"golang.org/x/crypto/ssh"
	"log"
)

const permissionUserIns = "authedUserJson"

func setPermission(user *User, fp, form string) *ssh.Permissions {
	marshal, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}
	return &ssh.Permissions{Extensions: map[string]string{permissionUserIns: string(marshal), "from": form, "fingerprint": fp}}
}
func getPermissionUser(conn *ssh.ServerConn) (*User, error) {
	userString := conn.Permissions.Extensions[permissionUserIns]
	user := new(User)
	err := json.Unmarshal([]byte(userString), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
