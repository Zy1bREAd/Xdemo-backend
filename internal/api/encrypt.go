package api

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// 使用bcrypt进行加密、解密
func EncryptWithHash(pwd string) ([]byte, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("加密密码With Hash时发生错误 ", err)
		return nil, err
	}
	return hashPwd, nil
}

// 校验密码的哈希值是否一致
func ComparePasswordWithHash(pwd string, hashPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd))
	return err == nil
}
