package notifiers

import (
	"fmt"
	"log"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/spf13/viper"
)

type ShoutrrrNotifier struct {
	URLs []string
}

func (shoutr *ShoutrrrNotifier) Open(configPath string) bool {
	urlsPath := fmt.Sprintf("%s.urls", configPath)
	if !viper.IsSet(urlsPath) {
		log.Println("Need 'urls' for Shoutrrr")
		return false
	}
	shoutr.URLs = viper.GetStringSlice(urlsPath)
	log.Println("Shoutrrr version:", shoutrrr.Version())
	return true
}

func (shoutr *ShoutrrrNotifier) Message(message string) bool {
	sender, err := shoutrrr.CreateSender(shoutr.URLs...)
	if err != nil {
		log.Println("Could not create Shoutrrr sender:", err)
		return false
	}

	errs := sender.Send(message, &types.Params{})
	if errs != nil {
		for _, err := range errs {
			if err != nil {
				log.Println("Something went wrong sending Shoutrrr:", err)
			}
		}
	}
	return true
}

func (shoutr *ShoutrrrNotifier) Close() bool {
	return true
}
