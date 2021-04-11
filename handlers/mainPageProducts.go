package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// GET /mainProductsPage
func MainProductPage(c *fiber.Ctx) error {
	fmt.Println("mainProductsPage endpoint")
	var cursor uint64
	keys, cursor, err := ProductsClient.Scan(RedisCtx, cursor, "*", 1000000).Result()
	CheckRedisErr(err)

	var onSale []Product
	var newArrivals []Product
	for _, key := range keys {
		// we need to query redis for the object from the key.
		v, e := ProductsClient.Get(RedisCtx, key).Result()
		CheckRedisErr(e)
		var p Product
		json.Unmarshal([]byte(v), &p)
		if p.OnSale {
			onSale = append(onSale, p)
		} else {
			newArrivals = append(newArrivals, p)
		}
	}
	var products ProductReturn
	products.OnSale = make([]Product, 0)
	products.NewArrivals = make([]Product, 0)
	products.OnSale = onSale
	products.NewArrivals = newArrivals
	return c.JSON(products)
}
