package main

import (
	"fmt"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/nlopes/slack"
)

func run(slackClient *slack.Client) int {
	// Redisに繋ぐ
	c, err := redis.DialURL(getRedisURL())
	if err != nil {
		log.Println("Error in dialing redis")
	} else {
		log.Printf("Connected to redis(address: %s)", getRedisURL())
	}
	defer c.Close()

	// 一度全てのキャッシュを消去する
	if reply, err := c.Do("FLUSHALL"); err != nil {
		log.Println(reply, err)
		return 1
	} else {
		log.Println("FLUSHALL: ", reply)
	}

	// 一定周期でスクレイピングをする
	ginzaTravelInfo := make(chan GinzaTravelInfo)
	go fetchGinzaTravelInfo(ginzaTravelInfo)

	for {
		select {
		case latestStatus := <-ginzaTravelInfo:
			s, err := redis.String(c.Do("GET", latestStatus.dateTime))
			// 最新の運行情報が既に保存済みか
			if err == nil {
				fmt.Printf("redis GET: %#v\n", s)
				continue
			}
			// redisにSET
			res, err := c.Do("SET", latestStatus.dateTime, latestStatus.content)
			if err != nil {
				log.Println("Error in setting travel info from radis: ", err)
				return 1
			}
			log.Println("Info set to radis: ", res)
			log.Println("New info will be send to slack...")

			// joinしているchannel全てにmessageを送る
			msg := latestStatus.dateTime + " : " + latestStatus.content
			err = postMessageToAll(slackClient, "銀座線の運行情報", msg, "#FF932E")
			if err != nil {
				log.Println("Error in posting message to slack: ", err)
			}
		}
	}
}

func main() {
	slackClient := slack.New(getSlackAPIToken())
	os.Exit(run(slackClient))
}
