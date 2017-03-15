package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"strings"

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

	// stubをsomeTravelInfo([]MetroTravelInfo)にセットする
	var someTravelInfo []MetroTravelInfo
	metroData := makeMetroData()
	for _, railwayLine := range metroData {
		travelInfo := MetroTravelInfo{
			railwayLine: railwayLine,
		}
		someTravelInfo = append(someTravelInfo, travelInfo)
	}
	// 5秒毎にスクレイピングする
	duration := 5 * time.Second
	metroCh := collectMetroTravelInfo(someTravelInfo, duration)

	// Real Time Messaging APIとの接続をgoroutineで持っとく
	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		// Slack RTM API
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				if strings.Contains(ev.Text, "@"+ev.BotID) {
					reply := getRandomReply()
					rtm.SendMessage(rtm.NewOutgoingMessage(reply, ev.Channel))
				}
			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}

		// スクレイピング
		case latestStatus := <-metroCh:
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
			err = postMessageToAll(slackClient, latestStatus.railwayLine.name+"の運行情報", msg, latestStatus.railwayLine.colorCode)
			if err != nil {
				log.Println("Error in posting message to slack: ", err)
			}
		}
	}
}

func init() {
	slackClient := slack.New(getSlackAPIToken())
	os.Exit(run(slackClient))
}
