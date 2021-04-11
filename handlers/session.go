package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/sha3"
)

var DisableSendingEmail = true

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
	fmt.Println("session middleware called")
	fmt.Println(c.Cookies("username"), c.Cookies("password"))

	wwwAuthentication := c.Get("WWW-Authentication")
	authorization := c.Get("Authorization")
	switch wwwAuthentication {
	case "upass":
		// skip to CreateUser Handler (GET and POST methods)
		if string(c.Request().Header.Method()) == "POST" {
			CreateUser(c)
		} else if string(c.Request().Header.Method()) == "GET" {
			return c.Next()
		}
		return c.Next()

	case "gtoken":
		fmt.Println("handler.session code may need some attention for gtoken")
		fmt.Println(wwwAuthentication, authorization)
		c.Next()

	case "token":
		fmt.Println("handler.session code may need some attention for token")
		fmt.Println(wwwAuthentication, authorization)
		c.Next()
	default:
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		userString := RandStringBytes(16)
		claims["user"] = userString
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		fmt.Println(wwwAuthentication, authorization)
		t, err := token.SignedString([]byte(Envs["ACCESS_TOKEN_SECRET"]))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.Set("WWW-Authentication", "gtoken")
		c.Set("Authorization", t)
		redisErr := GtokenClient.Set(RedisCtx, userString, t, 0).Err()
		CheckRedisErr(c, redisErr)
		return c.Next()
		//return c.JSON(fiber.Map{"token": t, "type": "gtoken"})
	}

	// placeholder code:
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = "placeholder"
	t, _ := token.SignedString([]byte(Envs["ACCESS_TOKEN_SECRET"]))
	return c.JSON(fiber.Map{"token": t})
}

func CreateUser(c *fiber.Ctx) error {
	wwwAuthentication := c.Get("WWW-Authentication")
	authorization := c.Get("Authorization")
	switch wwwAuthentication {
	case "upass":
		if authorization == "none" {
			username := c.Cookies("username")
			password := c.Cookies("password")
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
						"    https://www.techsales.dev/verify/" + link + "\r\n")

					// Sending email.
					err := smtp.SendMail(smtpHost+":"+smtpPort, auth, Envs["SMTP_ACCOUNT"], to, msg)
					if err != nil {
						fmt.Println("Error sending email", err)
					}
					fmt.Println("Email Sent Successfully!")
				}
			}
			// Else, user needs to supply username/password
		}
	}
	return c.SendStatus(200)
}

func VerifyUserLogin(c *fiber.Ctx) error {
	wwwAuthentication := c.Get("WWW-Authentication")
	authorization := c.Get("Authorization")
	switch wwwAuthentication {
	case "upass":
		if authorization == "none" {
			username := c.Cookies("username")
			password := c.Cookies("password")
			if username != "" && password != "" {
				var cursor uint64
				keys, cursor, err := UserClient.Scan(RedisCtx, cursor, "*", 1000000).Result()
				if err != nil {
					fmt.Println("something wrong with keys")
					panic(err)
				}
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
							userString := RandStringBytes(16)
							claims["user"] = userString
							claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
							fmt.Println(wwwAuthentication, authorization)
							t, err := token.SignedString([]byte(Envs["ACCESS_TOKEN_SECRET"]))
							if err != nil {
								return c.SendStatus(fiber.StatusInternalServerError)
							}
							c.Set("WWW-Authentication", "token")
							c.Set("Authorization", t)
							redisErr := UserTokensClient.Set(RedisCtx, userString, t, 0).Err()
							CheckRedisErr(c, redisErr)
						}
					}
				}
			}
			// Else, user needs to supply username/password
			return c.JSON(map[string]string{"status": "Invalid Username or Password"})
		}
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
			EmailPending.Del(RedisCtx, key)
			c.Status(200)
			return c.JSON(map[string]string{"status": "user created"})
		}
	}
	c.Status(403)
	return c.JSON(map[string]string{"status": "Forbidden"})
}
