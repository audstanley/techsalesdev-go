package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/sha3"
)

var DisableSendingEmail = false

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var RedisCtx = context.Background()

func Session(c *fiber.Ctx) error {
	// Middleware for user sessions
	fmt.Println("session middleware called")

	wwwAuthentication := c.Get("WWW-Authentication")
	authorization := c.Get("Authorization")
	fmt.Println(c.Get("WWW-Authentication"))
	switch wwwAuthentication {
	case "upass":
		// skip to CreateUseUserClientr Handler (for POST)
		if string(c.Request().Header.Method()) == "WAS_POST" {
			// DEPRECATED
			CreateUser(c)
		} else if string(c.Request().Header.Method()) == "POST" {
			// User is logging in.
			return c.Next()
		}
		return c.Next()

	case "gtoken":
		fmt.Println("handler.session code may need some attention for gtoken")
		fmt.Println(wwwAuthentication, authorization)
		return c.Next()

	case "token":
		fmt.Println("handler.session code may need some attention for token")
		fmt.Println(wwwAuthentication, authorization)
		c.Set("WWW-Authentication", "token")
		currentToken := c.Get("Authorization", "")
		if currentToken == "" {
			fmt.Println("middleware: bad request token")
			c.JSON(BadRequest)
		}
		token, err := jwt.Parse(c.Get("Authorization", ""), func(token *jwt.Token) (interface{}, error) {
			return []byte(Envs["ACCESS_TOKEN_SECRET"]), nil
		})
		if err != nil {
			fmt.Println("token err", err)
			return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: fiber.ErrBadRequest.Message}
		}
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println("token checks out:", claims["user"])
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString([]byte(Envs["ACCESS_TOKEN_SECRET"]))
		c.Set("WWW-Authentication", "token")
		c.Set("Authorization", t)
		claimUserStr, _ := claims["user"].(string) // cast interface{} to string [ignore err since we know it's ok]
		redisErr := UserTokensClient.Set(RedisCtx, claimUserStr, t, 0).Err()
		CheckRedisErr(c, redisErr)
		fmt.Println("boop")
		return c.Next()
	case "none":
		return c.Next()
	default:

		return c.Next()

	}

}

func CreateUser(c *fiber.Ctx) error {
	// WILL BE DEPRECATED
	wwwAuthentication := c.Get("WWW-Authentication")
	authorization := c.Get("Authorization")
	switch wwwAuthentication {
	case "upass":
		if authorization == "none" {
			username := c.Get("x-username")
			password := c.Get("x-password")
			username = strings.ToLower(username)
			if username != "" && password != "" {
				link := RandStringBytes(32)
				h := sha3.New512()
				h.Write([]byte(password))
				hexStr := hex.EncodeToString(h.Sum(nil))
				u := EmailPendingUser{
					Email:    username,
					Password: hexStr,
					Iat:      uint64(time.Now().Unix()),
					Link:     link,
				}
				// User Marshalled JSON
				um, _ := json.Marshal(u)
				// Send Data over to Redis
				err := EmailPending.Set(RedisCtx, username, string(um), 0).Err()
				CheckRedisErr(c, err)
				if DisableSendingEmail == false {
					// Send email to newly created user
					to := []string{username}

					// smtp server configuration.
					smtpHost := "smtp.gmail.com"
					smtpPort := "587"
					auth := smtp.PlainAuth("", Envs["SMTP_ACCOUNT"], Envs["SMTP_PASS"], smtpHost)

					msg := []byte("To: " + username + "\r\n" +
						"Subject: Verification Email\r\n" +
						"\r\n" +
						"You recently signed up for TechSales.dev.\r\n" +
						"Click here to verify your email address: \r\n" +
						"    https://api.techsales.dev/verify/" + link + "\r\n")

					// Sending email.
					err := smtp.SendMail(smtpHost+":"+smtpPort, auth, Envs["SMTP_ACCOUNT"], to, msg)
					if err != nil {
						fmt.Println("Error sending email", err)
					}
					fmt.Println("Email Sent Successfully!")
					c.Status(200)
					return c.JSON(map[string]string{"status": "OK"})
				}
			}
			// Else, user needs to supply username/password
		}
	}
	c.Status(400)
	noUser, _ := json.Marshal(map[string]string{"status": "Did not Create User"})
	return &fiber.Error{Code: 400, Message: string(noUser)}
}

func VerifyUserLogin(c *fiber.Ctx) error {
	wwwAuthentication := c.Get("WWW-Authentication")
	authorization := c.Get("Authorization")

	switch wwwAuthentication {
	case "upass":
		if authorization == "none" {
			username := c.Get("x-username")
			password := c.Get("x-password")
			username = strings.ToLower(username)
			if username != "" && password != "" {

				var cursor uint64
				keys, cursor, err := UserClient.Scan(RedisCtx, cursor, "*", 1000000).Result()
				CheckRedisErr(c, err)
				fmt.Println(keys)
				for _, key := range keys {
					// we need to query redis for the object from the key.
					v, e := UserClient.Get(RedisCtx, key).Result()
					CheckRedisErr(c, e)
					var u User
					json.Unmarshal([]byte(v), &u)
					if u.Email == username {
						h := sha3.New512()
						h.Write([]byte(password))
						hexStr := hex.EncodeToString(h.Sum(nil))
						if u.Password == hexStr {
							token := jwt.New(jwt.SigningMethodHS256)
							claims := token.Claims.(jwt.MapClaims)

							claims["user"] = u.Email
							claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
							fmt.Println(wwwAuthentication, authorization)
							t, err := token.SignedString([]byte(Envs["ACCESS_TOKEN_SECRET"]))
							if err != nil {
								return c.SendStatus(fiber.StatusInternalServerError)
							}
							c.Set("WWW-Authentication", "token")
							c.Set("Authorization", t)
							redisErr := UserTokensClient.Set(RedisCtx, u.Email, t, 0).Err()
							CheckRedisErr(c, redisErr)
							return c.JSON(map[string]string{"status": "OK"})
						}
					}
				}
			}
			c.Status(400)
			// Else, user needs to supply username/password
			return c.JSON(map[string]string{"status": "Invalid Username or Password"})
		}
	case "token":
		// The json web token gets verified in the "Session" function middleware
		// Otherwise, the API wont get to this response.
		return c.JSON(map[string]string{"status": "OK"})
	}
	return c.SendString("GET PLACEHOLDER - method DNE.")
}

// /verify/:link -> This Function Verifies the link created for a new user.
func Verify(c *fiber.Ctx) error {
	var cursor uint64
	keys, cursor, err := EmailPending.Scan(RedisCtx, cursor, "*", 1000000).Result()
	CheckRedisErr(c, err)

	for _, key := range keys {
		// we need to query redis for the object from the key.
		v, e := EmailPending.Get(RedisCtx, key).Result()
		CheckRedisErr(c, e)
		var u EmailPendingUser
		json.Unmarshal([]byte(v), &u)
		if u.Link == c.Params("link", "") {
			val, err := UserAddressesClient.Exists(RedisCtx, u.Email).Result()
			fmt.Println("userAddressClient.Exists in verify Redis result is value of:", val)
			CheckRedisErr(c, err)
			if val == 1 {
				var uc User
				uc.Email = u.Email
				uc.Password = u.Password
				uc.EmailVerification = true
				ucjm, _ := json.Marshal(uc)
				UserClient.Set(RedisCtx, u.Email, string(ucjm), 0).Err()
			} else {
				c.Status(400)
				return c.JSON(map[string]string{"status": "user already registered."})
			}

			// Update the adress saved user to note that they are no longer pending email verification
			valOfAddress, _ := UserAddressesClient.Get(RedisCtx, u.Email).Result()
			updateAddr := FullUserSigningUp{}
			json.Unmarshal([]byte(valOfAddress), &updateAddr)
			updateAddr.Pending = false
			jm, _ := json.Marshal(updateAddr)

			UserAddressesClient.Set(RedisCtx, u.Email, string(jm), 0)

			// Update the user's Wallet (no longer pending)
			wallet := EtheriumWallet{}
			walletStr, _ := UserWalletsClient.Get(RedisCtx, u.Email).Result()
			json.Unmarshal([]byte(walletStr), &wallet)
			wallet.Pending = false
			jmWallet, _ := json.Marshal(wallet)
			// update the wallet, so they are no longer in pending status
			UserAddressesClient.Set(RedisCtx, u.Email, string(jmWallet), 0).Err()

			// }
			EmailPending.Del(RedisCtx, key)
			c.Status(200)
			return c.JSON(map[string]string{"status": "user validated. please close this window and log in."})
		}
	}
	c.Status(403)
	return c.JSON(map[string]string{"status": "Forbidden"})
}
