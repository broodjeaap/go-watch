# Go Watch
[![Build Status](https://drone.broodjeaap.net/api/badges/broodjeaap/go-watch/status.svg)](https://drone.broodjeaap.net/broodjeaap/go-watch)

A change detection server that can notify through various services written in Go

# Intro

GoWatch works through filters, a filter performs operations on the input it recieves.  
Here is an example of a 'Watch' that calculates the lowest and average price of 4090s on NewEgg and notifies the user if the lowest price changed:  
![NewEgg 4090](docs/images/newegg_4090.png)  

Note that everything, including scheduling/storing/notifying, is a `filter`.  

`Schedule` is a [cron](#cron) filter with a '@every 15m' value, so this will run every 15 minutes.  

`NewEgg Fetch` is a [Get URL](#get-url) filter with a 'https://www.newegg.com/p/pl?N=100007709&d=4090&isdeptsrh=1&PageSize=96' value, it's output will be the HTTP response.  

`Select Price` is a [CSS](#css) filter with the value '.item-container .item-action strong[class!="item-buying-choices-price"]' value, it's output will be the html elements containing the prices.  
An [XPath](#xpath) filter could also have been used.  

`Sanitize` is a [Replace](#replace) filter, using a regular expression ('[^0-9]') it removes anything that's not a number.  

`Avg` is an [Average](#average) filter, it calculates the average value of its inputs.  

`Min` is a [Minimum](#minimum) filter, it calculates the minimum value of its inputs.  

`Store Avg` and `Store Min` are [Store](#store) filters, they store its input values in the database.  

`Diff` is a [Different Than Last](#different-then-last) filter, only passing on the inputs that are different then the last value stored in the database.  

`Notify` is a [Notify](#notify) filter, if there are any inputs to this filter, it will execute a template and send the result to a user defined 'notifier' (Telegram/Discord/Matrix/etc).

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

Matches the inputs to a value, supports regular expressions.

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

## type script compilation

`tsc static/*.ts --lib es2020,dom --watch --downlevelIteration`

# Dependencies

The following libaries are used in Go-Watch:  
- [Gin](https://github.com/gin-gonic/gin) for HTTP server
    - [multitemplate](https://github.com/gin-contrib/multitemplate) for template inheritance
- [Cascadia](https://pkg.go.dev/github.com/andybalholm/cascadia) for CSS selectors
- [htmlquery](https://pkg.go.dev/github.com/antchfx/htmlquery) for XPath selectors
- [validator](https://pkg.go.dev/github.com/go-playground/validator/v10@v10.11.0) for user user input validation
- [tgbotapi](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5@v5.5.1) for Telegram
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