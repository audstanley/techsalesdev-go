package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func ProductImage(c *fiber.Ctx) error {
	productId := c.Params("productId", "")
	pStr, err := ProductsClient.Get(RedisCtx, productId).Result()
	CheckRedisErr(c, err)
	if pStr != "" {
		p := Product{}
		json.Unmarshal([]byte(pStr), &p)
		return c.JSON(p)
	}
	return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: fiber.ErrBadRequest.Message}
}

func AddProduct(c *fiber.Ctx) error {
	productId := c.Params("productId", "")
	fmt.Println(productId)

	token, err := jwt.Parse(c.Get("Authorization", ""), func(token *jwt.Token) (interface{}, error) {
		return []byte(Envs["ACCESS_TOKEN_SECRET"]), nil
	})
	if err != nil {
		fmt.Println("token err", err)
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: fiber.ErrBadRequest.Message}
	}
	claims := token.Claims.(jwt.MapClaims)
	claimUserStr, _ := claims["user"].(string)
	fmt.Println("add product from user:", claimUserStr)

	i, err := UserClient.Exists(RedisCtx, claimUserStr).Result()
	CheckRedisErr(c, err)
	if i == 0 {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: fiber.ErrBadRequest.Message}
	}
	// get the product
	pStr, err := ProductsClient.Get(RedisCtx, productId).Result()
	CheckRedisErr(c, err)
	// get the user's cart
	cartStr, err := UsersCartClient.Get(RedisCtx, claimUserStr).Result()
	CheckRedisErr(c, err)
	fmt.Println("cart from redis:", cartStr)
	p := ProductForCart{Amount: 1}
	json.Unmarshal([]byte(pStr), &p)
	cart := Cart{}
	json.Unmarshal([]byte(cartStr), &cart)

	productFound := false
	for i, k := range cart.Products {
		if k.ProductId == productId {
			cart.Products[i].Amount += 1
			productFound = true
		}
	}
	if !productFound {
		cart.Products = append(cart.Products, p)
	}
	// update the total cost.
	cart.Total += p.Cost
	fmt.Println("newCart", cart)
	jm, _ := json.Marshal(cart)
	// put the cart back into the database
	err = UsersCartClient.Set(RedisCtx, claimUserStr, string(jm), 0).Err()
	return c.JSON(cart)
}

func RemoveProduct(c *fiber.Ctx) error {
	productId := c.Params("productId", "")
	fmt.Println(productId)

	token, err := jwt.Parse(c.Get("Authorization", ""), func(token *jwt.Token) (interface{}, error) {
		return []byte(Envs["ACCESS_TOKEN_SECRET"]), nil
	})
	if err != nil {
		fmt.Println("token err", err)
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: fiber.ErrBadRequest.Message}
	}
	claims := token.Claims.(jwt.MapClaims)
	claimUserStr, _ := claims["user"].(string)
	fmt.Println("remove product from user:", claimUserStr)

	i, err := UserClient.Exists(RedisCtx, claimUserStr).Result()
	CheckRedisErr(c, err)
	if i == 0 {
		return &fiber.Error{Code: fiber.ErrForbidden.Code, Message: fiber.ErrForbidden.Message}
	}
	// get the product
	pStr, err := ProductsClient.Get(RedisCtx, productId).Result()
	CheckRedisErr(c, err)
	// get the user's cart
	cartStr, err := UsersCartClient.Get(RedisCtx, claimUserStr).Result()
	CheckRedisErr(c, err)
	fmt.Println("cart from redis:", cartStr)
	p := ProductForCart{Amount: 1}
	json.Unmarshal([]byte(pStr), &p)
	cart := Cart{}
	json.Unmarshal([]byte(cartStr), &cart)

	productFound := false
	for i, k := range cart.Products {
		if k.ProductId == productId {
			if cart.Products[i].Amount > 0 {
				cart.Products[i].Amount -= 1
			}
			if cart.Products[i].Amount == 0 {
				// slower, but maintains order
				copy(cart.Products[i:], cart.Products[i+1:])
				cart.Products[len(cart.Products)-1] = ProductForCart{}
				cart.Products = cart.Products[:len(cart.Products)-1]
			}
			productFound = true
		}
	}
	if productFound {
		// update the total cost.
		cart.Total -= p.Cost
	}
	fmt.Println("newCart", cart)
	jm, _ := json.Marshal(cart)
	// put the cart back into the database
	err = UsersCartClient.Set(RedisCtx, claimUserStr, string(jm), 0).Err()
	return c.JSON(cart)
}

func GetCart(c *fiber.Ctx) error {
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
	cart := Cart{}
	json.Unmarshal([]byte(cartStr), &cart)
	return c.JSON(cart)
}
