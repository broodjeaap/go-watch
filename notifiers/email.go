package notifiers

import (
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

func (email *EmailNotifier) Open() bool {
	if viper.IsSet("notifiers.email.server") {
		email.Server = viper.GetString("notifiers.email.server")
	} else {
		log.Println("Need 'server' var for email notifier")
		return false
	}
	if viper.IsSet("notifiers.email.port") {
		email.Port = viper.GetInt("notifiers.email.port")
	} else {
		log.Println("Need 'port' var for email notifier")
		return false
	}
	if viper.IsSet("notifiers.email.user") {
		email.User = viper.GetString("notifiers.email.user")
	} else {
		log.Println("Need 'user' var for email notifier")
		return false
	}
	if viper.IsSet("notifiers.email.from") {
		email.From = viper.GetString("notifiers.email.from")
	} else {
		log.Println("Need 'from' var for email notifier")
		return false
	}
	if viper.IsSet("notifiers.email.to") {
		email.To = viper.GetString("notifiers.email.to")
	} else {
		log.Println("Need 'to' var for email notifier")
		return false
	}
	if viper.IsSet("notifiers.email.password") {
		email.Password = viper.GetString("notifiers.email.password")
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
