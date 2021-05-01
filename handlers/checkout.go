package handlers

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func GetConfirmationStatus(c *fiber.Ctx) error {
	token, err := jwt.Parse(c.Get("Authorization", ""), func(token *jwt.Token) (interface{}, error) {
		return []byte(Envs["ACCESS_TOKEN_SECRET"]), nil
	})
	if err != nil {
		fmt.Println("token err", err)
		return &fiber.Error{Code: fiber.ErrForbidden.Code, Message: fiber.ErrForbidden.Message}
	}
	claims := token.Claims.(jwt.MapClaims)
	claimUserStr, _ := claims["user"].(string)

	cartStr, err := ConfirmationCodesClient.Get(RedisCtx, c.Params("code", "")).Result()
	cart := Cart{}
	json.Unmarshal([]byte(cartStr), cart)
	// if the user isn't logged in, they wont be able to check their order
	if claimUserStr == cart.Email {
		return c.JSON(cart)
	}
	return &fiber.Error{Code: fiber.ErrForbidden.Code, Message: fiber.ErrForbidden.Message}
}

func Checkout(c *fiber.Ctx) error {
	token, err := jwt.Parse(c.Get("Authorization", ""), func(token *jwt.Token) (interface{}, error) {
		return []byte(Envs["ACCESS_TOKEN_SECRET"]), nil
	})
	if err != nil {
		fmt.Println("token err", err)
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: fiber.ErrBadRequest.Message}
	}
	claims := token.Claims.(jwt.MapClaims)
	claimUserStr, _ := claims["user"].(string)

	i, err := UserClient.Exists(RedisCtx, claimUserStr).Result()
	CheckRedisErr(c, err)
	if i == 0 {
		return &fiber.Error{Code: fiber.ErrForbidden.Code, Message: fiber.ErrForbidden.Message}
	}
	// get the user's cart
	cartStr, err := UsersCartClient.Get(RedisCtx, claimUserStr).Result()
	CheckRedisErr(c, err)
	fmt.Println("cart from redis:", cartStr)
	cart := Cart{Placed: true} // make the order placed
	// Make a confirmation code for the completed order
	cart.ConfirmationCode = RandStringBytes(32)
	// get the user's shipping info
	userShipping := UserShipping{}
	usersAddresses, err := UserAddressesClient.Get(RedisCtx, claimUserStr).Result()
	CheckRedisErr(c, err)
	json.Unmarshal([]byte(usersAddresses), &userShipping)
	// add shipping into to cart
	cart.ShipTo = userShipping
	cart.Email = claimUserStr
	json.Unmarshal([]byte(cartStr), &cart)
	usersOrdersStr, err := OrdersClient.Get(RedisCtx, claimUserStr).Result()
	CheckRedisErr(c, err)
	// add the cart to the orders

	o := Orders{}
	json.Unmarshal([]byte(usersOrdersStr), o)
	o.Orders = append(o.Orders, cart)
	jmOrders, _ := json.Marshal(o)
	// save the updated orders to the database
	err = OrdersClient.Set(RedisCtx, claimUserStr, string(jmOrders), 0).Err()
	CheckRedisErr(c, err)

	// add the confirmation code to the database
	jmCart, _ := json.Marshal(cart)
	err = ConfirmationCodesClient.Set(RedisCtx, cart.ConfirmationCode, string(jmCart), 0).Err()

	// email the user the order
	if !DisableSendingEmail {
		// Send email to user
		to := []string{claimUserStr}

		// smtp server configuration.
		smtpHost := "smtp.gmail.com"
		smtpPort := "587"
		auth := smtp.PlainAuth("", Envs["SMTP_ACCOUNT"], Envs["SMTP_PASS"], smtpHost)

		msg := []byte("To: " + claimUserStr + "\r\n" +
			"Subject: TechSales.dev Order has been placed\r\n" +
			"\r\n" +
			"You placed an order on our website.\r\n" +
			"You can check your order her: \r\n" +
			"    https://api.techsales.dev/confirmationCode/" + cart.ConfirmationCode + "\r\n")

		// Sending email.
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, Envs["SMTP_ACCOUNT"], to, msg)
		if err != nil {
			fmt.Println("Error sending email", err)
		}
		fmt.Println("Email Sent Successfully!")
	}

	// return the orders - the last order in the list should have the confirmation number.
	return c.JSON(o)
}
