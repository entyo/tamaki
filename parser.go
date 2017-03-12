package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"strings"
	"time"
)

// GinzaTravelInfo は銀座線の運行情報を表す
type GinzaTravelInfo struct {
	dateTime string
	content  string
}

func fetchGinzaTravelInfo(ch chan GinzaTravelInfo) {
	// http.Request の生成
	const url = "http://www.tokyometro.jp/unkou/history/ginza.html"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	// 行くぞーーー！！！
	for {
		time.Sleep(5 * time.Second)

		// リクエストを投げる
		resp, err := client.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			log.Fatalln(err)
		}

		// goqueryでがんばる
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatalln(err)
		}

		var travelInfo []GinzaTravelInfo
		doc.Find("tr").Each(func(i int, s *goquery.Selection) {
			// 日付時刻
			th := s.Find("th").Text()
			// 内容
			td := s.Find("td").Text()

			if th != "" && td != "" {
				travelInfo = append(travelInfo, GinzaTravelInfo{dateTime: th, content: td})
			}
		})

		//  最新の情報をチャネルで送る
		if len(travelInfo) > 0 {
			latest := travelInfo[0]
			latest.dateTime = strings.TrimSpace(latest.dateTime)
			latest.content = strings.TrimSpace(latest.content)
			log.Printf("Latest info from a website: %s %s\n", latest.dateTime, latest.content)
			ch <- latest
		}

	}
}
