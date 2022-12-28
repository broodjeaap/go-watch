package notifiers

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	Server   string
	Port     int
	From     string
	User     string
	To       string
	Password string
	Debug    bool
}

func (email *EmailNotifier) Open(configPath string) bool {
	serverPath := fmt.Sprintf("%s.server", configPath)
	if viper.IsSet(serverPath) {
		email.Server = viper.GetString(serverPath)
	} else {
		log.Println("Need 'server' var for email notifier")
		return false
	}
	portPath := fmt.Sprintf("%s.port", configPath)
	if viper.IsSet(portPath) {
		email.Port = viper.GetInt(portPath)
	} else {
		log.Println("Need 'port' var for email notifier")
		return false
	}
	userPath := fmt.Sprintf("%s.user", configPath)
	if viper.IsSet(userPath) {
		email.User = viper.GetString(userPath)
	} else {
		log.Println("Need 'user' var for email notifier")
		return false
	}
	fromPath := fmt.Sprintf("%s.from", configPath)
	if viper.IsSet(fromPath) {
		email.From = viper.GetString(fromPath)
	} else {
		log.Println("Need 'from' var for email notifier")
		return false
	}
	toPath := fmt.Sprintf("%s.to", configPath)
	if viper.IsSet(toPath) {
		email.To = viper.GetString(toPath)
	} else {
		log.Println("Need 'to' var for email notifier")
		return false
	}
	passwordPath := fmt.Sprintf("%s.password", configPath)
	if viper.IsSet(passwordPath) {
		email.Password = viper.GetString(passwordPath)
	} else {
		log.Println("Need 'password' var for email notifier")
		return false
	}
	log.Printf("Configured email: %s, %s", email.From, email.To)
	return true
}

func (email *EmailNotifier) Message(message string) bool {
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {email.From},
		"To":      {email.To},
		"Subject": {"GoWatch"},
	})
	m.SetBody("text/html", message)

	d := gomail.NewDialer(email.Server, email.Port, email.User, email.Password)

	err := d.DialAndSend(m)
	if err != nil {
		log.Println("Could not send email:", err)
		return false
	}
	return true
}

func (email *EmailNotifier) Close() bool {
	return true
}
