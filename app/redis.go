package app

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisClient(cnf *config.Config) *redis.Client {
	host := cnf.Env.GetString("REDIS_HOST")
	port := cnf.Env.GetString("REDIS_PORT")
	db := cnf.Env.GetInt("REDIS_DB")

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		DB:   db,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("error while connect to redis", err)
	}

	log.Println("Connected to redis")
	return client
}
