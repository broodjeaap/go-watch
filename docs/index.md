# GoWatch
[![Build Status](https://drone.broodjeaap.net/api/badges/broodjeaap/go-watch/status.svg)](https://drone.broodjeaap.net/broodjeaap/go-watch)

A change detection server that can notify through various services, written in Go
# Intro

GoWatch works through filters, a filter performs operations on the input it recieves.  
Here is an example of a 'Watch' that calculates the lowest and average price of 4090s on NewEgg and notifies the user if the lowest price changed:  
![NewEgg 4090](images/newegg_4090.png)  

Note that everything, including scheduling/storing/notifying, is a `filter`.  

`Schedule` is a [cron](#cron) filter with a '@every 15m' value, so this will run every 15 minutes.  

`NewEgg Fetch` is a [Get URL](#get-url) filter with a 'https://www.newegg.com/p/pl?N=100007709&d=4090&isdeptsrh=1&PageSize=96' value, it's output will be the HTTP response.  

`Select Price` is a [CSS](#css) filter with the value '.item-container .item-action strong[class!="item-buying-choices-price"]' value, it's output will be the html elements containing the prices.  
An [XPath](#xpath) filter could also have been used.  

`Sanitize` is a [Replace](#replace) filter, using a regular expression ('[^0-9]') it removes anything that's not a number.  

`Avg` is an [Average](#average) filter, it calculates the average value of its inputs.  

`Min` is a [Minimum](#minimum) filter, it calculates the minimum value of its inputs.  

`Average` and `Minimum` are [Store](#store) filters, they store its input values in the database.  

`Diff` is a [Different Than Last](#different-than-last) filter, only passing on the inputs that are different then the last value stored in the database.  

`Notify` is a [Notify](#notify) filter, if there are any inputs to this filter, it will execute a template and send the result to a user defined 'notifier' (Telegram/Discord/etc).
# Run

## Binary

Download 

## Docker

Easiest way to get started is with the prebuilt docker image `ghcr.io/broodjeaap/go-watch:latest`, first get a config template:  
`docker run --rm ghcr.io/broodjeaap/go-watch:latest -printConfig 2> config.yaml`  

Or:  
`docker run --rm -v $PWD:/config ghcr.io/broodjeaap/go-watch:latest -writeConfig /config/config.yaml`

After modifying the config to fit your needs, start the docker container
```
docker run \
    -p 8080:8080 \
    -v $PWD/:/config \
    ghcr.io/broodjeaap/go-watch:latest
```
### Database

By default, GoWatch will use an SQLite database, stored in the `/config` directory for the docker image, which is probably fine for most use cases.  

But you can use another database by changing the `database.dsn` value in the config or `GOWATCH_DATABASE_DSN` environment variable, for example with a PostgreSQL database:  
```
version: "3"

services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    container_name: go-watch
    environment:
    - GOWATCH_DATABASE_DSN=postgres://gorm:gorm@db:5432/gorm
    volumes:
    - /host/path/to/config:/config
    ports:
    - "8080:8080"
    depends_on:
      db:
        condition: service_healthy
  db:
    image: postgres:15
    environment:
    - POSTGRES_USER=gorm
    - POSTGRES_PASSWORD=gorm
    - POSTGRES_DB=gorm
    volumes:
    - /host/path/to/db:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -d $${POSTGRES_DB} -U $${POSTGRES_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
```

#### Pruning

An automatic database prune job that removes repeating values can be scheduled by adding a cron schedule to the config:  
```
database:
  dsn: "/config/watch.db"
  prune: "@every 1h"
```

### Proxy

GoWatch has some basic proxy support, using the config we can point GoWatch to a proxy server:  
```
proxy:
  proxy_url: http://proxy.com:1234
```

When using the docker image, the `HTTP_PROXY` and `HTTPS_PROXY` environment variables can also be used:  
```
services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    container_name: go-watch
    environment:
    - HTTP_PROXY=http://proxy.com:1234
    - HTTPS_PROXY=http://proxy.com:1234
```
#### Proxy pools

Proxy 'pools' are not directly supported by GoWatch, but can still be set up by using a proxy, for example with [Squid](http://www.squid-cache.org/):  
```
services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    container_name: go-watch
    environment:
    - HTTP_PROXY=http://squid_proxy:3128
    - HTTPS_PROXY=http://squid_proxy:3128
  squid_proxy:
    image: sameersbn/squid:latest
    volumes:
    - /path/to/squid.conf:/etc/squid/squid.conf
```

And in the `squid.conf` the proxy pool would be defined like this:  
```
cache_peer proxy1.com parent 3128 0 round-robin no-query
cache_peer proxy2.com parent 3128 0 round-robin no-query login=user:pass
```

### Browserless

Some websites don't send all content on the first request, it's added later through javascript, Amazon does this for example.  
To still be able to watch products from these websites, GoWatch supports [Browserless](https://www.browserless.io/), the Browserless URL can be added to the config:  
```
browserless:
  url: http://your.browserless:3000/content
```

Or as an environment variable, for example in a docker-compose:  
```
version: "3"

services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    container_name: go-watch
    environment:
    - GOWATCH_BROWSERLESS_URL=http://browserless:3000/content
    volumes:
    - /host/path/to/config:/config
    ports:
    - "8080:8080"
  browserless:
    image: browserless/chrome:latest
```

Note that the proxy environment variables can be added to the Browserless container to still allow for proxying.

### Authentication

GoWatch doesn't have built in authentication, but we can use a reverse proxy for that, for example through Traefik:  
```
version: "3"

services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    container_name: go-watch
    environment:
    - GOWATCH_DATABASE_DSN=postgres://gorm:gorm@db:5432/gorm
    volumes:
    - /host/path/to/config:/config
    ports:
    - "8181:8080"
    depends_on:
      db:
        condition: service_healthy
    labels:
    - "traefik.http.routers.gowatch.rule=Host(`192.168.178.254`)"
    - "traefik.http.routers.gowatch.middlewares=test-auth"
  db:
    image: postgres:15
    environment:
    - POSTGRES_USER=gorm
    - POSTGRES_PASSWORD=gorm
    - POSTGRES_DB=gorm
    volumes:
    - /host/path/to/db:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -d $${POSTGRES_DB} -U $${POSTGRES_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
    depends_on:
    - proxy
  proxy:
    image: traefik:v2.9.6
    command: --providers.docker
    labels:
    - "traefik.http.middlewares.test-auth.basicauth.users=broodjeaap:$$2y$$10$$aUvoh7HNdt5tvf8PYMKaaOyCLD3Uel03JtEIPxFEBklJE62VX4rD6"
    ports:
    - "8080:80"
    volumes:
    - /var/run/docker.sock:/var/run/docker.sock
```

Change the `Host` label to the correct ip/hostname and generate a user/password string with [htpasswd](https://httpd.apache.org/docs/2.4/programs/htpasswd.html) for the `basicauth.users` label, note that the `$` character is escaped with `$$`

# Filters

GoWatch comes with many filters that should, hopefully, be enough to allow for most use cases.  

## Cron

The `cron` filter is used to schedule when your watch will run.  
It uses the [cron](https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.0#section-readme) package to schedule go routines, some common examples would be:  
- `@every 15m`: will trigger every 15 minutes starting on server start.
- `@hourly`: will trigger every beginning of hour.
- `30 * * * *`: will trigger every hour on the half hour.  

For more detailed instructions you can check its documentation.

## Get URL

Fetches the given URL and outputs the HTTP response.  
For more complicated requests, POSTing/headers/login, use the [HTTP functionality](https://github.com/vadv/gopher-lua-libs/tree/master/http#client-1) in the Lua filter.  
During editing, http requests are cached, so not to trigger any DOS protection on your sources.

## Get URLs

Fetches every URL given as input and outputs every HTTP response.  
During editing, http requests are cached, so not to trigger any DOS protection on your sources.

## CSS

Use a [CSS selector](https://www.w3schools.com/cssref/css_selectors.php) to filter your http responses.  
The [Cascadia](https://github.com/andybalholm/cascadia) package is used for this filter, check the docs to see what is and isn't supported.

## XPath

Use an [XPath](https://www.w3schools.com/xml/xpath_intro.asp) to filter your http responses.  
The [XPath](https://github.com/antchfx/xpath) package is used for this filter, check the docs to see what is and isn't supported.

## JSON

Use a this to filter your JSON responses, the [gjson](https://github.com/tidwall/gjson) package is used for this filter.  
Some common examples would be:  
- product.price
- items.3
- products.#.price

## Replace

Simple replace filter, supports regular expressions.  
If the `With` value is empty, it will just remove matching text.

## Match

Searches for the regex, outputs every match.

## Substring

Substring allows for a [Python like](https://learnpython.com/blog/substring-of-string-python/) substring selection.  
For the input string 'Hello World!':  
- `:5`: Hello
- `6:`: World!
- `6,0,7`: WHo
- `-6:`: World!
- `-6:,:5`: World!Hello

## Contains

Inputs pass if they contain the given value, basically same as Match but more basic.

## Store

Stores each input value in the database under its own name, should probably limit this to single inputs (after Minimum/Maximum/Average filters).

## Notify

Executes the given template and sends the resulting string as a message to the given notifier(s).  
It uses the [Golang templating language](https://pkg.go.dev/text/template), the outputs of all the filters can be used by the name of the filters.  
So if you have a `Min` filter like in the example, it can be referenced in the template by using `{{ .Min }}`.  
The name of the watch is also included under `.WatchName`.  

To configure notifiers see the [notifiers](#notifiers) section.

## Math

### Sum

Sums the inputs together, nonnumerical values are skipped.
### Minimum
Outputs the lowest value of the inputs, nonnumerical values are skipped.
### Maximum
Outputs the highest value of the inputs, nonnumerical values are skipped.

### Average
Outputs the average of the inputs, nonnumerical values are skipped.

### Count
Outputs the number of inputs.
### Round
Outputs the inputs rounded to the given decimals, nonnumerical valuesa are skipped.
## Condition

### Different Than Last

Passes an input if it is different than the last stored value.

### Lower Than Last

Passes an input if it is lower than the last stored value.

### Lowest

Passes an input if it is lower than all previous stored values.

### Lower Than

Passes an input if it is lower than a given value.

### Higher Than Last

Passes an input if it is higher than the last stored value.

### Highest
Passes an input if it is higher than all previous stored values.

### Higher Than

Passes an input if it is higher than a given value.

## Lua

The Lua filter wraps [gopher-lua](https://github.com/yuin/gopher-lua), with [gopher-lua-libs](https://github.com/vadv/gopher-lua-libs) to greatly extend the capabilities of the Lua VM.  
A basic script that just passes all inputs to the output looks like this:  
```
for i,input in pairs(inputs) do
	table.insert(outputs, input)
end
```

Both `inputs` and `outputs` are convenience tables provided by GoWatch to make Lua scripting a bit easier.
There is also a `logs` table that can be used the same way as the `outputs` table (`table.insert(logs, 'this will be logged')`) to provide some basic logging.  

Much of the functionality that is provided through individual filters in GoWatch can also be done from Lua.  
The gopher-lua-libs provide an [http](https://github.com/vadv/gopher-lua-libs/tree/master/http) lib, whose output can be parsed with the [xmlpath](https://github.com/vadv/gopher-lua-libs/tree/master/xmlpath) or [json](https://github.com/vadv/gopher-lua-libs/tree/master/json) libs and then filtered with a [regular expression](https://github.com/vadv/gopher-lua-libs/tree/master/regexp) or some regular Lua scripting to then finally be turned into a ready to send notification through a [template](https://github.com/vadv/gopher-lua-libs/tree/master/template).  

# Notifiers

The basic form is:  
```
notifiers:
  <notifier-name>:
    type: "<notifier-type>"
    other: ""
    values: ""
    etc: ""
```
See below for possible types.
## Shoutrrr

[Shoutrrr](https://containrrr.dev/shoutrrr/v0.5/) can be used to notify many different services, check their docs for a [list](https://containrrr.dev/shoutrrr/v0.5/services/overview/) of which ones.  
An example config for sending notifications through Shoutrrr:   
```
notifiers:
  Shoutrrr-telegram-discord:
    type: "shoutrrr"
    urls:
    - telegram://<token>@telegram?chats=<channel-1-id>,<chat-2-id>
    - discord://<token>@<webhookid>
    - etc...
database:
  dsn: "watch.db"
  prune: "@every 1h"
```

## Apprise

[Apprise](https://github.com/caronc/apprise) can also be used to send notifications, it supports many different services/protocols, but it requires access to an [Apprise API](https://github.com/caronc/apprise-api).  
Luckily there is a [docker image](https://hub.docker.com/r/caronc/apprise) available that we can add to our compose:  
```
version: "3"

services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    container_name: go-watch
    volumes:
    - /host/path/to/:/config
    ports:
    - "8080:8080"
  apprise:
    image: caronc/apprise:latest
```

And the notifier config:  
```
notifiers:
  apprise:
    type: "apprise"
    url: "http://apprise:8000/notify"
    urls:
    - "tgram://<bot_token>/<chat_id>/"
    - "discord://<WebhookID>/<WebhookToken>/"
database:
  dsn: "watch.db"
  prune: "@every 1h"
```

## Telegram

Get a bot token from the [@BotFather](https://core.telegram.org/bots/tutorial) and get the [chatID](https://www.alphr.com/find-chat-id-telegram/) of the chat want to send the notifications to.  
An example config for sending notifications through Telegram:  
```
notifiers:
  Telegram:
    token: "<token>"
    chat: "<chatID>"
    debug: false
database:
  dsn: "watch.db"
  prune: "@every 1h"
```

## Discord

To get a token, userID and/or serverID you first have to enable [Developer Mode and create a new application](https://www.ionos.com/digitalguide/server/know-how/creating-discord-bot/).  
Then you can right click on your username in any chat to copy your user ID or right click on a server/channel to get the server/channel ID.  
An example config for sending DM notifications through Discord:  
```
notifiers:
  Discord:
    type: "discord"
    token: "<token>"
    userID: "<userID>"
    debug: false
database:
  dsn: "watch.db"
  prune: "@every 1h"
```

An example config for sending channel notifications:  
```
notifiers:
  Discord:
    type: "discord"
    token: "<token>"
    server:
      ID: "<serverID>"
      channel: "<channelID>"
    debug: false
database:
  dsn: "watch.db"
  prune: "@every 1h"
```

Both a userID and serverID/channelID is also possible.

## Email

An example config for sending email notifications through a SMTP relay server:  
```
notifiers:
  Email-to-at-email-com:
    type: "email"
    server: "smtp.relay.com"
    port: "465"
    from: "from@email.com"
    user: "apikey"
    password: "-"
    to: "to@email.com"
database:
  dsn: "watch.db"
  prune: "@every 1h"
```

# Dev
## type script compilation

`tsc static/*.ts --lib es2020,dom --watch --downlevelIteration`

# Dependencies

The following libaries are used in Go-Watch:  
- [Gin](https://github.com/gin-gonic/gin) for HTTP server
    - [multitemplate](https://github.com/gin-contrib/multitemplate) for template inheritance
- [Cascadia](https://pkg.go.dev/github.com/andybalholm/cascadia) for CSS selectors
- [htmlquery](https://pkg.go.dev/github.com/antchfx/htmlquery) for XPath selectors
- [validator](https://pkg.go.dev/github.com/go-playground/validator/v10@v10.11.0) for user user input validation
- Notifiers
  - [Shoutrrr](https://github.com/containrrr/shoutrrr/) for many different services
  - [tgbotapi](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5@v5.5.1) for Telegram
  - [discordgo](https://pkg.go.dev/github.com/bwmarrin/discordgo) for Discord
  - [gomail](https://pkg.go.dev/gopkg.in/gomail.v2?utm_source=godoc) for Email
- [cron](https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.0) for job scheduling
- [viper](https://pkg.go.dev/github.com/spf13/viper@v1.12.0) for config management
- [gjson](https://pkg.go.dev/github.com/tidwall/gjson@v1.14.2) for JSON selectors
- [gopher-lua](https://github.com/yuin/gopher-lua) for Lua scripting
    - [gopher-lua-libs](https://pkg.go.dev/github.com/vadv/gopher-lua-libs@v0.4.0) for expanding the Lua scripting functionality
- [net](https://pkg.go.dev/golang.org/x/net) for http fetching
- [gorm](https://pkg.go.dev/gorm.io/gorm@v1.23.8) for database abstraction
    - [sqlite](https://pkg.go.dev/gorm.io/driver/sqlite@v1.3.6)
    - [postgres](https://github.com/go-gorm/postgres)
    - [mysql](https://github.com/go-gorm/mysql)
    - [sqlserver](https://github.com/go-gorm/sqlserver)
- [bootstrap](https://getbootstrap.com/)