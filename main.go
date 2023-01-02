package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"broodjeaap.net/go-watch/notifiers"
)

//go:embed templates static watchTemplates config.tmpl
var EMBED_FS embed.FS

type Web struct {
	router          *gin.Engine
	templates       multitemplate.Renderer
	cron            *cron.Cron
	urlCache        map[string]string
	cronWatch       map[uint]cron.EntryID
	db              *gorm.DB
	notifiers       map[string]notifiers.Notifier
	startupWarnings []string
}

func newWeb() *Web {
	web := &Web{
		urlCache:        make(map[string]string, 5),
		startupWarnings: make([]string, 0, 10),
	}
	web.init()
	return web
}

func (web *Web) init() {
	web.urlCache = make(map[string]string, 5)
	if !viper.GetBool("gin.debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	web.validateProxyURL()
	web.initDB()
	web.initRouter()
	web.initCronJobs()
	web.initNotifiers()
}

func (web *Web) startupWarning(m ...any) {
	warning := fmt.Sprint(m...)
	log.Println(warning)
	web.startupWarnings = append(web.startupWarnings, warning)
}

func (web *Web) validateProxyURL() {
	if viper.IsSet("proxy.proxy_url") {
		_, err := url.Parse(viper.GetString("proxy.proxy_url"))
		if err != nil {
			web.startupWarning("Could not parse proxy url, check config")
			return
		}
	}
}

func (web *Web) initDB() {
	dsn := "./watch.db"
	if viper.IsSet("database.dsn") {
		dsn = viper.GetString("database.dsn")
	}
	var db *gorm.DB
	var err error
	if strings.HasPrefix(dsn, "sqlserver") {
		db, err = gorm.Open(sqlserver.Open(dsn))
		log.Println("Using SQLServer database")
	} else if strings.HasPrefix(dsn, "postgres") {
		db, err = gorm.Open(postgres.Open(dsn))
		log.Println("Using PostgreSQL database")
	} else if strings.HasPrefix(dsn, "mysql") {
		db, err = gorm.Open(mysql.Open(dsn))
		log.Println("Using MySQL database")
	} else {
		db, err = gorm.Open(sqlite.Open(dsn))
		log.Println("Using sqlite database at:", dsn)
	}
	if db == nil {
		log.Panicln("Could not recognize database.dsn: ", dsn)
	}
	if err != nil {
		log.Panicln("Could not start DB: ", err)
	}
	web.db = db
	web.db.AutoMigrate(&Watch{}, &Filter{}, &FilterConnection{}, &FilterOutput{})
}
func (web *Web) initRouter() {
	web.router = gin.Default()

	staticFS, err := fs.Sub(EMBED_FS, "static")
	if err != nil {
		log.Fatalln("Could not load static embed fs")
	}
	web.router.StaticFS("/static", http.FS(staticFS))

	web.initTemplates()
	web.router.HTMLRender = web.templates

	web.router.GET("/", web.index)

	web.router.GET("/watch/view/:id", web.watchView)
	web.router.GET("/watch/edit/:id", web.watchEdit)
	web.router.GET("/watch/create", web.watchCreate)
	web.router.POST("/watch/create", web.watchCreatePost)
	web.router.POST("/watch/update", web.watchUpdate)
	web.router.POST("/watch/delete", web.deleteWatch)
	web.router.GET("/watch/export/:id", web.exportWatch)
	web.router.POST("/watch/import/:id", web.importWatch)

	web.router.GET("/cache/view", web.cacheView)
	web.router.POST("/cache/clear", web.cacheClear)

	web.router.GET("/notifiers/view", web.notifiersView)
	web.router.POST("/notifiers/test", web.notifiersTest)

	web.router.SetTrustedProxies(nil)
}

func (web *Web) initTemplates() {
	web.templates = multitemplate.NewRenderer()

	templatesFS, err := fs.Sub(EMBED_FS, "templates")
	if err != nil {
		log.Fatalln("Could not load templates embed fs")
	}

	web.templates.Add("index", template.Must(template.ParseFS(templatesFS, "base.html", "index.html")))

	web.templates.Add("watchView", template.Must(template.ParseFS(templatesFS, "base.html", "watch/view.html")))
	web.templates.Add("watchCreate", template.Must(template.ParseFS(templatesFS, "base.html", "watch/create.html")))
	web.templates.Add("watchEdit", template.Must(template.ParseFS(templatesFS, "base.html", "watch/edit.html")))

	web.templates.Add("cacheView", template.Must(template.ParseFS(templatesFS, "base.html", "cache/view.html")))

	web.templates.Add("notifiersView", template.Must(template.ParseFS(templatesFS, "base.html", "notifiers.html")))

	web.templates.Add("500", template.Must(template.ParseFS(templatesFS, "base.html", "500.html")))
}

func (web *Web) initCronJobs() {
	var cronFilters []Filter
	web.db.Model(&Filter{}).Find(&cronFilters, "type = 'cron' AND var2 = 'yes'")
	web.cronWatch = make(map[uint]cron.EntryID, len(cronFilters))
	web.cron = cron.New()
	for i := range cronFilters {
		cronFilter := &cronFilters[i]
		entryID, err := web.cron.AddFunc(cronFilter.Var1, func() { triggerSchedule(cronFilter.WatchID, web, &cronFilter.ID) })
		if err != nil {
			web.startupWarning("Could not start job for Watch: ", cronFilter.WatchID)
			continue
		}
		log.Println("Started CronJob for WatchID", cronFilter.WatchID, "with schedule:", cronFilter.Var1)
		web.cronWatch[cronFilter.ID] = entryID
	}
	if viper.IsSet("database.prune") {
		pruneSchedule := viper.GetString("database.prune")
		_, err := web.cron.AddFunc(pruneSchedule, web.pruneDB)
		if err != nil {
			web.startupWarning("Could not parse database.prune:", err)
		}
		log.Println("Started DB prune cronjob:", pruneSchedule)
	}
	web.cron.Start()
}

func (web *Web) initNotifiers() {
	web.notifiers = make(map[string]notifiers.Notifier, 5)
	if !viper.IsSet("notifiers") {
		web.startupWarning("No notifiers set!")
		return
	}
	notifiersMap := viper.GetStringMap("notifiers")
	for name := range notifiersMap {
		notifierPath := fmt.Sprintf("notifiers.%s", name)
		notifierMap := viper.GetStringMapString(notifierPath)

		notifierType, exists := notifierMap["type"]
		if !exists {
			web.startupWarning(fmt.Sprintf("No 'type' for '%s' notifier!", name))
			continue
		}
		success := false
		var notifier notifiers.Notifier
		switch notifierType {
		case "telegram":
			{
				notifier = &notifiers.TelegramNotifier{}
				success = notifier.Open(notifierPath)
				break
			}
		case "discord":
			{
				notifier = &notifiers.DiscordNotifier{}
				success = notifier.Open(notifierPath)
				break
			}
		case "email":
			{
				notifier = &notifiers.EmailNotifier{}
				success = notifier.Open(notifierPath)
				break
			}
		case "shoutrrr":
			{
				notifier = &notifiers.ShoutrrrNotifier{}
				success = notifier.Open(notifierPath)
				break
			}
		case "apprise":
			{
				notifier = &notifiers.AppriseNotifier{}
				success = notifier.Open(notifierPath)
				break
			}
		case "file":
			{
				notifier = &notifiers.FileNotifier{}
				success = notifier.Open(notifierPath)
				break
			}
		default:
			{
				web.startupWarning("Did not recognize notifier type:", notifierType)
			}
		}
		if success {
			web.notifiers[name] = notifier
		} else {
			web.startupWarning("Could not add notifier:", name)
		}
	}
}

func (web *Web) pruneDB() {
	log.Println("Starting database pruning")
	var storeNames []string
	web.db.Model(&FilterOutput{}).Distinct().Pluck("name", &storeNames)
	for _, storeName := range storeNames {
		log.Println("Pruning:", storeName)
		var values []FilterOutput
		web.db.Model(&FilterOutput{}).Order("time asc").Find(&values, fmt.Sprintf("name = '%s'", storeName))
		IDs := make([]uint, 0, len(values))
		for i := range values {
			if i > len(values)-3 {
				break
			}
			a := values[i]
			b := values[i+1]
			c := values[i+2]
			if a.Value == b.Value && b.Value == c.Value {
				IDs = append(IDs, b.ID)
			}
		}
		web.db.Delete(&FilterOutput{}, IDs)
	}
}

func (web *Web) notify(notifierKey string, message string) {
	if notifierKey == "All" {
		for _, notifier := range web.notifiers {
			notifier.Message(message)
		}
	} else {
		notifier, exists := web.notifiers[notifierKey]
		if !exists {
			log.Println("Could not find notifier with key:", notifierKey)
			return
		}
		notifier.Message(message)
	}
}

func (web *Web) run() {
	web.router.Run("0.0.0.0:8080")
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
		entryID, exists := web.cronWatch[filter.ID]
		if !exists {
			log.Println("No cron entry for filter", filter.ID, filter.Name)
			continue
		}
		entry := web.cron.Entry(entryID)
		watchMap[filter.WatchID].CronEntry = &entry
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"watches":  watches,
		"warnings": web.startupWarnings,
	})
}

func (web *Web) notifiersView(c *gin.Context) {
	c.HTML(http.StatusOK, "notifiersView", web.notifiers)
}

func (web *Web) notifiersTest(c *gin.Context) {
	notifierName := c.PostForm("notifier_name")
	notifier, exists := web.notifiers[notifierName]
	if !exists {
		c.Redirect(http.StatusSeeOther, "/notifiers/view")
		return
	}
	notifier.Message("GoWatch Test")
	c.Redirect(http.StatusSeeOther, "/notifiers/view")
}

func (web *Web) watchCreate(c *gin.Context) {
	templateFiles, err := EMBED_FS.ReadDir("watchTemplates")
	if err != nil {
		log.Fatalln("Could not load templates from embed FS")
	}
	templates := make([]string, 0, len(templateFiles))
	templates = append(templates, "None")
	for _, template := range templateFiles {
		templateFile := template.Name()
		templateName := templateFile[:len(templateFile)-len(".json")]
		templates = append(templates, templateName)
	}
	c.HTML(http.StatusOK, "watchCreate", gin.H{
		"templates": templates,
	})
}

func (web *Web) watchCreatePost(c *gin.Context) {
	var watch Watch
	errMap, err := bindAndValidateWatch(&watch, c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500", errMap)
		return
	}

	templateID, err := strconv.Atoi(c.PostForm("template"))
	if err != nil {
		log.Println(err)
		templateID = 0
	}

	if templateID == 0 { // empty new watch
		web.db.Create(&watch)
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/watch/edit/%d", watch.ID))
		return
	}

	var jsn []byte
	if templateID == -1 { // watch from url template
		url := c.PostForm("url")
		if len(url) == 0 {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		resp, err := http.Get(url)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		jsn = body
	} else {
		templateFiles, err := EMBED_FS.ReadDir("watchTemplates")
		if err != nil {
			log.Fatalln("Could not load templates from embed FS")
		}

		if templateID >= len(templateFiles) {
			log.Println("/watch/create POSTed with", templateID, "but only", len(templateFiles), "templates")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		template := templateFiles[templateID-1] // -1 because of "None" option
		templatePath := fmt.Sprintf("watchTemplates/%s", template.Name())
		_jsn, err := EMBED_FS.ReadFile(templatePath)
		if err != nil {
			log.Println("Could not read template from embed.FS:", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		jsn = _jsn
	}

	export := WatchExport{}
	if err := json.Unmarshal(jsn, &export); err != nil {
		log.Println("Could not unmarshel JSON:", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	web.db.Create(&watch)

	filterMap := make(map[uint]*Filter)
	for i := range export.Filters {
		filter := &export.Filters[i]
		filterMap[filter.ID] = filter
		filter.ID = 0
		filter.WatchID = watch.ID
	}
	if len(export.Filters) > 0 {
		tx := web.db.Create(&export.Filters)
		if tx.Error != nil {
			log.Println("Create filters transaction failed:", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	for i := range export.Connections {
		connection := &export.Connections[i]
		connection.ID = 0
		connection.WatchID = watch.ID
		connection.OutputID = filterMap[connection.OutputID].ID
		connection.InputID = filterMap[connection.InputID].ID
	}
	if len(export.Connections) > 0 {
		tx := web.db.Create(&export.Connections)
		if tx.Error != nil {
			log.Println("Create connections transaction failed:", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

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
	web.db.Model(&Filter{}).Find(&cronFilters, "watch_id = ? AND type = 'cron' AND var2 = 'yes'", id)
	for _, filter := range cronFilters {
		entryID, exist := web.cronWatch[filter.ID]
		if exist {
			web.cron.Remove(entryID)
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
		entryID, exists := web.cronWatch[filter.ID]
		if !exists {
			log.Println("Could not find entry for filter", filter.ID, filter.Name)
			continue
		}
		entry := web.cron.Entry(entryID)
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
	for name := range names {
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

	notifiers := make([]string, 1)
	notifiers = append(notifiers, "All")
	for notifier := range web.notifiers {
		notifiers = append(notifiers, notifier)
	}

	buildFilterTree(filters, connections)
	processFilters(filters, web, &watch, true, nil)

	c.HTML(http.StatusOK, "watchEdit", gin.H{
		"Watch":       watch,
		"Filters":     filters,
		"Connections": connections,
		"Values":      values,
		"Notifiers":   notifiers,
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
	web.db.Model(&Filter{}).Where("watch_id = ? AND type = 'cron' AND var2 = 'yes'", watch.ID).Find(&cronFilters)
	for _, filter := range cronFilters {
		entryID, exist := web.cronWatch[filter.ID]
		if exist {
			web.cron.Remove(entryID)
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
			if *filter.Var2 == "no" {
				continue
			}
			entryID, err := web.cron.AddFunc(filter.Var1, func() { triggerSchedule(filter.WatchID, web, &filter.ID) })
			if err != nil {
				log.Println("Could not start job for Watch: ", filter.WatchID, err)
				continue
			}
			log.Println("Started CronJob for WatchID", filter.WatchID, "FilterID", filter.ID, "with schedule:", filter.Var1)
			web.cronWatch[filter.ID] = entryID
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
	web.urlCache = make(map[string]string, 5)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
		entryID, exist := web.cronWatch[filter.ID]
		if exist {
			web.cron.Remove(entryID)
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
		for i := range export.Filters {
			filter := &export.Filters[i]
			if filter.Type == "cron" {
				entryID, err := web.cron.AddFunc(filter.Var1, func() { triggerSchedule(filter.WatchID, web, &filter.ID) })
				if err != nil {
					log.Println("Could not start job for Watch: ", filter.WatchID)
					continue
				}
				log.Println("Started CronJob for WatchID", filter.WatchID, "with schedule:", filter.Var1)
				web.cronWatch[filter.ID] = entryID
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
	writeConfFlag := flag.String("writeConfig", "-", "Path to write template config to")
	var printConfigFlag bool
	flag.BoolVar(&printConfigFlag, "printConfig", false, "Print the template config to stdout")
	flag.Parse()

	if *writeConfFlag != "-" {
		conf, err := EMBED_FS.ReadFile("config.tmpl")
		if err != nil {
			log.Fatalln("Could not read config.tmpl")
		}
		os.WriteFile(*writeConfFlag, conf, 0666)
		log.Println("Wrote template config to:", *writeConfFlag)
		return
	}

	if printConfigFlag {
		conf, err := EMBED_FS.ReadFile("config.tmpl")
		if err != nil {
			log.Fatalln("Could not read config.tmpl")
		}
		log.SetFlags(0)
		log.Println(string(conf))
		return
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("GOWATCH")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Could not load config file, using env/defaults")
	}

	web := newWeb()
	web.run()
}
