package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/smtp"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/sha3"
)

func CheckRedisErr(c *fiber.Ctx, e error) error {
	if e != nil {
		fmt.Println("REDIS ERROR", e)
		fmt.Println("Make sure to set the environment variable: REDIS_PASS")
		fmt.Println("Make sure the backend REDIS_ENDPOINT is pointed correctly (127.0.0.1:6379?)")
		return &fiber.Error{Code: fiber.ErrInternalServerError.Code, Message: fiber.ErrInternalServerError.Message}
	}
	return nil
}

func CheckRedisErrWithoutContext(e error) {
	if e != nil {
		fmt.Println("REDIS ERROR", e)
		fmt.Println("Make sure to set the environment variable: REDIS_PASS")
		fmt.Println("Make sure the backend REDIS_ENDPOINT is pointed correctly (127.0.0.1:6379?)")
		panic(e)
	}
}

// Post request for when user email is forgotten
func ForgotPassword(c *fiber.Ctx) error {
	fmt.Println("ForgotPassword endpoint")
	emailBody := JustEmail{}
	email := ForgotPasswordEmailStruct{}
	email.Iat = uint64(time.Now().Unix())
	email.Link = RandStringBytes(32)
	err := c.BodyParser(&emailBody)
	if err != nil {
		return c.JSON(map[string]string{"status": "You need to POST an email with the request."})
	}
	email.Email = emailBody.Email

	// make sure the user exists
	iEmailExists, err := UserClient.Exists(RedisCtx, email.Email).Result()
	CheckRedisErr(c, err)
	if iEmailExists == 0 {
		return c.JSON(map[string]string{"status": "There is no account for " + email.Email})
	}

	forgotPasswordEmailStructAsString, err := json.Marshal(email)
	err = ForgotPasswordsClient.Set(RedisCtx, email.Email, forgotPasswordEmailStructAsString, 0).Err()
	CheckRedisErr(c, err)

	if DisableSendingEmail == false {
		// Send email to newly created user
		to := []string{email.Email}

		// smtp server configuration.
		smtpHost := "smtp.gmail.com"
		smtpPort := "587"
		auth := smtp.PlainAuth("", Envs["SMTP_ACCOUNT"], Envs["SMTP_PASS"], smtpHost)

		msg := []byte("To: " + email.Email + "\r\n" +
			"Subject: TechSales.dev Reset Password Request was made\r\n" +
			"\r\n" +
			"Your password for TechSales.dev was requested for reset.\r\n" +
			"Click here to verify your reset your password: \r\n" +
			"    https://www.techsales.dev/forgotPassword/" + "PLACEHOLDER" + "\r\n" +
			"if this was not you, you can ignore this email.\r\n")

		// Sending email.
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, Envs["SMTP_ACCOUNT"], to, msg)
		if err != nil {
			fmt.Println("Error sending email", err)
		}
		fmt.Println("Email Sent Successfully!")
		return c.JSON(map[string]string{"status": "email sent to " + email.Email + "."})
	}
	return c.JSON(map[string]string{"status": "email sent to " + email.Email + ". (THOUGH handlers.DisableSendingEmail is false [so not really sent])"})
}

func ForgotPasswordLinkGet(c *fiber.Ctx) error {
	link := c.Params("link", "")
	if link != "" {
		var cursor uint64
		keys, cursor, err := ForgotPasswordsClient.Scan(RedisCtx, cursor, "*", 1000000).Result()
		CheckRedisErr(c, err)
		// not efficient, but good for now
		for _, key := range keys {
			v, e := ForgotPasswordsClient.Get(RedisCtx, key).Result()
			CheckRedisErr(c, e)
			var f ForgotPasswordEmailStruct
			json.Unmarshal([]byte(v), &f)
			if f.Link == link {
				return c.JSON(map[string]string{"status": "ready for a reset", "post": link})
			}
		}
	}
	return c.JSON(map[string]string{"status": "Nothing Here"})
}

func CheckPasswordRegex(c *fiber.Ctx, s *string) bool {
	match, _ := regexp.Match("[a-zA-Z0-9$@#$%\"\\^&\\*(){}\\[\\]_+|=~-]{8,}", []byte(*s))
	return match
}

func ForgotPasswordLinkPost(c *fiber.Ctx) error {
	link := c.Params("link", "")
	if link != "" {
		var cursor uint64
		keys, cursor, err := ForgotPasswordsClient.Scan(RedisCtx, cursor, "*", 1000000).Result()
		CheckRedisErr(c, err)
		// not efficient, but good for now
		for _, key := range keys {
			v, e := ForgotPasswordsClient.Get(RedisCtx, key).Result()
			CheckRedisErr(c, e)
			var f ForgotPasswordEmailStruct
			json.Unmarshal([]byte(v), &f)

			if f.Link == link {
				u := PasswordChecker{}
				u.Email = f.Email
				u.EmailVerification = true

				err = c.BodyParser(&u)
				if err != nil {
					fmt.Println("body parsing err", err)
					c.Status(400)
					return c.JSON(map[string]string{"status": "Bad Request"})
				} else {
					//Check Both Passwords
					if u.Password1 == u.Password2 && CheckPasswordRegex(c, &u.Password1) {
						// Change Password
						u2 := User{Email: f.Email, EmailVerification: true}
						h := sha3.New512()
						h.Write([]byte(u.Password1))
						hexStr := hex.EncodeToString(h.Sum(nil))
						u2.Password = hexStr
						marshaledUser, _ := json.Marshal(u)
						err = UserClient.Set(RedisCtx, f.Email, marshaledUser, 0).Err()
						CheckRedisErr(c, err)
						// Remove the password link from redis Db: 5
						err = ForgotPasswordsClient.Del(RedisCtx, f.Email).Err()
						CheckRedisErr(c, err)
						return c.JSON(map[string]string{"status": "success"})
					} else {
						// Why in the world wouldn't the passwords match?
						// Why in the world wouldn't the password pass the regular expression match?
						// The request is NOT coming from our frontend, so
						// the request is nefarious, and therefore the user is "attempting to brew coffee with a teapot"
						return &fiber.Error{Code: fiber.ErrTeapot.Code, Message: "The server refuses the attempt to brew coffee with a teapot."}
					}
				}
			}
		}
	}
	return c.JSON(map[string]string{"status": "Nothing Here"})
}
