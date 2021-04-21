package main

import (
	"fmt"
	"log"
	"main/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "www-authentication, authorization, x-username, x-password",
		ExposeHeaders:    "*",
		AllowCredentials: true,
	}))

	// Sessions Handler
	app.Use(handlers.Session)
	//app.Post("/session", handlers.CreateUser)
	app.Post("/session", handlers.VerifyUserLogin)
	// Verify a user from SMTP
	app.Get("/verify/:link", handlers.Verify)
	// Main page products
	app.Get("/mainProductPage", handlers.MainProductPage)
	// When the user forgot a password
	app.Post("/forgotPassword", handlers.ForgotPassword)
	// Where to get reset password link from
	app.Get("/forgotPassword/:link", handlers.ForgotPasswordLinkGet)
	// Where to post to actually reset password
	app.Post("/forgotPassword/:link", handlers.ForgotPasswordLinkPost)
	// Where to fully signup
	app.Post("/signup", handlers.SignUpARealUser)
	// Get the cart for the user (must have jwt)
	app.Get("/cart", handlers.GetCart)
	// Checkout with a simple post request (must be a user)
	app.Post("/checkout", handlers.Checkout)
	// FOR TESTING ONLY!.  View the ballance for all walllets
	app.Get("/wallets", handlers.Wallets)
	// Checkout with a simple post request (must be a user)
	app.Post("/confirmationCode/:code", handlers.GetConfirmationStatus)
	// Products based on category
	app.Get("/:category", handlers.Categories)
	// Adding a product based on the id
	app.Get("/add/:productId", handlers.AddProduct)
	// Removing a product based on the id
	app.Get("/remove/:productId", handlers.RemoveProduct)
	// Get the product directly (for image) [this will be needed for the cart page]
	app.Get("/product/:productId", handlers.ProductImage)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", handlers.Envs["API_DEV_PORT"])))
}
