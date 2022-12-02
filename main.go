package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"broodjeaap.net/go-watch/notifiers"
)

var baseHTML = filepath.Join("templates", "base.html")
var indexHTML = filepath.Join("templates", "index.html")
var newWatchHTML = filepath.Join("templates", "newWatch.html")

type Web struct {
	router    *gin.Engine
	templates multitemplate.Renderer
	cron      *cron.Cron
	urlCache  map[string]string
	cronWatch map[uint]cron.Entry
	db        *gorm.DB
	notifiers map[string]notifiers.Notifier
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
	web.initNotifiers()
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
	web.router.GET("/watch/export/:id", web.exportWatch)
	web.router.POST("/watch/import/:id", web.importWatch)

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
		web.cronWatch[cronFilter.ID] = web.cron.Entry(entryID)
	}
	web.cron.Start()
}

func (web *Web) initNotifiers() {
	web.notifiers = make(map[string]notifiers.Notifier, 5)
	if viper.IsSet("telegram") {
		telegramBot := notifiers.TelegramNotifier{}
		telegramBot.Open()
		web.notifiers["Telegram"] = telegramBot

	}
}

func (web *Web) notify(notifierKey string, message string) {
	notifier, exists := web.notifiers[notifierKey]
	if !exists {
		log.Println("Could not find notifier with key:", notifierKey)
	}
	notifier.Message(message)
}

func (web *Web) run() {
	web.router.Run("0.0.0.0:8080")
}

type WatchEntry struct {
	Watch *Watch
	Entry *cron.Entry
}

func (web *Web) index(c *gin.Context) {
	watches := []Watch{}
	web.db.Find(&watches)

	watchMap := make(map[uint]*Watch, len(watches))
	for i := 0; i < len(watches); i++ {
		watchMap[watches[i].ID] = &watches[i]
	}
	// this doesn't work with multiple schedule filters per watch, but meh
	var filters []Filter
	web.db.Model(&Filter{}).Find(&filters, "type = 'cron'")
	for _, filter := range filters {
		entry, exists := web.cronWatch[filter.ID]
		if !exists {
			log.Println("No cron entry for filter", filter.ID, filter.Name)
			continue
		}
		watchMap[filter.WatchID].CronEntry = &entry
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
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/edit/%d", watch.ID))
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

	var cronFilters []Filter
	web.db.Model(&Filter{}).Find(&cronFilters, "watch_id = ? AND type = 'cron'", id)
	for _, filter := range cronFilters {
		entry, exist := web.cronWatch[filter.ID]
		if exist {
			web.cron.Remove(entry.ID)
			delete(web.cronWatch, filter.ID)
		}
	}
	web.db.Delete(&Filter{}, "watch_id = ?", id)

	web.db.Delete(&Watch{}, id)
	c.Redirect(http.StatusSeeOther, "/")
}

func (web *Web) watchView(c *gin.Context) {
	id := c.Param("id")

	var watch Watch
	web.db.Model(&Watch{}).First(&watch, id)

	var cronFilters []Filter
	web.db.Model(&Filter{}).Find(&cronFilters, "watch_id = ? AND type = 'cron'", id)
	for _, filter := range cronFilters {
		entry, exists := web.cronWatch[filter.ID]
		if !exists {
			log.Println("Could not find entry for filter", filter.ID, filter.Name)
			continue
		}
		watch.CronEntry = &entry
	}

	var values []FilterOutput
	web.db.Model(&FilterOutput{}).Order("time asc").Where("watch_id = ?", watch.ID).Find(&values)

	valueMap := make(map[string][]FilterOutput, len(values))
	names := make(map[string]bool, 5)
	for _, value := range values {
		names[value.Name] = true
		valueMap[value.Name] = append(valueMap[value.Name], value)
	}

	colorMap := make(map[string]int, len(names))
	index := 0
	for name, _ := range names {
		colorMap[name] = index % 16 // only 16 colors
		index += 1
	}

	//data := make([]map[string]string, len(valueMap))

	c.HTML(http.StatusOK, "watchView", gin.H{
		"Watch":    watch,
		"ValueMap": valueMap,
		"colorMap": colorMap,
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
	processFilters(filters, web, &watch, true)

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

	// stop/delete cronjobs running for this watch
	var cronFilters []Filter
	web.db.Model(&Filter{}).Where("watch_id = ? AND type = 'cron'", watch.ID).Find(&cronFilters)
	for _, filter := range cronFilters {
		entry, exist := web.cronWatch[filter.ID]
		if exist {
			web.cron.Remove(entry.ID)
			delete(web.cronWatch, filter.ID)
		} else {
			log.Println("Tried removing cron entry but ID not found ", filter.ID)
			log.Println(web.cronWatch)
		}
	}

	web.db.Delete(&Filter{}, "watch_id = ?", watch.ID)

	filterMap := make(map[uint]*Filter)
	if len(newFilters) > 0 {
		for i := range newFilters {
			filter := &newFilters[i]
			filterMap[filter.ID] = filter
			filter.ID = 0
		}

		web.db.Create(&newFilters)

		for i := range newFilters {
			filter := &newFilters[i]
			if filter.Type != "cron" {
				continue
			}
			entryID, err := web.cron.AddFunc(filter.Var1, func() { triggerSchedule(filter.WatchID, web) })
			if err != nil {
				log.Println("Could not start job for Watch: ", filter.WatchID)
				continue
			}
			log.Println("Started CronJob for WatchID", filter.WatchID, "FilterID", filter.ID, "with schedule:", filter.Var1)
			web.cronWatch[filter.ID] = web.cron.Entry(entryID)
		}
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

func (web *Web) exportWatch(c *gin.Context) {
	watchID := c.Param("id")
	export := WatchExport{}
	var watch Watch
	web.db.Model(&Watch{}).Find(&watch, watchID)
	web.db.Model(&Filter{}).Find(&export.Filters, "watch_id = ?", watchID)
	web.db.Model(&FilterConnection{}).Find(&export.Connections, "watch_id = ?", watchID)

	c.Header("Content-Disposition", "attachment; filename="+watch.Name+".json")
	c.Header("Content-Type", c.Request.Header.Get("Content-Type"))

	c.JSON(http.StatusOK, export)
}

func (web *Web) importWatch(c *gin.Context) {
	watchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	file, err := c.FormFile("json")

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	jsn, _ := ioutil.ReadAll(openedFile)

	export := WatchExport{}

	if err := json.Unmarshal(jsn, &export); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// stop/delete cronjobs running for this watch
	var cronFilters []Filter
	web.db.Model(&Filter{}).Where("watch_id = ? AND type = 'cron'", watchID).Find(&cronFilters)
	for _, filter := range cronFilters {
		entry, exist := web.cronWatch[filter.ID]
		if exist {
			web.cron.Remove(entry.ID)
			delete(web.cronWatch, filter.ID)
		}
	}

	filterMap := make(map[uint]*Filter)
	for i := range export.Filters {
		filter := &export.Filters[i]
		filterMap[filter.ID] = filter
		filter.ID = 0
		filter.WatchID = uint(watchID)
	}
	web.db.Delete(&Filter{}, "watch_id = ?", watchID)

	if len(export.Filters) > 0 {
		tx := web.db.Create(&export.Filters)
		if tx.Error != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		for _, filter := range export.Filters {
			if filter.Type == "cron" {
				entryID, err := web.cron.AddFunc(filter.Var1, func() { triggerSchedule(filter.WatchID, web) })
				if err != nil {
					log.Println("Could not start job for Watch: ", filter.WatchID)
					continue
				}
				log.Println("Started CronJob for WatchID", filter.WatchID, "with schedule:", filter.Var1)
				web.cronWatch[filter.ID] = web.cron.Entry(entryID)
			}
		}
	}

	web.db.Delete(&FilterConnection{}, "watch_id = ?", watchID)
	for i := range export.Connections {
		connection := &export.Connections[i]
		connection.ID = 0
		connection.WatchID = uint(watchID)
		connection.OutputID = filterMap[connection.OutputID].ID
		connection.InputID = filterMap[connection.InputID].ID
	}
	if len(export.Connections) > 0 {
		tx := web.db.Create(&export.Connections)
		if tx.Error != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/edit/%d", watchID))
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

	web := newWeb()
	web.run()
}
