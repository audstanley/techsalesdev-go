package handlers

import (
	"context"
	"encoding/json"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gofiber/fiber/v2"
)

func WeiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

// this is just for testing purposes
func Wallets(c *fiber.Ctx) error {
	ctx := context.Background()
	client, err := ethclient.Dial(Envs["ETHERIUM_NETWORK"])
	if err != nil {
		log.Fatal(err)
	}

	a := make([]EtheriumWallet, 0)
	var cursor uint64
	keys, cursor, err := UserWalletsClient.Scan(RedisCtx, cursor, "*", 100000).Result()
	CheckRedisErr(c, err)
	for _, key := range keys {
		v, err := UserWalletsClient.Get(RedisCtx, key).Result()
		CheckRedisErr(c, err)

		var w EtheriumWallet
		amount, _ := client.BalanceAt(ctx, common.HexToAddress("0x051Fc21738F3CDE0F6370AF3860CF8eAA617E3B4"), nil)
		w.Ballance = *(WeiToEther(amount))
		json.Unmarshal([]byte(v), &w)
		a = append(a, w)
	}

	return c.JSON(a)
}
