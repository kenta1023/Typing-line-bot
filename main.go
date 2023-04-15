package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	http.HandleFunc("/", lineHandler)
	fmt.Println("start http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func lineHandler(w http.ResponseWriter, r *http.Request) {
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(r)
	if err != linebot.ErrInvalidSignature {
		w.WriteHeader(400)
	} else {
		w.WriteHeader(500)
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			case *linebot.TextMessage:
				if message.Text == "フリック入力" {
					sendQuestion(event, bot)
				}
			}

		}
	}
}

// 問題を送信
func sendQuestion(event *linebot.Event, bot *linebot.Client) {
	question := generateQuestion()
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(question)).Do()
	if err != nil {
		log.Print(err)
	}
}

// ランダムで問題を選択
func generateQuestion() string {
	questionList := []string{"テスト", "今日の天気は晴れのち曇り", "今日も良い一日になりそうだ", "毎日のラジオ体操が日課です"}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := rand.Intn(len(questionList))

	return questionList[randomNum]
}
