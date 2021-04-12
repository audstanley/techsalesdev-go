package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Categories(c *fiber.Ctx) error {
	cat := c.Params("category", "")
	fmt.Println("category", cat+".")
	if cat == "" {
		return fiber.ErrBadRequest
	}
	if cat == "pcb" {
		a := make([]Product, 0)
		var cursor uint64
		keys, cursor, err := PCBClient.Scan(RedisCtx, cursor, "*", 100000).Result()
		CheckRedisErr(c, err)
		for _, key := range keys {
			v, err := PCBClient.Get(RedisCtx, key).Result()
			CheckRedisErr(c, err)

			var p Product
			json.Unmarshal([]byte(v), &p)
			a = append(a, p)
		}
		return c.JSON(a)
	} else if cat == "wires" {
		a := make([]Product, 0)
		var cursor uint64
		keys, cursor, err := WiresClient.Scan(RedisCtx, cursor, "*", 100000).Result()
		CheckRedisErr(c, err)
		for _, key := range keys {
			v, err := WiresClient.Get(RedisCtx, key).Result()
			CheckRedisErr(c, err)

			var p Product
			json.Unmarshal([]byte(v), &p)
			a = append(a, p)
		}
		return c.JSON(a)
	} else if cat == "diodes" {
		a := make([]Product, 0)
		var cursor uint64
		keys, cursor, err := DiodesClient.Scan(RedisCtx, cursor, "*", 100000).Result()
		CheckRedisErr(c, err)
		for _, key := range keys {
			v, err := DiodesClient.Get(RedisCtx, key).Result()
			CheckRedisErr(c, err)

			var p Product
			json.Unmarshal([]byte(v), &p)
			a = append(a, p)
		}
		return c.JSON(a)
	} else if cat == "caps" {
		a := make([]Product, 0)
		var cursor uint64
		keys, cursor, err := CapsClient.Scan(RedisCtx, cursor, "*", 100000).Result()
		CheckRedisErr(c, err)
		for _, key := range keys {
			v, err := CapsClient.Get(RedisCtx, key).Result()
			CheckRedisErr(c, err)

			var p Product
			json.Unmarshal([]byte(v), &p)
			a = append(a, p)
		}
		return c.JSON(a)
	}
	return c.JSON(fiber.Map{"status": "nothing here"})
}
