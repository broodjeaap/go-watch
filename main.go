package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

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

func (web Web) newWatch(c *gin.Context) {
	c.HTML(http.StatusOK, "newWatch", gin.H{})
}

func (web Web) createWatch(c *gin.Context) {
	interval, err := strconv.Atoi(c.PostForm("interval"))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}

	watch := &Watch{
		Name:     c.PostForm("name"),
		Interval: interval,
	}
	web.db.Create(watch)
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/view/%d", watch.ID))
}

func (web Web) deleteWatch(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}

	web.db.Delete(&Watch{}, id)
	c.Redirect(http.StatusSeeOther, "/")
}

func (web Web) viewWatch(c *gin.Context) {
	id := c.Param("id")

	var watch Watch
	web.db.Model(&Watch{}).Preload("URLs.Queries.Filters").First(&watch, id)
	c.HTML(http.StatusOK, "viewWatch", watch)
}

func (web Web) createURL(c *gin.Context) {
	watch_id, err := strconv.ParseUint(c.PostForm("watch_id"), 10, 64)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return // TODO response
	}
	name := c.PostForm("name")
	if name == "" {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}
	url := c.PostForm("url")
	if url == "" {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}

	url_model := &URL{
		WatchID: uint(watch_id),
		Name:    name,
		URL:     url,
	}
	web.db.Create(url_model)
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/view/%d", watch_id))
}

func (web Web) createQuery(c *gin.Context) {
	watch_id, err := strconv.ParseUint(c.PostForm("watch_id"), 10, 64)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return // TODO response
	}
	url_id, err := strconv.ParseUint(c.PostForm("url_id"), 10, 64)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return // TODO response
	}
	name := c.PostForm("name")
	if name == "" {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}
	typ := c.PostForm("type")
	if typ == "" {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}
	query := c.PostForm("query")
	if query == "" {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}

	query_model := &Query{
		URLID: uint(url_id),
		Name:  name,
		Type:  typ,
		Query: query,
	}
	web.db.Create(query_model)
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/view/%d", watch_id))
}

func (web Web) createFilter(c *gin.Context) {
	query_id, err := strconv.ParseUint(c.PostForm("query_id"), 10, 64)
	if err != nil {
		log.Print(err)
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return // TODO response
	}
	name := c.PostForm("name")
	if name == "" {
		log.Print(name)
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}
	typ := c.PostForm("type")
	if typ == "" {
		log.Print(typ)
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}
	from := c.PostForm("from")
	if from == "" {
		log.Print(from)
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}
	to := c.PostForm("to")
	log.Print("To:", to)
	filter_model := &Filter{
		QueryID: uint(query_id),
		Name:    name,
		Type:    typ,
		From:    from,
		To:      to,
	}
	web.db.Create(filter_model)
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/query/edit/%d", query_id))
}

func (web Web) editQuery(c *gin.Context) {
	query_id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return // TODO response
	}
	var query Query
	web.db.Preload("URL.Watch").Preload("Filters").Preload("URL").First(&query, query_id)

	c.HTML(http.StatusOK, "editQuery", gin.H{
		"Query":         query,
		"currentResult": getQueryResult(&query),
	})
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

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename, baseHTML)
	if err != nil {
		log.Print(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
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
	db.AutoMigrate(&Watch{}, &URL{}, &Query{}, &Filter{})

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
	templates.AddFromFiles("newWatch", "templates/base.html", "templates/newWatch.html")
	templates.AddFromFiles("viewWatch", "templates/base.html", "templates/viewWatch.html")
	templates.AddFromFiles("editQuery", "templates/base.html", "templates/editQuery.html")
	router.HTMLRender = templates

	router.GET("/", web.index)
	router.GET("/watch/new", web.newWatch)
	router.POST("/watch/create", web.createWatch)
	router.POST("/watch/delete", web.deleteWatch)
	router.GET("/watch/view/:id/", web.viewWatch)
	router.POST("/url/create/", web.createURL)
	router.POST("/query/create/", web.createQuery)
	router.GET("/query/edit/:id", web.editQuery)
	router.POST("/filter/create/", web.createFilter)

	router.Run("0.0.0.0:8080")
}
