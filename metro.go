package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"strings"
	"time"
)

// MetroRailwayLine は東京メトロの路線情報を表す
type MetroRailwayLine struct {
	name               string
	colorCode          string
	operationStatusURL string
}

// MetroTravelInfo は http://www.tokyometro.jp/unkou/history/* から得られる運行情報を表す
type MetroTravelInfo struct {
	dateTime    string
	content     string
	railwayLine MetroRailwayLine
}

func makeMetroData() []MetroRailwayLine {
	metroData := []MetroRailwayLine{
		MetroRailwayLine{
			name:               "銀座線",
			colorCode:          "#FF9500",
			operationStatusURL: "http://www.tokyometro.jp/unkou/history/ginza.html",
		},
		MetroRailwayLine{
			name:               "丸の内線",
			colorCode:          "#F62E36",
			operationStatusURL: "http://www.tokyometro.jp/unkou/history/marunouchi.html",
		},
		MetroRailwayLine{
			name:               "日比谷線",
			colorCode:          "#B5B5AC",
			operationStatusURL: "http://www.tokyometro.jp/unkou/history/hibiya.html",
		},
		MetroRailwayLine{
			name:               "東西線",
			colorCode:          "#009BBF",
			operationStatusURL: "http://www.tokyometro.jp/unkou/history/touzai.html",
		},
		MetroRailwayLine{
			name:               "千代田線",
			colorCode:          "#00BB85",
			operationStatusURL: "http://www.tokyometro.jp/unkou/history/chiyoda.html",
		},
		MetroRailwayLine{
			name:               "有楽町線",
			colorCode:          "#C1A470",
			operationStatusURL: "http://www.tokyometro.jp/unkou/history/yurakucho.html",
		},
		MetroRailwayLine{
			name:               "半蔵門線",
			colorCode:          "#8F76D6",
			operationStatusURL: "http://www.tokyometro.jp/unkou/history/hanzoumon.html",
		},
		MetroRailwayLine{
			name:               "南北線",
			colorCode:          "#00AC9B",
			operationStatusURL: "http://www.tokyometro.jp/unkou/history/nanboku.html",
		},
		MetroRailwayLine{
			name:               "副都心線",
			colorCode:          "#9C5E31",
			operationStatusURL: "http://www.tokyometro.jp/unkou/history/fukutoshin.html",
		},
	}

	return metroData
}

// 複数のWebページそれぞれについて、運行情報が更新されたらある一つのchannelに流すgoroutineをつくる
func collectMetroTravelInfo(someInfo []MetroTravelInfo, duration time.Duration) chan MetroTravelInfo {
	c := make(chan MetroTravelInfo)
	for _, info := range someInfo {
		go updateMetroTravelInfo(c, info, duration)
	}

	return c
}

func updateMetroTravelInfo(ch chan MetroTravelInfo, travelInfo MetroTravelInfo, interval time.Duration) {
	// http.Request の生成
	req, err := http.NewRequest(http.MethodGet, travelInfo.railwayLine.operationStatusURL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	const agent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"
	req.Header.Set("User-Agent", agent)

	// 行くぞーーー！！！
	for {
		time.Sleep(interval)

		// リクエストを投げる
		client := &http.Client{}
		resp, err := client.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			log.Fatalln(err)
		}

		// goqueryでがんばる
		doc, err := goquery.NewDocument(travelInfo.railwayLine.operationStatusURL)
		if err != nil {
			log.Fatalln(err)
		}
		var someInfo []MetroTravelInfo
		doc.Find("tr").Each(func(i int, s *goquery.Selection) {
			// 日付時刻
			th := s.Find("th").Text()
			// 内容
			td := s.Find("td").Text()

			if th != "" && td != "" {
				newInfo := MetroTravelInfo{
					dateTime:    th,
					content:     td,
					railwayLine: travelInfo.railwayLine,
				}
				someInfo = append(someInfo, newInfo)
			}
		})

		//  最新の情報をチャネルで送る
		if len(someInfo) > 0 {
			latest := someInfo[0]
			latest.dateTime = strings.TrimSpace(latest.dateTime)
			latest.content = strings.TrimSpace(latest.content)
			log.Printf("Latest info from a website(%s): %s %s\n", latest.railwayLine.name, latest.dateTime, latest.content)
			ch <- latest
		}

	}
}
