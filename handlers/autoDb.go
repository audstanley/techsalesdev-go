package handlers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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

	// Create a wallet for this test user:
	privateKey, err := crypto.GenerateKey()
	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	wallet := EtheriumWallet{Pending: true}
	wallet.Private = hexutil.Encode(privateKeyBytes)
	wallet.Public = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	wallet.Email = email
	jmWallet, _ := json.Marshal(&wallet)
	err = UserWalletsClient.Set(RedisCtx, email, string(jmWallet), 0).Err()
	CheckRedisErrWithoutContext(err)
}
