package handlers

import "github.com/gofiber/fiber/v2"

func SignUp(c *fiber.Ctx) error {

	return c.JSON(map[string]string{"status": "PLACEHOLDER"})
}
