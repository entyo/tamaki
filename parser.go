package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"io/ioutil"
	"time"
)

func fetchGinzaTravelInfo(ch chan map[string]string) {
	// http.Request の生成
	url := "http://www.tokyometro.jp/unkou/history/ginza.html"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	//  行くぞーーー！！！
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

		// bodyの確認
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("body: " + string(body) + "\n")

		// goqueryでがんばる
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatalln(err)
		}

		// 偏差値の低い型 map[string]string
		var travelInfos []map[string]string
		doc.Find("tr").Each(func(i int, s *goquery.Selection) {
			// 日付時刻
			th := s.Find("th").Text()
			// 内容
			td := s.Find("td").Text()

			if th != "" && td != "" {
				travelInfos = append(travelInfos, map[string]string{th: td})
			}
		})

		//  最新の情報をチャネルで送る
		if len(travelInfos) > 0 {
			ch <- travelInfos[0]
		}

	}
}
