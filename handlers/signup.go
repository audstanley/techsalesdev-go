package handlers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/sha3"
)

var BadRequest = `{ "status": "Bad Request" }`

// need to check regex better.
func CheckFullUserRegex(c *fiber.Ctx, u *FullUserSigningUp) error {

	if u.Address1 == "" || u.City == "" || u.Email == "" || u.FirstName == "" || u.LastName == "" || u.State == "" {
		fmt.Println("Signup is missing a something")
		c.Status(400)
		return &fiber.Error{Code: 400, Message: BadRequest}
	}
	if !(u.Password1 == u.Password2 && CheckPasswordRegex(c, &u.Password1)) {
		c.Status(400)
		return &fiber.Error{Code: 400, Message: BadRequest}
	}
	return nil
}

func SignUpARealUser(c *fiber.Ctx) error {
	fullUser := FullUserSigningUp{}
	err := c.BodyParser(&fullUser)
	fullUser.Pending = true

	if err != nil {
		return &fiber.Error{Code: 400, Message: `{ "status": "Could not parse body" }`}
	}
	err = CheckFullUserRegex(c, &fullUser)
	if err != nil {
		return err
	}
	// make sure the user isn't already pending email
	val, err := EmailPending.Exists(RedisCtx, fullUser.Email).Result()
	CheckRedisErr(c, err)
	if val == 1 {
		return &fiber.Error{Code: 400, Message: `{ "status": "Email is already pending" }`}
	}

	// make sure the user isn't already a user
	val, err = UserClient.Exists(RedisCtx, fullUser.Email).Result()
	CheckRedisErr(c, err)
	if val == 1 {
		// If the user already exists, but is updating their profile, we will want to
		// check the jwt to let them update their account. Eventually (maybe a TODO).
		return &fiber.Error{Code: 400, Message: `{ "status": "This account already exists" }`}
	}

	// Put the user in the address database.
	h := sha3.New512()
	h.Write([]byte(fullUser.Password1))
	hexStr := hex.EncodeToString(h.Sum(nil))
	fullUser.Password1 = hexStr
	fullUser.Password2 = hexStr
	fullUserMarshalled, err := json.Marshal(fullUser)
	err = UserAddressesClient.Set(RedisCtx, fullUser.Email, string(fullUserMarshalled), 0).Err()

	link := RandStringBytes(32)

	emailUser := EmailPendingUser{
		Email:    fullUser.Email,
		Password: hexStr,
		Iat:      uint64(time.Now().Unix()),
		Link:     link,
	}
	// User Marshalled JSON
	um, _ := json.Marshal(emailUser)
	// Send Data over to Redis
	err = EmailPending.Set(RedisCtx, fullUser.Email, string(um), 0).Err()
	CheckRedisErr(c, err)

	// Create an Etherium wallet for the new user
	// Generate a key for the user to send funds to:
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
	wallet.Email = fullUser.Email
	jmWallet, _ := json.Marshal(&wallet)
	err = UserWalletsClient.Set(RedisCtx, fullUser.Email, string(jmWallet), 0).Err()
	CheckRedisErr(c, err)

	// send email to user that they have a new account
	if DisableSendingEmail == false {
		// Send email to newly created user
		to := []string{fullUser.Email}

		// smtp server configuration.
		smtpHost := "smtp.gmail.com"
		smtpPort := "587"
		auth := smtp.PlainAuth("", Envs["SMTP_ACCOUNT"], Envs["SMTP_PASS"], smtpHost)

		msg := []byte("To: " + fullUser.Email + "\r\n" +
			"Subject: Verification Email\r\n" +
			"\r\n" +
			"You recently signed up for TechSales.dev.\r\n" +
			"Click here to verify your email address: \r\n" +
			"    https://www.techsales.dev/verify/" + link + "\r\n")

		// Sending email.
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, Envs["SMTP_ACCOUNT"], to, msg)
		if err != nil {
			fmt.Println("Error sending email", err)
		}
		fmt.Println("Email Sent Successfully!")
		c.Status(200)
		return c.JSON(map[string]string{"status": "OK"})
	}
	return c.JSON(map[string]string{"status": "OK"})
}
