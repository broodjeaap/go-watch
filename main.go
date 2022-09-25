package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var baseHTML = filepath.Join("templates", "base.html")
var indexHTML = filepath.Join("templates", "index.html")
var newWatchHTML = filepath.Join("templates", "newWatch.html")

type Web struct {
	//Bot *tgbotapi.BotAPI
	db *gorm.DB
}

func (web Web) index(c *gin.Context) {
	//msg := tgbotapi.NewMessage(viper.GetInt64("telegram.chat"), message)
	//web.Bot.Send(msg)
	watches := []Watch{}
	web.db.Find(&watches)
	c.HTML(http.StatusOK, "index", watches)
}

func (web Web) watchCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "watchCreate", gin.H{})
}

func (web Web) watchCreatePost(c *gin.Context) {
	var watch Watch
	errMap, err := bindAndValidateWatch(&watch, c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500", errMap)
		return
	}
	web.db.Create(&watch)
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/%d", watch.ID))
}

func (web Web) deleteWatch(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("watch_id"))
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	web.db.Delete(&Watch{}, id)
	c.Redirect(http.StatusSeeOther, "/")
}

func (web Web) watchView(c *gin.Context) {
	id := c.Param("id")

	var watch Watch
	web.db.Model(&Watch{}).First(&watch, id)

	var filters []Filter
	web.db.Model(&Filter{}).Where("watch_id = ?", watch.ID).Find(&filters)

	var connections []FilterConnection
	web.db.Model(&FilterConnection{}).Where("watch_id = ?", watch.ID).Find(&connections)

	c.HTML(http.StatusOK, "watchView", gin.H{
		"Watch":       watch,
		"Filters":     filters,
		"Connections": connections,
	})
}

func (web Web) watchUpdate(c *gin.Context) {
	var watch Watch
	bindAndValidateWatch(&watch, c)

	web.db.Save(&watch)

	var newFilters []Filter
	var filtersJson = c.PostForm("filters")
	if err := json.Unmarshal([]byte(filtersJson), &newFilters); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var newConnections []FilterConnection
	var connectionsJson = c.PostForm("connections")
	if err := json.Unmarshal([]byte(connectionsJson), &newConnections); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var oldFilters []Filter
	web.db.Model(&Filter{}).Where("watch_id = ?", watch.ID).Find(&oldFilters)
	filterMap := make(map[uint]*Filter)
	for i := range newFilters {
		filter := &newFilters[i]
		filterMap[filter.ID] = filter
		filter.ID = 0
	}
	web.db.Delete(&Filter{}, "watch_id = ?", watch.ID)

	if len(newFilters) > 0 {
		web.db.Create(&newFilters)
	}

	web.db.Delete(&FilterConnection{}, "watch_id = ?", watch.ID)
	for i := range newConnections {
		connection := &newConnections[i]
		connection.OutputID = filterMap[connection.OutputID].ID
		connection.InputID = filterMap[connection.InputID].ID
	}
	if len(newConnections) > 0 {
		web.db.Create(&newConnections)
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/%d", watch.ID))
}

/*
func (web Web) viewWatch(c *gin.Context) {
	id := c.Param("id")

	var watch Watch
	web.db.Model(&Watch{}).First(&watch, id)

	var filters []Filter
	web.db.Model(&Filter{}).Where("watch_id = ?", watch.ID).Find(&filters)

	filterMap := make(map[uint]*Filter)
	for i := range filters {
		filterMap[filters[i].ID] = &filters[i]
	}

	currentLayerFilters := []*Filter{}
	for i := range filters {
		filter := &filters[i]
		if filter.ParentID != nil {
			parent := filterMap[*filter.ParentID]
			parent.Filters = append(parent.Filters, *filter)
		} else {
			currentLayerFilters = append(currentLayerFilters, filter)
		}
	}

	bftFilters := []FilterDepth{}
	nextLayerFilters := []*Filter{}
	depth := 0
	for len(nextLayerFilters) > 0 || len(currentLayerFilters) > 0 {
		for len(currentLayerFilters) > 0 {
			filter := currentLayerFilters[0]
			bftFilters = append(bftFilters, FilterDepth{
				Filter: filter,
				Depth:  make([]struct{}, depth),
			})
			for _, filter := range filter.Filters {
				nextLayerFilters = append(nextLayerFilters, &filter)
			}
			currentLayerFilters = currentLayerFilters[1:]
		}
		depth += 1
		currentLayerFilters = nextLayerFilters
		nextLayerFilters = []*Filter{}
	}

	for i := range bftFilters {
		fd := &bftFilters[i]
		fd.RevDepth = make([]struct{}, depth-len(fd.Depth))
	}
	numberOfColumns := depth + 4
	c.HTML(http.StatusOK, "viewWatch", gin.H{
		"Watch":       watch,
		"Filters":     bftFilters,
		"MaxDepth":    depth,
		"Columns":     make([]struct{}, numberOfColumns),
		"ColumnWidth": 100 / numberOfColumns,
	})
}
*/

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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Could not load config file")
	}

	db, _ := gorm.Open(sqlite.Open(viper.GetString("database.dsn")))
	db.AutoMigrate(&Watch{}, &Filter{}, &FilterConnection{})

	//bot, _ := tgbotapi.NewBotAPI(viper.GetString("telegram.token"))

	//bot.Debug = true

	//log.Printf("Authorized on account %s", bot.Self.UserName)

	//go passiveBot(bot)

	web := Web{
		//bot,
		db,
	}
	router := gin.Default()

	router.Static("/static", "./static")

	templates := multitemplate.NewRenderer()
	templates.AddFromFiles("index", "templates/base.html", "templates/index.html")
	templates.AddFromFiles("watchCreate", "templates/base.html", "templates/watch/create.html")
	templates.AddFromFiles("watchView", "templates/base.html", "templates/watch/view.html")

	templates.AddFromFiles("500", "templates/base.html", "templates/500.html")
	router.HTMLRender = templates

	router.GET("/", web.index)

	router.GET("/watch/:id", web.watchView)
	router.GET("/watch/new", web.watchCreate)
	router.POST("/watch/create", web.watchCreatePost)
	router.POST("/watch/update", web.watchUpdate)
	router.POST("/watch/delete", web.deleteWatch)

	router.Run("0.0.0.0:8080")
}
