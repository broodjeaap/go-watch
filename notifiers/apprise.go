package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type AppriseNotifier struct {
	URL    string
	Title  string
	Type   string
	Format string
	URLs   []string
}

func (apprise *AppriseNotifier) Open(configPath string) bool {
	urlPath := fmt.Sprintf("%s.url", configPath)
	if !viper.IsSet(urlPath) {
		log.Println("Need 'url' for Apprise")
		return false
	}
	apprise.URL = viper.GetString(urlPath)

	urlsPath := fmt.Sprintf("%s.urls", configPath)
	if !viper.IsSet(urlsPath) {
		log.Println("Need 'urls' for Apprise")
		return false
	}

	apprise.Title = "GoWatch Notification"
	titlePath := fmt.Sprintf("%s.title", configPath)
	if viper.IsSet(titlePath) {
		apprise.Title = viper.GetString(titlePath)
	}
	apprise.Type = "info"
	typePath := fmt.Sprintf("%s.mtype", configPath)
	if viper.IsSet(typePath) {
		apprise.Type = viper.GetString(typePath)
	}
	apprise.Format = "text"
	formatPath := fmt.Sprintf("%s.format", configPath)
	if viper.IsSet(formatPath) {
		apprise.Format = viper.GetString(formatPath)
	}
	apprise.URLs = viper.GetStringSlice(urlsPath)
	log.Println("Apprise notifier:", apprise.URL, apprise.Type, apprise.Format)
	return true
}

type ApprisePostData struct {
	Title  string   `json:"title"`
	Type   string   `json:"type"`
	Format string   `json:"format"`
	URLs   []string `json:"urls"`
	Body   string   `json:"body"`
}

func (apprise *AppriseNotifier) Message(message string) bool {
	data := ApprisePostData{
		URLs:   apprise.URLs,
		Title:  apprise.Title,
		Type:   apprise.Type,
		Format: apprise.Format,
		Body:   message,
	}
	jsn, err := json.Marshal(data)
	if err != nil {
		log.Panicln("Could not create JSON post data:", err)
		return false
	}

	resp, err := http.Post(apprise.URL, "application/json", bytes.NewBuffer(jsn))
	if err != nil {
		log.Println("Could not send Apprise notification:", err)
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not parse Apprise response:", err)
		return false
	}
	log.Println(string(body))
	return true
}

func (apprise *AppriseNotifier) Close() bool {
	return true
}
