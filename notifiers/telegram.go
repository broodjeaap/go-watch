package notifiers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

type TelegramNotifier struct {
	Bot   *tgbotapi.BotAPI
	Token string
	Debug bool
}

func (telegram *TelegramNotifier) Open() bool {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("notifiers.telegram.token"))
	if err != nil {
		log.Println("Could not start Telegram notifier:\n", err)
		return false
	}
	telegram.Bot = bot
	bot.Debug = viper.GetBool("notifiers.telegram.debug")
	log.Printf("Authorized telegram bot: %s", bot.Self.UserName)
	return true
}

func (telegram *TelegramNotifier) Message(message string) bool {
	msg := tgbotapi.NewMessage(viper.GetInt64("notifiers.telegram.chat"), message)
	_, err := telegram.Bot.Send(msg)
	if err != nil {
		log.Println("Could not send Telegram message:\n", err)
		return false
	}
	return true
}

func (telegram *TelegramNotifier) Close() bool {
	return true
}
