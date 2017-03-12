package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/nlopes/slack"
)

func getRandomReply() string {
	switch rand.Intn(6) {
	case 0:
		return "ぴゃいっ！!"
	case 1:
		return "えっ あの なんで撫でるんですか?"
	case 2:
		return "そ そんなつらい顔しないで下さい…"
	case 3:
		return "わーーーっ ハードル上げないでよーーっ"
	case 4:
		return "わーっ パパお帰りっ"
	default:
		return "富田林田舎ちゃうもんっ…"
	}
}

func postMessage(client *slack.Client, channelID string, pretext string, text string, color string) (err error) {
	params := slack.PostMessageParameters{
		AsUser: true,
	}
	attachment := slack.Attachment{
		Color:   color,
		Pretext: pretext,
		Text:    text,
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := client.PostMessage(channelID, "", params)
	if err != nil {
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	return nil
}

func getSlackAPIToken() string {
	return os.Getenv("SLACK_API_TOKEN")
}

// botの所属しているchannel全てにmessageを送る
func postMessageToAll(client *slack.Client, pretext string, text string, color string) (err error) {
	groups, err := client.GetGroups(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	for _, group := range groups {
		err := postMessage(client, group.ID, pretext, text, color)
		if err != nil {
			return err
		}
	}
	return nil
}
