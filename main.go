package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	RedisAddr     string `default:"localhost:6379" split_words:"true"`
	RedisPassword string `default:"" split_words:"true"`
	RedisDB       int    `default:"0" split_words:"true"`
	ChannelName   string `default:"alerts" split_words:"true"`
	Debug         bool   `default:"false" split_words:"true"`
}

var (
	redisClient *redis.Client
	conf        config
)

func init() {

	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal(err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("%s connected", redisClient)
	}

	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {

	r := gin.Default()
	r.POST("/webhook", webhook)
	r.Run()
}

func webhook(c *gin.Context) {

	dataBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
	}
	dataString := string(dataBytes)

	if conf.Debug {
		log.Println(dataString)
	}

	err = redisClient.Publish(context.Background(), conf.ChannelName, dataString).Err()
	if err != nil {
		log.Println(err)
	}

}
