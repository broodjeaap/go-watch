# GoWatch
[![Build Status](https://drone.broodjeaap.net/api/badges/broodjeaap/go-watch/status.svg)](https://drone.broodjeaap.net/broodjeaap/go-watch)

A change detection server that can notify through various services, written in Go

Some out-of-the-box highlights:  
- Create watches by connecting filters in a DAG
- A small runtime footprint, a basic instance uses around 20MB of memory
- Supports Lua scripting to filter/modify/reduce your data any way you want
- Send notifications through Discord, Matrix, Slack, Telegram and many more services

# Index
- [Intro](#intro)
- [Run](#run)
  - [Binary](#binary)
  - [Docker](#docker)
    - [Compose Templates](#compose-templates)
- [Config](#config)
  - [Database](#database)
    - [Pruning](#pruning)
  - [Proxy](#proxy)
    - [Proxy Pools](#proxy-pools)
    - [Tor](#tor)
  - [Reverse Proxy](#reverse-proxy)
  - [Browserless](#browserless)
  - [Authentication](#Authentication)
- [Filters](#filters)
  - [Schedule](#schedule)
  - [Get URL](#get-url)
  - [Get URLs](#get-urls)
  - [CSS](#css)
  - [XPath](#xpath)
  - [JSON](#json)
  - [Replace](#replace)
  - [Match](#match)
  - [Substring](#substring)
  - [Contains](#contains)
  - [Store](#store)
  - [Expect](#expect)
  - [Notify](#notify)
  - [Math](#math)
    - [Sum](#sum)
    - [Minimum](#minimum)
    - [Maximum](#maximum)
    - [Average](#average)
    - [Count](#count)
    - [Round](#round)
  - [Condition](#condition)
    - [Different Than Last](#different-than-last)
    - [Lower Than Last](#lower-than-last)
    - [Lowest](#lowest)
    - [Lower Than](#lower-than)
    - [Higher Than Last](#higher-than-last)
    - [Highest](#highest)
    - [Higher Than](#higher-than)
  - [Browserless](#browserless-1)
    - [Browserless Get URL](#browserless-get-url)
    - [Browserless Get URLs](#browserless-get-urls)
    - [Browserless Function](#browserless-function)
    - [Browserless Function On Results](#browserless-function-on-results)
  - [Lua](#lua)
- [Notifiers](#notifiers)
  - [Shoutrrr](#shoutrrr)
  - [Apprise](#apprise)
  - [File](#file)
- [Build/Development](#builddevelopment)
  - [Typescript compilation](#type-script-compilation)
- [Dependencies](#dependencies)

# Intro

GoWatch works through filters, a filter performs operations on the input it recieves.  
Here is an example of a 'Watch' that calculates the lowest and average price of 4090s on NewEgg and notifies the user if the lowest price changed:  
![NewEgg 4090](docs/images/newegg_4090.png)  

Note that everything, including scheduling/storing/notifying, is a `filter`.  

`Schedule` is a [cron](#cron) filter with a '@every 15m' value, this will run every 15 minutes.  

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

Download the binary for your platform from the [releases page](https://github.com/broodjeaap/go-watch/releases), for example for Linux:  
`wget https://github.com/broodjeaap/go-watch/releases/download/1.0/go-watch-1.0-linux-amd64 -O ./gowatch`

And make it executable:  
`chmod +x ./gowatch`

Download the config template:  
`wget https://raw.githubusercontent.com/broodjeaap/go-watch/master/config.tmpl -O ./config.yaml`

Or use the binary to generate it:  
```bash
./gowatch -printConfig 2> config.yaml
# or 
./gowatch -writeConfig config.yaml
```

And modify it to fit your needs, then simply run:  
`./gowatch`

## Docker

Probably the easiest way to get started is with the prebuilt docker image `ghcr.io/broodjeaap/go-watch:latest`, first get a config template:  
`docker run --rm ghcr.io/broodjeaap/go-watch:latest -printConfig 2> config.yaml`  

Or:  
`docker run --rm -v $PWD:/config ghcr.io/broodjeaap/go-watch:latest -writeConfig /config/config.yaml`

After modifying the config to fit your needs, start the docker container
```bash
docker run \
    -p 8080:8080 \
    -v $PWD/:/config \
    ghcr.io/broodjeaap/go-watch:latest
```
### Compose templates

There are a few docker-compose templates in the [docs/compose](https://github.com/broodjeaap/go-watch/tree/master/docs/compose) directory that can be downloaded and used as starting points.  
For example, if you want to set up GoWatch with Browserless, Apprise and a PostgreSQL database backend:  
`wget https://raw.githubusercontent.com/broodjeaap/go-watch/master/docs/compose/apprise-browserless-postgresql.yml -O ./docker-compose.yml`

# Config
## Database

By default, GoWatch will use an SQLite database, stored in the `/config` directory for the docker image.  
If you have only a few watches with schedules of minutes+ then SQLite is probably fine.
But with more watches, especially with shorter schedules, Gorm will start logging warnings about `SLOW SQL`.  
Which are just warnings, but at that point it's probably better to switch to another database.

You can use another database by changing the `database.dsn` value in the config or `GOWATCH_DATABASE_DSN` environment variable, for example with a PostgreSQL database:  
```yaml
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

### Pruning

An automatic database prune job that removes repeating values, this can be scheduled by adding a cron schedule to the config or with the `GOWATCH_SCHEDULE_DELAY` environment variable:  
```yaml
database:
  dsn: "/config/watch.db"
  prune: "@every 1h"
```


## Startup CronJob delay

If there are multiple watches set up with the same schedule then if GoWatch is restarted, all these watches will trigger at the same time, which causes a short burst of activity.  
It might be preferable to spread out these schedules a bit, this can be done by setting `schedule.delay` in the config or with the `GOWATCH_SCHEDULE_DELAY` environment variable:  
```yaml
schedule:
  delay: "5s"
```

## Proxy

An HTTP/HTTPS proxy can be configured in the config or through the `GOWATCH_PROXY_URL` environment variable:  
```yaml
proxy:
  url: http://proxy.com:1234
```
This will not work automatically for requests made through Lua filters, but when using the docker image, the `HTTP_PROXY` and `HTTPS_PROXY` environment variables can also be used which will route all traffic through the proxy:  
```yaml
version: "3"

services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    container_name: go-watch
    environment:
    - HTTP_PROXY=http://proxy.com:1234
    - HTTPS_PROXY=http://proxy.com:1234
```
### Proxy pools

Proxy 'pools' can be created by configuring the proxy that GoWatch points to, for example with [Squid](http://www.squid-cache.org/):  
```yaml
version: "3"

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

And in the `squid.conf` the proxy pool would be defined with [cache_peer](http://www.squid-cache.org/Doc/config/cache_peer/)s like this:  
```conf
cache_peer proxy1.com parent 3128 0 round-robin no-query
cache_peer proxy2.com parent 3128 0 round-robin no-query login=user:pass
```

An example `squid.conf` can be found in [docs/proxy/squid-1.conf](docs/proxy/squid-1.conf).

### Tor

[Tor](https://www.torproject.org/) can also be used to proxy your requests, for example with the [tor-privoxy](https://github.com/dockage/tor-privoxy) container:  
```yaml
version: "3"

services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    environment:
    - HTTP_PROXY=http://tor-privoxy:8118
    - HTTPS_PROXY=http://tor-privoxy:8118
    volumes:
    - ./tmp:/config
    ports:
    - "8080:8080"
  tor-privoxy:
    image: dockage/tor-privoxy
```

To test if it's working, add a [Get URL](#get-url) filter with a `https://check.torproject.org/api/ip` value, and check the result.

## Reverse Proxy

GoWatch can be run behind a reverse proxy, if it's hosted under a subdomain (https://gowatch.domain.tld), no changes to the config are needed.  
But if you want to run GoWatch under a path (https://domain.tld/gowatch), you can set the `gin.urlprefix` value in the config or the `GOWATCH_GIN_URLPREFIX` environment variable can be used.  
```yaml
gin:
  urlprefix: "/gowatch"
```

## Browserless

Some websites (Amazon for example) don't send all content on the first request, it's added later through javascript.  
To still be able to watch products from these types of websites, GoWatch supports [Browserless](https://www.browserless.io/), the Browserless URL can be added to the config:  
```yaml
browserless:
  url: http://your.browserless:3000
```

Or as an environment variable, for example in a docker-compose:  
```yaml
version: "3"

services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    container_name: go-watch
    environment:
    - GOWATCH_BROWSERLESS_URL=http://browserless:3000
    volumes:
    - /host/path/to/config:/config
    ports:
    - "8080:8080"
  browserless:
    image: browserless/chrome:latest
```

To use Browserless, the [Browserless Get URL](#browserless-get-url), [Browserless Get URLs](#browserless-get-urls), [Browserless Function](#browserless-function) or [Browserless Function on result](#browserless-function-on-results) filters must be used. 

Note that for Browserless request to be proxied, Browserless needs to be configured to do so:  
```yaml
version: "3"

services:
  app:
    image: ghcr.io/broodjeaap/go-watch:latest
    container_name: go-watch
    environment:
    - GOWATCH_PROXY_URL=http://tor-privoxy:8118 
    - GOWATCH_BROWSERLESS_URL=http://browserless:3000
    volumes:
    - /host/path/to/config:/config
    ports:
    - "8080:8080"
  tor-privoxy:
    image: dockage/tor-privoxy
  browserless:
    image: browserless/chrome:latest
    environment:
    - DEFAULT_LAUNCH_ARGS=["--proxy-server=socks5://tor-privoxy:9050"]
```

## Authentication

GoWatch doesn't have built in authentication, but we can use a reverse proxy for that, for example through Traefik:  
```yaml
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

GoWatch comes with many filters that cover most typical use cases.  

## Schedule

The `Schedule` filter is used to schedule when your watch will run.  
It uses the [cron](https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.0#section-readme) package to schedule Go routines, some common examples would be:  
- `@every 15m`: will trigger every 15 minutes starting on server start.
- `@hourly`: will trigger on the hour.
- `30 * * * *`: will trigger every hour on the half hour.  

More detailed instructions can be found in its documentation.

Optionally one or more 'jitter' duration strings can be added:  
- `@every 15m + 10m`: Will trigger every 15 to 25 minutes
- `@every 15m + 5m + 5m`: Same as above, but more centered around 20 minutes

## Get URL

Fetches the given URL and outputs the HTTP response.  
For more complicated requests, POSTing/headers/login, use the [HTTP functionality](https://github.com/vadv/gopher-lua-libs/tree/master/http#client-1) in the Lua filter (snippets for these requests are availble from the web UI).  
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
- `product.price`
- `items.3`
- `products.#.price`

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

Inputs pass if they contain the given regex.

## Store

Stores each input value in the database under its own name.  
It's recommended to do this after reducing inputs to a single value (Minimum/Maximum/Average/etc).

## Expect

Outputs a value when it has no inputs, useful to do something (notify) when something goes wrong with your Watch.  
Will only trigger once and can be set to wait multiple times before triggering.

## Notify

Executes the given template and sends the resulting string as a message to the given notifier(s).  
It uses the [Golang templating language](https://pkg.go.dev/text/template), filters are available by their name, so for the filter named `Min` in the intro:  
- `{{ .Min }}` gets the results (Multiple values get joined by `, `)
- `{{ .Min_Type }}` gets the type of the filter
- `{{ .Min_Var1 }}` gets the first variable, useful for Get URL filters or Schedule filters
- `{{ .Min_Var2 }}` gets the second variable

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
Outputs the inputs rounded to the given decimals, nonnumerical values are skipped.
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

## Browserless


### Browserless Get URL

Fetches the given URL through [Browserless](#browserless) and outputs the HTTP response.  
Will log an error if no Browserless instance is configured.

### Browserless Get URLs

Fetches every URL given as input through [Browserless](#browserless) and outputs every HTTP response.  
Will log an error if no Browserless instance is configured.

### Browserless Function

Executes the given [Puppeteer](https://github.com/puppeteer/puppeteer) [function](https://www.browserless.io/docs/function) in a Browserless session.  

### Browserless Function On Results
Executes the given [Puppeteer](https://github.com/puppeteer/puppeteer) [function](https://www.browserless.io/docs/function) in a Browserless session for every result.  

## Lua

The Lua filter wraps [gopher-lua](https://github.com/yuin/gopher-lua), with [gopher-lua-libs](https://github.com/vadv/gopher-lua-libs) to greatly extend the capabilities of the Lua VM.  
A basic script that just passes all inputs to the output looks like this:  
```lua
for i,input in pairs(inputs) do
	table.insert(outputs, input)
end
```

Both `inputs` and `outputs` are convenience tables provided by GoWatch to make Lua scripting a bit easier.
There is also a `logs` table that can be used the same way as the `outputs` table (`table.insert(logs, 'this will be logged')`) to provide some basic logging.  

Much of the functionality that is provided through individual filters in GoWatch can also be done from Lua.  
The gopher-lua-libs provide an [http](https://github.com/vadv/gopher-lua-libs/tree/master/http) lib, whose output can be parsed with the [xmlpath](https://github.com/vadv/gopher-lua-libs/tree/master/xmlpath) or [json](https://github.com/vadv/gopher-lua-libs/tree/master/json) libs and then filtered with a [regular expression](https://github.com/vadv/gopher-lua-libs/tree/master/regexp) or some regular Lua scripting to then finally be turned into a ready to send notification through a [template](https://github.com/vadv/gopher-lua-libs/tree/master/template).  

# Notifiers

## Shoutrrr

[Shoutrrr](https://containrrr.dev/shoutrrr/v0.5/) can be used to notify many different services, check their docs for a [list](https://containrrr.dev/shoutrrr/v0.5/services/overview/) of which ones.  
An example config for sending notifications through Shoutrrr:   
```yaml
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

[Apprise](https://github.com/caronc/apprise) is another option to send notifications, it supports many different services/protocols, but it requires access to an [Apprise API](https://github.com/caronc/apprise-api).  
Luckily there is a [docker image](https://hub.docker.com/r/caronc/apprise) available that we can add to our compose:  
```yaml
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
```yaml
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

## File

GoWatch can also simply append your notification text to a file:
```yaml
notifiers:
  File:
    type: "file"
    path: /config/notifications.log
```

# Build/Development

For local development, clone this repository:  
`git clone https://github.com/broodjeaap/go-watch`

And build the binary:  
`go build -o ./gowatch`

Or:  
`go run .`

Or if you have [Air](https://github.com/cosmtrek/air) set up, just:  
`air`

## type script compilation

`tsc static/*.ts --lib es2020,dom --target ES5 --watch --downlevelIteration`

# Dependencies

The following libaries are used in Go-Watch:  
- [Gin](https://github.com/gin-gonic/gin) for HTTP server
    - [multitemplate](https://github.com/gin-contrib/multitemplate) for template inheritance
- [Cascadia](https://pkg.go.dev/github.com/andybalholm/cascadia) for CSS selectors
- [htmlquery](https://pkg.go.dev/github.com/antchfx/htmlquery) for XPath selectors
- [validator](https://pkg.go.dev/github.com/go-playground/validator/v10@v10.11.0) for user user input validation
- [Shoutrrr](https://github.com/containrrr/shoutrrr/) for built in notifiers
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