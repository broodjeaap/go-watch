package main

import (
	"bytes"
	"compress/gzip"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	router          *gin.Engine                   // gin router instance
	templates       multitemplate.Renderer        // multitemplate instance
	cron            *cron.Cron                    // cron instance
	urlCache        map[string]string             // holds url -> http response
	cronWatch       map[uint]cron.EntryID         // holds cronFilter.ID -> EntryID
	db              *gorm.DB                      // gorm db instance
	notifiers       map[string]notifiers.Notifier // holds notifierName -> notifier
	startupWarnings []string                      // simple list of warnings/errors found during startup, displayed on / page
}

// NewWeb creates a new web instance and calls .init() before returning it
func newWeb() *Web {
	web := &Web{
		urlCache:        make(map[string]string, 5),
		startupWarnings: make([]string, 0, 10),
	}
	web.init()
	return web
}

// init initializes DB, routers, cron jobs and notifiers
func (web *Web) init() {
	web.urlCache = make(map[string]string, 5)
	if !viper.GetBool("gin.debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	web.validateProxyURL()
	web.initDB()
	web.initRouter()
	web.initNotifiers()
	go web.initCronJobs()
}

// startupWarning is a helper function to add a message to web.startupWarnings and print it to stdout
func (web *Web) startupWarning(m ...any) {
	warning := fmt.Sprint(m...)
	log.Println(warning)
	web.startupWarnings = append(web.startupWarnings, warning)
}

func (web *Web) addCronJobIfCronFilter(filter *Filter, startup bool) {
	if filter.ID == 0 {
		return
	}
	if filter.Type != "cron" {
		return
	}
	if filter.Var2 != nil && *filter.Var2 == "no" {
		return
	}
	entryID, err := web.cron.AddFunc(filter.Var1, func() { triggerSchedule(filter.WatchID, web, &filter.ID) })
	if err != nil {
		if startup {
			web.startupWarning("Could not start job for Watch: ", filter.WatchID)
		} else {
			log.Println("Could not start job for Watch: ", filter.WatchID)
		}
		return
	}
	log.Println("Started CronJob for WatchID", filter.WatchID, "with schedule:", filter.Var1)
	web.cronWatch[filter.ID] = entryID
}

// validateProxyURL calls url.Parse with the proxy.proxy_url, if there is an error, it's added to startupWarnings
func (web *Web) validateProxyURL() {
	if viper.IsSet("proxy.proxy_url") {
		_, err := url.Parse(viper.GetString("proxy.proxy_url"))
		if err != nil {
			web.startupWarning("Could not parse proxy url, check config")
			return
		}
	}
}

// initDB initializes the database with the database.dsn value.
func (web *Web) initDB() {
	dsn := "./watch.db"
	if viper.IsSet("database.dsn") {
		dsn = viper.GetString("database.dsn")
	}

	conf := &gorm.Config{}
	conf.PrepareStmt = true

	var db *gorm.DB
	var err error
	if strings.HasPrefix(dsn, "sqlserver") {
		db, err = gorm.Open(sqlserver.Open(dsn), conf)
		log.Println("Using SQLServer database")
	} else if strings.HasPrefix(dsn, "postgres") {
		db, err = gorm.Open(postgres.Open(dsn), conf)
		log.Println("Using PostgreSQL database")
	} else if strings.HasPrefix(dsn, "mysql") {
		db, err = gorm.Open(mysql.Open(dsn), conf)
		log.Println("Using MySQL database")
	} else {
		db, err = gorm.Open(sqlite.Open(dsn), conf)
		log.Println("Using sqlite database at:", dsn)
	}
	if db == nil {
		log.Panicln("Could not recognize database.dsn: ", dsn)
	}
	if err != nil {
		log.Panicln("Could not start DB: ", err)
	}

	// retry connection to the db a couple times with exp retry time
	web.db = db
	delay := time.Duration(1) * time.Second
	maxDelay := time.Duration(8) * time.Second
	for {
		err := web.db.AutoMigrate(&Watch{}, &Filter{}, &FilterConnection{}, &FilterOutput{})
		if err == nil {
			log.Println("Connected to DB!")
			break
		}
		if delay >= maxDelay {
			web.startupWarning("Could not initialize database", err)
			break
		}
		log.Println("Could not connect to DB, retry in:", delay.String())
		time.Sleep(delay)
		delay *= 2
	}
}

// initRouer initializes the GoWatch routes, binding web.func to a url path
func (web *Web) initRouter() {
	web.router = gin.Default()

	staticFS, err := fs.Sub(EMBED_FS, "static")
	if err != nil {
		log.Fatalln("Could not load static embed fs")
	}
	web.router.StaticFS("/static", http.FS(staticFS))
	web.router.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(staticFS))

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

	web.router.GET("/backup/view", web.backupView)
	web.router.GET("/backup/create", web.backupCreate)
	web.router.GET("/backup/test/:id", web.backupTest)

	web.router.SetTrustedProxies(nil)
}

// initTemplates initializes the templates from EMBED_FS/templates
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

	web.templates.Add("backupView", template.Must(template.ParseFS(templatesFS, "base.html", "backup/view.html")))
	web.templates.Add("backupTest", template.Must(template.ParseFS(templatesFS, "base.html", "backup/test.html")))

	web.templates.Add("500", template.Must(template.ParseFS(templatesFS, "base.html", "500.html")))
}

// initCronJobs reads any 'cron' type filters from the database, and starts a cron job for each
func (web *Web) initCronJobs() {
	var cronFilters []Filter

	// type cron and enabled = yes
	web.db.Model(&Filter{}).Find(&cronFilters, "type = 'cron' AND var2 = 'yes'")

	web.cronWatch = make(map[uint]cron.EntryID, len(cronFilters))

	web.cron = cron.New()
	web.cron.Start()

	// db prune job is started if there is a database.prune set
	if viper.IsSet("database.prune") {
		pruneSchedule := viper.GetString("database.prune")
		_, err := web.cron.AddFunc(pruneSchedule, web.pruneDB)
		if err != nil {
			web.startupWarning("Could not parse database.prune:", err)
		}
		log.Println("Started DB prune cronjob:", pruneSchedule)
	}

	// backup job is started if there is a schedule and path
	if viper.IsSet("database.backup.schedule") && viper.IsSet("database.backup.path") {
		backupSchedule := viper.GetString("database.backup.schedule")
		_, err := web.cron.AddFunc(backupSchedule, web.scheduledBackup)
		if err != nil {
			web.startupWarning("Could not parse database.backup.schedule:", err)
		}
		log.Println("Backup schedule set:", backupSchedule)
	}

	// add some delay to cron jobs, so watches with the same schedule don't
	// 'burst' at the same time after restarting GoWatch
	cronDelayStr := viper.GetString("schedule.delay")
	cronDelay, delayErr := time.ParseDuration(cronDelayStr)
	if delayErr == nil {
		log.Println("Delaying job startup by:", cronDelay.String())
	} else {
		web.startupWarning("Could not parse schedule.delay: ", cronDelayStr)
	}

	// for every cronFilter, add a new cronjob with the schedule in filter.var1
	for i := range cronFilters {
		cronFilter := &cronFilters[i]
		web.addCronJobIfCronFilter(cronFilter, true)
		if delayErr == nil {
			time.Sleep(cronDelay)
		}
	}
}

// initNotifiers initializes the notifiers configured in the config
func (web *Web) initNotifiers() {
	web.notifiers = make(map[string]notifiers.Notifier, 5)
	if !viper.IsSet("notifiers") {
		web.startupWarning("No notifiers set!")
		return
	}

	// iterates over the map of notifiers, key being the name of the notifier
	notifiersMap := viper.GetStringMap("notifiers")
	for name := range notifiersMap {

		// should probably use notifiersMap.Sub(name), but if it aint broke...
		notifierPath := fmt.Sprintf("notifiers.%s", name)
		notifierMap := viper.GetStringMapString(notifierPath)

		notifierType, exists := notifierMap["type"]
		if !exists {
			web.startupWarning(fmt.Sprintf("No 'type' for '%s' notifier!", name))
			continue
		}

		// create an empty notifier and a success flag,
		// so we can add it to the map at the end instead of each switch case
		success := false
		var notifier notifiers.Notifier
		switch notifierType {
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

// pruneDB is called by the pruneDB cronjob, it removes repeating values from the database.
func (web *Web) pruneDB() {
	log.Println("Starting database pruning")

	// for every unique (watch.ID, storeFilter.Name)
	var storeFilters []Filter
	web.db.Model(&FilterOutput{}).Distinct("watch_id", "name").Find(&storeFilters)
	for _, storeFilter := range storeFilters {
		// get all the values out of the database
		var values []FilterOutput
		tx := web.db.Model(&FilterOutput{}).Order("time asc").Find(&values, fmt.Sprintf("watch_id = %d AND name = '%s'", storeFilter.WatchID, storeFilter.Name))
		if tx.Error != nil {
			continue
		}

		// a value can be deleted if it's the same as the previous and next value
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
		if len(IDs) > 0 {
			log.Println("Pruned: ", storeFilter.WatchID, "-", storeFilter.Name, "removed", len(IDs), "values")
			web.db.Delete(&FilterOutput{}, IDs)
		} else {
			log.Println("Nothing to prune for", storeFilter.WatchID, "-", storeFilter.Name)
		}
	}
}

// notify sends a message to the notifier with notifierKey name, or all notifiers if key is 'All'
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

// run simply calls router.Run
func (web *Web) run() {
	web.router.Run("0.0.0.0:8080")
}

// index (/) displays all watches with name, last run, next run, last value
func (web *Web) index(c *gin.Context) {
	watches := []Watch{}
	web.db.Find(&watches)

	// make a map[watch.ID] -> watch so after this we can add data to watches in O(1)
	watchMap := make(map[uint]*Watch, len(watches))
	for i := 0; i < len(watches); i++ {
		watchMap[watches[i].ID] = &watches[i]
	}
	// get the schedule for this watch, so we can display last/next run
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

	// get the last value stored, also doesn't really work with multiple values but meh again
	rows, err := web.db.Table("watches").
		Select("watches.id, max(filter_outputs.time) as time, filter_outputs.value").
		Joins("left join filter_outputs on filter_outputs.watch_id = watches.id").
		Order("filter_outputs.name").
		Group("watches.id, time, filter_outputs.value, filter_outputs.name").
		Rows()

	if err != nil {
		log.Println(err)
	} else {
		for rows.Next() {
			var watchID uint
			var _time string
			var value string
			err := rows.Scan(&watchID, &_time, &value)
			if err != nil {
				log.Println(err)
				continue
			}
			watchMap[watchID].LastValue = value
		}
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"watches":  watches,
		"warnings": web.startupWarnings,
	})
}

// notifiersView (/notifiers/view) shows the notifiers and a test button
func (web *Web) notifiersView(c *gin.Context) {
	c.HTML(http.StatusOK, "notifiersView", web.notifiers)
}

// notifiersTest (/notifiers/test) sends a test message to notifier_name
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

// watchCreate (/watch/create) allows user to create a new watch
// A name and an optional template can be picked
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

// watchCreatePost (/watch/create) is where a new watch create will be submitted to
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
		return // nothing else to do
	}

	// get the template either from a url or from one of the template files
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
	} else if templateID == -2 { // watch from file upload
		file, err := c.FormFile("file")

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		jsn, _ = ioutil.ReadAll(openedFile)
	} else { // selected one of the templates
		templateFiles, err := EMBED_FS.ReadDir("watchTemplates")
		if err != nil {
			log.Fatalln("Could not load templates from embed FS")
		}

		if templateID > len(templateFiles) {
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

	// create the watch
	web.db.Create(&watch)

	// the IDs of filters and connections have to be 0 when they are added to the database
	// otherwise they will overwrite whatever filters/connections happened to have the same ID
	// so we set them all to 0, but keep a map of 'old filter ID' -> filter
	filterMap := make(map[uint]*Filter)
	for i := range export.Filters {
		filter := &export.Filters[i]
		filterMap[filter.ID] = filter
		filter.ID = 0
		filter.WatchID = watch.ID
	}
	if len(export.Filters) > 0 {
		// after web.db.Create, the filters will have their new IDs
		tx := web.db.Create(&export.Filters)
		if tx.Error != nil {
			log.Println("Create filters transaction failed:", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	for i := range export.Filters {
		filter := &export.Filters[i]
		web.addCronJobIfCronFilter(filter, false)
	}

	// we again set all the connection.ID to 0,
	// but then also swap the old filterIDs to the new IDs the filters got after db.Create
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

// deleteWatch (/watch/delete) removes a watch and it's cronjobs
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

// watchView (/watch/view) shows the watch page with a graph and/or a table of stored values
func (web *Web) watchView(c *gin.Context) {
	id := c.Param("id")

	var watch Watch
	web.db.Model(&Watch{}).First(&watch, id)

	// get the cron filter for this watch
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

	// get all the values from this watch from the database
	// split it in 2 groups, numerical and categorical
	var values []FilterOutput
	web.db.Model(&FilterOutput{}).Order("time asc").Where("watch_id = ?", watch.ID).Find(&values)
	numericalMap := make(map[string][]*FilterOutput, len(values))
	categoricalMap := make(map[string][]*FilterOutput, len(values))
	names := make(map[string]bool, 5)
	for i := range values {
		value := &values[i]
		names[value.Name] = true
		_, err := strconv.ParseFloat(value.Value, 64)
		if err == nil {
			numericalMap[value.Name] = append(numericalMap[value.Name], value)
		} else {
			// probably very inefficient to prepend, but want newest values at the top
			categoricalMap[value.Name] = append([]*FilterOutput{value}, categoricalMap[value.Name]...)
		}
	}

	// give value groups a color, defined in templates/watch/view.html
	colorMap := make(map[string]int, len(names))
	index := 0
	for name := range names {
		colorMap[name] = index % 16 // only 16 colors
		index += 1
	}

	c.HTML(http.StatusOK, "watchView", gin.H{
		"Watch":          watch,
		"numericalMap":   numericalMap,
		"categoricalMap": categoricalMap,
		"colorMap":       colorMap,
	})
}

// watchEdit (/watch/edit) shows the node diagram of a watch, allowing the user to modify it
func (web *Web) watchEdit(c *gin.Context) {
	id := c.Param("id")

	var watch Watch
	web.db.Model(&Watch{}).First(&watch, id)

	var filters []Filter
	web.db.Model(&Filter{}).Where("watch_id = ?", watch.ID).Find(&filters)

	var connections []FilterConnection
	web.db.Model(&FilterConnection{}).Where("watch_id = ?", watch.ID).Find(&connections)

	notifiers := make([]string, 1)
	notifiers = append(notifiers, "All")
	for notifier := range web.notifiers {
		notifiers = append(notifiers, notifier)
	}

	// add all the filters to
	buildFilterTree(filters, connections)
	processFilters(filters, web, &watch, true, nil)

	c.HTML(http.StatusOK, "watchEdit", gin.H{
		"Watch":       watch,
		"Filters":     filters,
		"Connections": connections,
		"Notifiers":   notifiers,
	})
}

// watchUpdate (/watch/update) is where /watch/edit POSTs to
func (web *Web) watchUpdate(c *gin.Context) {
	// So this function is a simple/lazy way of implementing a watch update.
	// the watch is the only thing that gets 'updated', the rest is just wiped from the database
	// and reinserted
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
			web.addCronJobIfCronFilter(filter, false)
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

// cacheView (/cache/view) shows the items in the web.urlCache
func (web *Web) cacheView(c *gin.Context) {
	c.HTML(http.StatusOK, "cacheView", web.urlCache)
}

// cacheClear (/cache/clear) clears all items in web.urlCache
func (web *Web) cacheClear(c *gin.Context) {
	web.urlCache = make(map[string]string, 5)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// backupView (/backup/view) lists the stored backups
func (web *Web) backupView(c *gin.Context) {
	if !viper.IsSet("database.backup") {
		c.HTML(http.StatusOK, "backupView", gin.H{"Error": "database.backup not set"})
		return
	}
	if !viper.IsSet("database.backup.schedule") {
		c.HTML(http.StatusOK, "backupView", gin.H{"Error": "database.backup.schedule not set"})
		return
	}
	if !viper.IsSet("database.backup.path") {
		c.HTML(http.StatusOK, "backupView", gin.H{"Error": "database.backup.path not set"})
		return
	}

	backupPath := viper.GetString("database.backup.path")

	backupDir, err := filepath.Abs(filepath.Dir(backupPath))
	if err != nil {
		c.HTML(http.StatusOK, "backupView", gin.H{"Error": err})
		return
	}

	filesInBackupDir, err := ioutil.ReadDir(backupDir)
	if err != nil {
		c.HTML(http.StatusOK, "backupView", gin.H{"Error": err})
		return
	}

	filePaths := make([]string, 0, len(filesInBackupDir))
	for _, fileInBackupDir := range filesInBackupDir {
		fullPath := filepath.Join(backupDir, fileInBackupDir.Name())
		filePaths = append(filePaths, fullPath)
	}

	c.HTML(http.StatusOK, "backupView", gin.H{"Backups": filePaths})
}

// backupCreate (/backup/create) creates a new backup
func (web *Web) backupCreate(c *gin.Context) {
	if !viper.IsSet("database.backup") {
		c.HTML(http.StatusBadRequest, "backupView", gin.H{"Error": "database.backup not set"})
		return
	}
	if !viper.IsSet("database.backup.path") {
		c.HTML(http.StatusBadRequest, "backupView", gin.H{"Error": "database.backup.path not set"})
		return
	}
	backupDir := filepath.Dir(viper.GetString("database.backup.path"))
	backupName := fmt.Sprintf("gowatch_%s.gzip", time.Now().Format(time.RFC3339))
	backupName = strings.Replace(backupName, ":", "-", -1)

	backupPath := filepath.Join(backupDir, backupName)
	err := web.createBackup(backupPath)
	if err != nil {
		c.HTML(http.StatusBadRequest, "backupView", gin.H{"Error": err})
		return
	}
	c.Redirect(http.StatusSeeOther, "/backup/view")
}

func (web *Web) scheduledBackup() {
	log.Println("Starting scheduled backup")
	backupPath := viper.GetString("database.backup.path")

	// compare abs backup path to abs dir path, if they are the same, it's a dir
	// avoids an Open(backupPath).Stat, which will fail if it's a file that doesn't exist
	absBackupPath, err := filepath.Abs(backupPath)
	if err != nil {
		log.Println("Could not get abs path of database.backup.path")
		return
	}

	backupDir, err := filepath.Abs(filepath.Dir(backupPath))
	if err != nil {
		log.Println("Could not get abs path of dir(database.backup.path)")
		return
	}
	if absBackupPath == backupDir {
		backupName := fmt.Sprintf("gowatch_%s.gzip", time.Now().Format(time.RFC3339))
		backupPath = filepath.Join(backupPath, backupName)
		log.Println(backupPath)
	} else {
		backupTemplate, err := template.New("backup").Parse(backupPath)
		if err != nil {
			log.Println("Could not parse backup path as template:", err)
			return
		}
		var backupNameBytes bytes.Buffer
		err = backupTemplate.Execute(&backupNameBytes, time.Now())
		if err != nil {
			log.Println("Could not execute backup template:", err)
			return
		}
		backupPath = backupNameBytes.String()
	}
	err = web.createBackup(backupPath)
	if err != nil {
		log.Println("Could not create scheduled backup:", err)
		return
	}
	log.Println("Backup succesful:", backupPath)
}

// createBackup is the function that actually creates the backup
func (web *Web) createBackup(backupPath string) error {
	backupFile, err := os.OpenFile(backupPath, os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer backupFile.Close()

	backupWriter := gzip.NewWriter(backupFile)
	defer backupWriter.Close()

	var watches []Watch
	tx := web.db.Find(&watches)
	if tx.Error != nil {
		return tx.Error
	}

	var filters []Filter
	tx = web.db.Find(&filters)
	if tx.Error != nil {
		return tx.Error
	}

	var connections []FilterConnection
	tx = web.db.Find(&connections)
	if tx.Error != nil {
		return tx.Error
	}

	var values []FilterOutput
	tx = web.db.Find(&values)
	if tx.Error != nil {
		return tx.Error
	}

	backup := Backup{
		Watches:     watches,
		Filters:     filters,
		Connections: connections,
		Values:      values,
	}

	jsn, err := json.Marshal(backup)
	if err != nil {
		return err
	}
	_, err = backupWriter.Write(jsn)

	if err != nil {
		return err
	}

	return nil
}

// backupTest (/backup/test) tests the selected backup file
func (web *Web) backupTest(c *gin.Context) {
	importID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if importID < 0 {
		c.Redirect(http.StatusSeeOther, "/backup/view")
		return
	}
	if !viper.IsSet("database.backup") {
		c.HTML(http.StatusOK, "backupTest", gin.H{"Error": "database.backup not set"})
		return
	}
	if !viper.IsSet("database.backup.schedule") {
		c.HTML(http.StatusOK, "backupTest", gin.H{"Error": "database.backup.schedule not set"})
		return
	}
	if !viper.IsSet("database.backup.path") {
		c.HTML(http.StatusOK, "backupTest", gin.H{"Error": "database.backup.path not set"})
		return
	}

	backupPath := viper.GetString("database.backup.path")

	backupDir, err := filepath.Abs(filepath.Dir(backupPath))
	if err != nil {
		c.HTML(http.StatusOK, "backupTest", gin.H{"Error": err})
		return
	}

	filesInBackupDir, err := ioutil.ReadDir(backupDir)
	if err != nil {
		c.HTML(http.StatusOK, "backupTest", gin.H{"Error": err})
		return
	}
	if importID >= len(filesInBackupDir) {
		c.Redirect(http.StatusSeeOther, "/backup/view")
		return
	}

	backupFileName := filesInBackupDir[importID]
	backupFullPath := filepath.Join(backupDir, backupFileName.Name())
	backupFile, err := os.Open(backupFullPath)
	if err != nil {
		c.HTML(http.StatusOK, "backupTest", gin.H{"Error": err})
		return
	}
	defer backupFile.Close()

	backupReader, err := gzip.NewReader(backupFile)
	if err != nil {
		c.HTML(http.StatusOK, "backupTest", gin.H{"Error": err})
		return
	}
	defer backupReader.Close()
	rawBytes, err := io.ReadAll(backupReader)

	var backup Backup
	json.Unmarshal(rawBytes, &backup)

	c.HTML(http.StatusOK, "backupTest", gin.H{
		"Backup":     backup,
		"BackupPath": backupFullPath,
	})
}

// exportWatch (/watch/export/:id) creates a json export of the current watch
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

// importWatch (/watch/import/:id) takes a json file and imports it to the current watch
func (web *Web) importWatch(c *gin.Context) {
	watchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	importType := c.PostForm("type")
	if !(importType == "clear" || importType == "add") {
		c.AbortWithError(http.StatusBadRequest, errors.New("Unknown Import Type"))
		return
	}
	clearFilters := importType == "clear"

	offsetX := 0
	offsetY := 0

	if !clearFilters {
		offsetX, _ = strconv.Atoi(c.PostForm("offset_x"))
		offsetY, _ = strconv.Atoi(c.PostForm("offset_y"))
		offsetX *= -1
		offsetY *= -1
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
	if !clearFilters {
		web.db.Model(&Filter{}).Where("watch_id = ? AND type = 'cron'", watchID).Find(&cronFilters)
	}
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
		filter.X += offsetX
		filter.Y += offsetY
		filter.WatchID = uint(watchID)
	}

	if clearFilters {
		web.db.Delete(&Filter{}, "watch_id = ?", watchID)
	}

	if len(export.Filters) > 0 {
		tx := web.db.Create(&export.Filters)
		if tx.Error != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		for i := range export.Filters {
			filter := &export.Filters[i]
			web.addCronJobIfCronFilter(filter, false)
		}
	}

	if clearFilters {
		web.db.Delete(&FilterConnection{}, "watch_id = ?", watchID)
	}
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
