package main

import (
	"encoding/json"
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
	var watch Watch
	errMap, err := bindAndValidateWatch(&watch, c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500", errMap)
		return
	}
	web.db.Create(&watch)
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

type FilterDepth struct {
	Filter Filter
	Depth  int
}

func (web Web) viewWatch(c *gin.Context) {
	id := c.Param("id")

	var watch Watch
	web.db.Model(&Watch{}).First(&watch, id)

	var filters []Filter
	web.db.Model(&Filter{}).Find(&filters)

	queuedFilters := []*Filter{}
	filterMap := make(map[uint]*Filter)
	for _, filter := range filters {
		filterMap[filter.ID] = &filter
		if filter.ParentID == nil {
			queuedFilters = append(queuedFilters, &filter)
		}
		s, _ := json.MarshalIndent(filter, "", "\t")
		fmt.Println(s)
	}

	for _, filter := range filterMap {
		if filter.Parent != nil {
			parent := filterMap[*filter.ParentID]
			parent.Filters = append(parent.Filters, *filter)
		}
	}

	nextFilters := []*Filter{}
	bftFilters := []FilterDepth{}
	depth := 0
	for len(queuedFilters) > 0 {
		for _, f1 := range queuedFilters {
			bftFilters = append(bftFilters, FilterDepth{
				Filter: *f1,
				Depth:  depth,
			})
			for _, f2 := range f1.Filters {
				nextFilters = append(nextFilters, &f2)
			}
		}
		log.Println(nextFilters)
		queuedFilters = nextFilters
		log.Println(queuedFilters)
		nextFilters = []*Filter{}
		log.Println(nextFilters)
		depth += 1
	}

	c.HTML(http.StatusOK, "viewWatch", gin.H{
		"Watch":    watch,
		"Filters":  bftFilters,
		"MaxDepth": depth,
	})
}

func (web Web) createFilter(c *gin.Context) {
	var filter Filter
	errMap, err := bindAndValidateFilter(&filter, c)
	if err != nil {
		log.Print(err)
		c.HTML(http.StatusBadRequest, "500", errMap)
		return
	}
	web.db.Create(&filter)
	c.Redirect(http.StatusSeeOther, "/group/edit")
}

func (web Web) updateFilter(c *gin.Context) {
	var filterUpdate FilterUpdate
	errMap, err := bindAndValidateFilterUpdate(&filterUpdate, c)
	if err != nil {
		log.Print(err)
		c.HTML(http.StatusBadRequest, "500", errMap)
		return
	}
	var filter Filter
	web.db.First(&filter, filterUpdate.ID)
	filter.Name = filterUpdate.Name
	filter.Type = filterUpdate.Type
	filter.Var1 = filterUpdate.From
	filter.Var2 = &filterUpdate.To
	web.db.Save(&filter)
	c.Redirect(http.StatusSeeOther, "/group/edit/")
}

func (web Web) deleteFilter(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("filter_id"))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/watch/new")
		return
	}

	group_id := c.PostForm("group_id")
	web.db.Delete(&Filter{}, id)
	c.Redirect(http.StatusSeeOther, "/group/edit/"+group_id)
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
	db.AutoMigrate(&Watch{}, &Filter{})

	filters := []Filter{}
	watch := Watch{
		Name:     "LG C2 42",
		Interval: 60,
		Filters:  filters,
	}
	db.Create(&watch)

	urlFilter := Filter{
		WatchID:  watch.ID,
		ParentID: nil,
		Parent:   nil,
		Name:     "PriceWatch Fetch",
		Type:     "url",
		Var1:     "https://tweakers.net/pricewatch/1799060/lg-c2-42-inch-donkerzilveren-voet-zwart.html",
	}
	db.Create(&urlFilter)

	xpathFilter := Filter{
		WatchID:  watch.ID,
		Watch:    watch,
		ParentID: &urlFilter.ID,
		Name:     "price select",
		Type:     "xpath",
		Var1:     "//td[@class='shop-price']",
	}
	db.Create(&xpathFilter)

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
	templates.AddFromFiles("editGroup", "templates/base.html", "templates/editGroup.html")

	templates.AddFromFiles("500", "templates/base.html", "templates/500.html")
	router.HTMLRender = templates

	router.GET("/", web.index)
	router.GET("/watch/new", web.newWatch)
	router.POST("/watch/create", web.createWatch)
	router.POST("/watch/delete", web.deleteWatch)
	router.GET("/watch/view/:id/", web.viewWatch)
	router.POST("/filter/create/", web.createFilter)
	router.POST("/filter/update/", web.updateFilter)
	router.POST("/filter/delete/", web.deleteFilter)

	router.Run("0.0.0.0:8080")
}
