package security

import (
	"devread/handle_log"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// ref > https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72

func HashAndSalt(pwd []byte) string {
	log, _ := handle_log.WriteLog()
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Error("Error hash and salt password ", zap.Error(err))
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}
