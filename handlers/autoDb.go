package handlers

import (
	"encoding/hex"
	"encoding/json"

	"golang.org/x/crypto/sha3"
)

func AutoVerifyAUserInTheDatabase() {
	email := "verifiedUser@gmail.com"
	p := "12345678"
	h := sha3.New512()
	h.Write([]byte(p))
	hexStr := hex.EncodeToString(h.Sum(nil))
	u := User{
		Email:             email,
		Password:          hexStr,
		EmailVerification: true,
	}
	um, _ := json.Marshal(u)
	UserClient.Set(RedisCtx, email, string(um), 0)
}
