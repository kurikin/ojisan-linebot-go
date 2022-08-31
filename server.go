package main

import (
	"log"
	"net/http"
	"os"

	"github.com/greymd/ojichat/generator"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Config struct {
	TargetName       string `docopt:"<name>"`
	EmojiNum         int    `docopt:"-e"`
	PunctuationLevel int    `docopt:"-p"`
}

func main() {
	config := generator.Config{
		TargetName:       "くりきん",
		EmojiNum:         2,
		PunctuationLevel: 1,
	}

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			ojisanMessage, err := generator.Start(config)
			if err != nil {
				log.Fatal(err)
			}

			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ojisanMessage)).Do(); err != nil {
				log.Print(err)
			}
		}
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
