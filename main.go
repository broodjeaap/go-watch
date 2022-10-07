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
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var baseHTML = filepath.Join("templates", "base.html")
var indexHTML = filepath.Join("templates", "index.html")
var newWatchHTML = filepath.Join("templates", "newWatch.html")

type Web struct {
	//Bot *tgbotapi.BotAPI
	router    *gin.Engine
	templates multitemplate.Renderer
	cron      *cron.Cron
	urlCache  map[string]string
	cronWatch map[uint]cron.Entry
	db        *gorm.DB
}

func newWeb() *Web {
	web := &Web{
		urlCache: make(map[string]string, 5),
	}
	web.init()
	return web
}

func (web *Web) init() {
	web.urlCache = make(map[string]string, 5)
	web.initDB()

	web.initRouter()
	web.initCronJobs()
}

func (web *Web) initDB() {
	db, err := gorm.Open(sqlite.Open(viper.GetString("database.dsn")))
	if err != nil {
		log.Panicln("Could not start DB: ", err)
	}
	web.db = db
	web.db.AutoMigrate(&Watch{}, &Filter{}, &FilterConnection{}, &FilterOutput{})
}
func (web *Web) initRouter() {
	web.router = gin.Default()

	web.router.Static("/static", "./static")

	web.initTemplates()
	web.router.HTMLRender = web.templates

	web.router.GET("/", web.index)

	web.router.GET("/watch/view/:id", web.watchView)
	web.router.GET("/watch/edit/:id", web.watchEdit)
	web.router.GET("/watch/new", web.watchCreate)
	web.router.POST("/watch/create", web.watchCreatePost)
	web.router.POST("/watch/update", web.watchUpdate)
	web.router.POST("/watch/delete", web.deleteWatch)

	web.router.GET("/cache/view", web.cacheView)
	web.router.POST("/cache/clear", web.cacheClear)
}

func (web *Web) initTemplates() {
	web.templates = multitemplate.NewRenderer()
	web.templates.AddFromFiles("index", "templates/base.html", "templates/index.html")
	web.templates.AddFromFiles("watchCreate", "templates/base.html", "templates/watch/create.html")
	web.templates.AddFromFiles("watchView", "templates/base.html", "templates/watch/view.html")
	web.templates.AddFromFiles("watchEdit", "templates/base.html", "templates/watch/edit.html")

	web.templates.AddFromFiles("cacheView", "templates/base.html", "templates/cache/view.html")

	web.templates.AddFromFiles("500", "templates/base.html", "templates/500.html")
}

func (web *Web) initCronJobs() {
	var cronFilters []Filter
	web.db.Model(&Filter{}).Find(&cronFilters, "type = 'cron'")
	web.cronWatch = make(map[uint]cron.Entry, len(cronFilters))
	web.cron = cron.New()
	for _, cronFilter := range cronFilters {
		entryID, err := web.cron.AddFunc(cronFilter.Var1, func() { triggerSchedule(cronFilter.WatchID, web) })
		if err != nil {
			log.Println("Could not start job for Watch: ", cronFilter.WatchID)
			continue
		}
		log.Println("Started CronJob for WatchID", cronFilter.WatchID, "with schedule:", cronFilter.Var1)
		web.cronWatch[cronFilter.WatchID] = web.cron.Entry(entryID)
	}
	web.cron.Start()
}

func (web *Web) run() {
	web.router.Run("0.0.0.0:8080")
}

type WatchEntry struct {
	Watch *Watch
	Entry *cron.Entry
}

func (web *Web) index(c *gin.Context) {
	//msg := tgbotapi.NewMessage(viper.GetInt64("telegram.chat"), message)
	//web.Bot.Send(msg)
	watches := []Watch{}
	web.db.Find(&watches)

	for i := 0; i < len(watches); i++ {
		entry := web.cronWatch[watches[i].ID]
		watches[i].CronEntry = &entry
	}

	c.HTML(http.StatusOK, "index", watches)
}

func (web *Web) watchCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "watchCreate", gin.H{})
}

func (web *Web) watchCreatePost(c *gin.Context) {
	var watch Watch
	errMap, err := bindAndValidateWatch(&watch, c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500", errMap)
		return
	}
	web.db.Create(&watch)
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/%d", watch.ID))
}

func (web *Web) deleteWatch(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("watch_id"))
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	web.db.Delete(&FilterConnection{}, "watch_id = ?", id)
	web.db.Delete(&FilterOutput{}, "watch_id = ?", id)
	web.db.Delete(&Filter{}, "watch_id = ?", id)

	web.db.Delete(&Watch{}, id)
	c.Redirect(http.StatusSeeOther, "/")
}

func (web *Web) watchView(c *gin.Context) {
	id := c.Param("id")

	var watch Watch
	web.db.Model(&Watch{}).First(&watch, id)
	entry, exists := web.cronWatch[watch.ID]
	if !exists {
		log.Println("Could not find entry for Watch", watch.ID)
		c.HTML(http.StatusNotFound, "watchView", gin.H{"error": "Entry not found"})
		return
	}
	watch.CronEntry = &entry

	var values []FilterOutput
	web.db.Model(&FilterOutput{}).Where("watch_id = ?", watch.ID).Find(&values)

	valueMap := make(map[string][]FilterOutput, len(values))
	for _, value := range values {
		valueMap[value.Name] = append(valueMap[value.Name], value)
	}

	c.HTML(http.StatusOK, "watchView", gin.H{
		"Watch":    watch,
		"ValueMap": valueMap,
	})
}

func (web *Web) watchEdit(c *gin.Context) {
	id := c.Param("id")

	var watch Watch
	web.db.Model(&Watch{}).First(&watch, id)

	var filters []Filter
	web.db.Model(&Filter{}).Where("watch_id = ?", watch.ID).Find(&filters)

	var connections []FilterConnection
	web.db.Model(&FilterConnection{}).Where("watch_id = ?", watch.ID).Find(&connections)

	var values []FilterOutput
	web.db.Model(&FilterOutput{}).Where("watch_id = ?", watch.ID).Find(&values)

	buildFilterTree(filters, connections)
	processFilters(filters, web, &watch, true, true)

	c.HTML(http.StatusOK, "watchEdit", gin.H{
		"Watch":       watch,
		"Filters":     filters,
		"Connections": connections,
		"Values":      values,
	})
}

func (web *Web) watchUpdate(c *gin.Context) {
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

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/edit/%d", watch.ID))
}

func (web *Web) cacheView(c *gin.Context) {
	c.HTML(http.StatusOK, "cacheView", web.urlCache)
}

func (web *Web) cacheClear(c *gin.Context) {
	url := c.PostForm("url")
	delete(web.urlCache, url)
	c.Redirect(http.StatusSeeOther, "/cache/view")
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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Could not load config file")
	}

	//bot, _ := tgbotapi.NewBotAPI(viper.GetString("telegram.token"))

	//bot.Debug = true

	//log.Printf("Authorized on account %s", bot.Self.UserName)

	//go passiveBot(bot)

	web := newWeb()
	web.run()
}
