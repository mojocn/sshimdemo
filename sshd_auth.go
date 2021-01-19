package main

//import (
//	"bytes"
//	"errors"
//	"fmt"
//	"github.com/sirupsen/logrus"
//	"golang.org/x/crypto/ssh"
//	"gorm.io/gorm"
//	"strings"
//)
//
//func authUserMfa(database *gorm.DB) func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
//	return func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
//		dbUser, err := getDbUser(database, conn)
//		if err != nil {
//			return nil, err
//		}
//		pt := "password-2fa"
//		if dbUser.MFAStatus == UserMfaStatusOff {
//			return nil, errors.New("用户必须开启MFA认证才能登录")
//		}
//		err = dbUser.Check2FA(string(password))
//		if err != nil {
//			logrus.WithError(err)
//			return nil, err
//		}
//		perm := &ssh.Permissions{
//			Extensions: map[string]string{
//				"authType": pt,
//				"user":     conn.User(),
//				"userID":   fmt.Sprintf("%d", dbUser.ID),
//			},
//		}
//		return perm, nil
//	}
//
//}
//func authKeyboard(database *gorm.DB) func(conn ssh.ConnMetadata, client ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
//	return func(conn ssh.ConnMetadata, client ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
//		dbUser, err := getDbUser(database, conn)
//		if err != nil {
//			return nil, err
//		}
//		pt := "password-keyboard"
//		if dbUser.MFAStatus == UserMfaStatusOff {
//			return nil, errors.New("用户必须开启MFA认证才能登录")
//		}
//		ans, err := client("", "", []string{"请输入您的堡垒机密码：", "请输入您2FA的6位数字:"}, []bool{false, true})
//		if err != nil {
//			return nil, fmt.Errorf("keyboard认证错误：%v", err)
//		}
//
//		err = dbUser.Check2FA(ans[1])
//		if err != nil {
//			return nil, fmt.Errorf("2fa错误：%v", err)
//		}
//		err = dbUser.CheckMyPassword(ans[0])
//		if err != nil {
//			return nil, fmt.Errorf("密码错误：%v", err)
//		}
//		perm := &ssh.Permissions{
//			Extensions: map[string]string{
//				//"authType": pt,
//				"user":     conn.User(),
//				"userID":   fmt.Sprintf("%d", dbUser.ID),
//			},
//		}
//		return perm, nil
//	}
//}
//
//func authUserPublicKey(database *gorm.DB) func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
//	return func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
//		dbUser, err := getDbUser(database, conn)
//		if err != nil {
//			return nil, err
//		}
//		err = publicKeyVerify(dbUser, key)
//		if err != nil {
//			return nil, err
//		}
//		perm := &ssh.Permissions{
//			Extensions: map[string]string{
//				"authType":  key.Type(),
//				"pubkey-fp": ssh.FingerprintSHA256(key),
//				"user":      conn.User(),
//				"userID":    fmt.Sprintf("%d", dbUser.ID),
//			},
//		}
//		return perm, nil
//	}
//}
//
//func getDbUser(database *gorm.DB, conn ssh.ConnMetadata) (dbUser *User, err error) {
//	dbUser = new(User)
//	err = database.Where("email = ?", conn.User()+"@mojotv.cn").Take(dbUser).Error
//	if err != nil {
//		return nil, err
//	}
//	return dbUser, nil
//}
//
//func publicKeyVerify(dbUser *User, inKey ssh.PublicKey) (err error) {
//	if dbUser.SshPublicKey == nil {
//		return errors.New("user's public key is empty")
//	}
//	myPublicKey := strings.TrimSpace(*dbUser.SshPublicKey)
//	dbPubKeyBytes := []byte(myPublicKey)
//	for len(dbPubKeyBytes) > 0 {
//		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(dbPubKeyBytes)
//		if err != nil {
//			logrus.WithError(err).Error("database user public key is invalid")
//		}
//		if bytes.Equal(pubKey.Marshal(), inKey.Marshal()) {
//			return nil
//		}
//		dbPubKeyBytes = rest
//	}
//	return errors.New("there is no matched public key in known keys")
//}

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//authPublicKeysOfGithub github.com 公钥身份任职
func authPublicKeysOfGithub(db *gorm.DB) func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	return func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		userName := conn.User()
		err := sshPublicKeysAuthByGithub(userName, key) //用过github.com api 获取用户名的公钥 校验
		if err != nil {
			return nil, err
		}
		one := new(User)
		err = db.Where("name = ?", userName).Take(one).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			one.Name = userName
			err := db.Save(one).Error
			if err != nil {
				log.Println(err)
			}
		}
		return setPermission(one, Fingerprint(key), "github"), nil
	}
}

//authKeyboard 用户普通匿名登录 记录用户信息
func authKeyboard(db *gorm.DB) func(conn ssh.ConnMetadata, challenge ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
	return func(conn ssh.ConnMetadata, challenge ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
		userName := conn.User()
		one := new(User)
		err := db.Where("name = ?", userName).Take(one).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			one.Name = userName
			err := db.Save(one).Error
			if err != nil {
				log.Println(err)
			}
		}
		return setPermission(one, "", "anon"), nil
	}
}

//Fingerprint 计算公钥指纹
func Fingerprint(k ssh.PublicKey) string {
	hash := md5.Sum(k.Marshal())
	r := fmt.Sprintf("% x", hash)
	return strings.Replace(r, " ", ":", -1)
}

//sshPublicKeysAuthByGithub 比较github的公钥
func sshPublicKeysAuthByGithub(user string, key ssh.PublicKey) error {
	publicKeys, err := fetchGithubPublicKeys(user)
	if err != nil {
		return err
	}
	for _, pbk := range publicKeys {
		if bytes.Equal(key.Marshal(), pbk.Marshal()) {
			return nil
		}
	}
	return fmt.Errorf("the key is not match any https://github.com/%s.keys", user)
}

//fetchGithubPublicKeys 获取当用户名的公钥
func fetchGithubPublicKeys(githubUser string) ([]ssh.PublicKey, error) {
	keyURL := fmt.Sprintf("https://github.com/%s.keys", githubUser)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	defer cancelFunc()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, keyURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("invalid response from github")
	}
	authorizedKeysBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body:%v", err)
	}
	var keys []ssh.PublicKey
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			return nil, fmt.Errorf("parsing key: %v", err) //errors.Wrap(err, "parsing key")
		}
		keys = append(keys, pubKey)
		authorizedKeysBytes = rest
	}
	return keys, nil
}
