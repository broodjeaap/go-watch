package notifiers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

type TelegramNotifier struct {
	Bot    *tgbotapi.BotAPI
	Token  string
	ChatID int64
	Debug  bool
}

func (telegram *TelegramNotifier) Open(configPath string) bool {
	tokenPath := fmt.Sprintf("%s.token", configPath)
	if !viper.IsSet(tokenPath) {
		log.Println("Telegram needs 'token' value set")
		return false
	}
	telegram.Token = viper.GetString(tokenPath)
	bot, err := tgbotapi.NewBotAPI(telegram.Token)
	if err != nil {
		log.Println("Could not start Telegram notifier:\n", err)
		return false
	}

	chatIDPath := fmt.Sprintf("%s.chat", configPath)
	if !viper.IsSet(chatIDPath) {
		log.Panicln("Telegram needs 'chat' ID value")
		return false
	}
	telegram.ChatID = viper.GetInt64(chatIDPath)
	telegram.Bot = bot
	bot.Debug = viper.GetBool("notifiers.telegram.debug")
	log.Printf("Authorized telegram bot: %s", bot.Self.UserName)
	return true
}

func (telegram *TelegramNotifier) Message(message string) bool {
	msg := tgbotapi.NewMessage(telegram.ChatID, message)
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
