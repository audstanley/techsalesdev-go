package handlers

import (
	redis "github.com/go-redis/redis/v8"
)

// hash : jwt
var GtokenClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         0,
	MaxRetries: 3,
})

// hash : jwt
var UserClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         1,
	MaxRetries: 3,
})

// md5Hash : { OnSale : []Product, NewArrivals : []Product }
var ProductsClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         2,
	MaxRetries: 3,
})

// md5Hash : { OnSale : []Product, NewArrivals : []Product }
var EmailPending = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         3,
	MaxRetries: 3,
})

var UserTokensClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         4,
	MaxRetries: 3,
})

var ForgotPasswordsClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         5,
	MaxRetries: 3,
})

var UserAddressesClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         6,
	MaxRetries: 3,
})

var PCBClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         7,
	MaxRetries: 3,
})

var WiresClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         8,
	MaxRetries: 3,
})

var DiodesClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         9,
	MaxRetries: 3,
})

var CapsClient = redis.NewClient(&redis.Options{
	Addr:       Envs["REDIS_ENDPOINT"],
	Password:   Envs["REDIS_PASS"],
	DB:         10,
	MaxRetries: 3,
})
