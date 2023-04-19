package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type data struct {
	id            string
	sentQuestion  string
	correctNumber int
	sendTime      time.Time
}

var sentData map[string]*data

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
				} else {
					checkAnswer(event, bot, message)
				}
			}

		}
	}
}

// 問題を送信 送信済みの問題を保存(saveSentData)
func sendQuestion(event *linebot.Event, bot *linebot.Client) {
	question := generateQuestion()
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(question)).Do()
	if err != nil {
		log.Print(err)
	}
	saveSentData(event, question)
}

// 送信済み問題を保存
func saveSentData(event *linebot.Event, question string) {
	sentData = make(map[string]*data)
	var id string

	//id取得したい
	if event.Source.Type == "group" {
		id = event.Source.GroupID
	} else {
		id = event.Source.UserID
	}
	//idがすでに存在　削除する
	delete(sentData, id)
	sentData[id] = &data{id: id, sentQuestion: question, correctNumber: 0, sendTime: time.Now()}
	fmt.Print("保管・送信済:")
	fmt.Println(sentData[id])
}

// ランダムで問題を選択
func generateQuestion() string {
	questionList := []string{"テスト", "今日の天気は晴れのち曇り", "今日も良い一日になりそうだ", "毎日のラジオ体操が日課です"}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := rand.Intn(len(questionList))

	return questionList[randomNum]
}

// 　回答をチェック
func checkAnswer(event *linebot.Event, bot *linebot.Client, message *linebot.TextMessage) {
	var id string
	var userID string
	var replyMessage string
	typeUserOrGroup := event.Source.Type

	if typeUserOrGroup == "group" {
		id = event.Source.GroupID
		userID = event.Source.UserID
	} else {
		id = event.Source.UserID
	}

	if data, ok := sentData[id]; ok {
		if data.sentQuestion == message.Text {
			// 正解の返答が来たあとの処理
			time := time.Since(data.sendTime)
			FormattedTime := fmt.Sprintf("%.3f秒", time.Seconds())
			//　userの場合の返信
			if typeUserOrGroup == "user" {
				replyMessage = "タイム:" + FormattedTime
			} else { //groupの場合
				//ユーザ名（ディスプレイ名）取得
				res, err := bot.GetProfile(userID).Do()
				if err != nil {
					log.Print(err)
				}
				data.correctNumber++
				name := res.DisplayName
				replyMessage = strconv.Itoa(data.correctNumber) + "位" + name + "\nタイム:" + FormattedTime
			}
			//データ送信
			_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
			if err != nil {
				log.Print(err)
			}
			fmt.Print("正解:")
			fmt.Println(sentData[id])

		} else {
			fmt.Println("関係のない会話 or 不正解:" + id)
		}
	} else {
		fmt.Println("関係のない会話 or 不正解:" + id)
	}
}
