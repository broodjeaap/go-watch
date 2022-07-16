package main

import (
	"log"
	"net/http"
	"text/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Web struct {
	Bot *tgbotapi.BotAPI
}

func (web Web) index(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")

	tmpl, _ := template.New("name").Parse("{{.Message}}")
	context := struct {
		Message string
	}{
		Message: message,
	}
	msg := tgbotapi.NewMessage(abc, message)
	web.Bot.Send(msg)
	tmpl.Execute(w, context)
}

func passiveBot(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func main() {
	bot, _ := tgbotapi.NewBotAPI("AAAA")

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	go passiveBot(bot)

	web := Web{
		bot,
	}

	http.HandleFunc("/", web.index)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
