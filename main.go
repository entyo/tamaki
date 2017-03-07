package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

func getSlackAPIToken() string {
	return os.Getenv("SLACK_API_TOKEN")
}

func run(api *slack.Client) int {
	// RTMコネクションを張っておく
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// 一定周期でGRT->スクレイピングをする
	var ginzaTravelInfo = make(chan map[string]string)
	go fetchGinzaTravelInfo(ginzaTravelInfo)

	for {
		select {
		case info := <-ginzaTravelInfo:
			log.Println(info)
		}

		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Println("Hello Event")

			case *slack.MessageEvent:
				log.Printf("Message: %v\n", ev)
				// rtm.SendMessage(rtm.NewOutgoingMessage("基礎知識なさすぎ", ev.Channel))

			case *slack.InvalidAuthEvent:
				log.Println("Invalid credentials")
				return 1
			}
		}
	}
}

func main() {
	api := slack.New(getSlackAPIToken())
	os.Exit(run(api))
}
