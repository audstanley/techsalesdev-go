package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// GET /john
func UserList(c *fiber.Ctx) error {
	msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
	return c.SendString(msg) // => Hello john 👋!
}
