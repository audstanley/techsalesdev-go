package main

import (
	"fmt"
	"log"
	"main/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Must have a .products/photos/ path with photos in the path.
	// These photos will get populated in the redis database.
	handlers.PopulateDatabaseWithImages()
	// This function will add user: verifiedUser@gmail.com to the userDb
	handlers.AutoVerifyAUserInTheDatabase()
	if handlers.DisableSendingEmail {
		fmt.Println("SMTP Emailing is Disabled, to enable email,\n\tset handlers.DisableSendingEmail to: false - and rebuild application")
	}

	app := fiber.New()
	// Sessions Handler
	app.Use(handlers.Session)
	app.Post("/session", handlers.CreateUser)
	app.Get("/session", handlers.VerifyUserLogin)
	app.Get("/verify/:link", handlers.Verify)
	app.Get("/mainProductPage", handlers.MainProductPage)
	app.Post("/forgotPassword", handlers.ForgotPassword)
	app.Get("/forgotPassword/:link", handlers.ForgotPasswordLinkGet)
	app.Post("/forgotPassword/:link", handlers.ForgotPasswordLinkPost)
	app.Post("/signup", handlers.SignUp)

	// app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: []byte(envs["ACCESS_TOKEN_SECRET"]),
	// }))

	//app.Get("/:name", handlers.UserList)

	// GET /john/75
	// app.Get("/:namefiber/:age", func(c *fiber.Ctx) error {
	// 	msg := fmt.Sprintf("👴 %s is %s years old", c.Params("name"), c.Params("age"))
	// 	return c.SendString(msg) // => 👴 john is 75 years old
	// })

	// // GET /dictionary.txt
	// app.Get("/:file.:ext", func(c *fiber.Ctx) error {
	// 	msg := fmt.Sprintf("📃 %s.%s", c.Params("file"), c.Params("ext"))
	// 	return c.SendString(msg) // => 📃 dictionary.txt
	// })

	// // GET /flights/LAX-SFO
	// app.Get("/flights/:from-:to", func(c *fiber.Ctx) error {
	// 	msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
	// 	return c.SendString(msg) // => 💸 From: LAX, To: SFO
	// })

	// // GET /api/register
	// app.Get("/api/*", func(c *fiber.Ctx) error {
	// 	msg := fmt.Sprintf("✋ %s", c.Params("*"))
	// 	return c.SendString(msg) // => ✋ register
	// })

	log.Fatal(app.Listen(fmt.Sprintf(":%s", handlers.Envs["API_DEV_PORT"])))
}
