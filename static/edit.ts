function onTypeChange(node: DiagramNode | null = null){
    let select = document.getElementById("typeInput") as HTMLSelectElement;
    let type = select.value;
    
    let var1Div = document.getElementById("var1Div") as HTMLDivElement;
    var1Div.innerHTML = "";
    let var1Label = document.getElementById("var1Label") as HTMLLabelElement;
    
    let var2Div = document.getElementById("var2Div") as HTMLDivElement;
    var2Div.innerHTML = "";
    let var2Label = document.getElementById("var2Label") as HTMLLabelElement;
    
    let var3Div = document.getElementById("var3Div") as HTMLDivElement;
    var3Div.innerHTML = "";
    let var3Label = document.getElementById("var3Label") as HTMLLabelElement;

    let var1Value = "";
    let var2Value = "";
    let var3Value = "";
    if (node != null){
        // @ts-ignore
        var1Value = node.meta.var1;
        // @ts-ignore
        var2Value = node.meta.var2;
        // @ts-ignore
        var3Value = node.meta.var3;
    }
    
    switch(type){
        case "gurl": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "URL";
            var1Input.placeholder = "https://shopping.website.com";
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Input.placeholder = ""
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            
            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Input.placeholder = ""
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "gurls": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "-";
            var1Input.placeholder = "From parents";
            var1Input.disabled = true;
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = true;
            var2Input.placeholder = ""
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Input.placeholder = ""
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "xpath": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "XPath";
            var1Input.placeholder = "//a[@class='price]";
            var1Div.appendChild(var1Input);

            let select = document.createElement("select");
            select.name = "var2";
            select.id = "var2Input";
            select.classList.add("form-control");
            let innerHTML = document.createElement("option");
            innerHTML.value = "inner"
            innerHTML.innerHTML = "innerHTML";
            select.appendChild(innerHTML);
            let attributes = document.createElement("option");
            attributes.value = "attr"
            attributes.innerHTML = "Attributes";
            select.appendChild(attributes);
            let node = document.createElement("option");
            node.value = "node"
            node.innerHTML = "Node";
            select.appendChild(node);
            var2Div.appendChild(select);
            var2Label.innerHTML = "Select";
            if (var2Value == ""){
                select.value = "inner";
            } else {
                select.value = var2Value;
            }

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Input.placeholder = ""
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "json": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "JSON";
            var1Input.placeholder = "products.#.price";
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "css": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "Selector";
            var1Input.placeholder = ".price";
            var1Div.appendChild(var1Input);

            let select = document.createElement("select");
            select.name = "var2";
            select.id = "var2Input";
            select.classList.add("form-control");
            let innerHTML = document.createElement("option");
            innerHTML.value = "inner"
            innerHTML.innerHTML = "innerHTML";
            select.appendChild(innerHTML);
            let attributes = document.createElement("option");
            attributes.value = "attr"
            attributes.innerHTML = "Attributes";
            select.appendChild(attributes);
            let node = document.createElement("option");
            node.value = "node"
            node.innerHTML = "Node";
            select.appendChild(node);
            var2Div.appendChild(select);
            var2Label.innerHTML = "Select";
            if (var2Value == ""){
                select.value = "inner";
            } else {
                select.value = var2Value;
            }

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "replace": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "Regex";
            var1Input.placeholder = "So[mM]e(thing|where)";
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = false;
            var2Label.innerHTML = "With";
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "match": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "Regex";
            var1Input.placeholder = "So[mM]e(thing|where)";
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "substring": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "Substring";
            var1Input.placeholder = ":20,25-40,45,47,49,-20:";
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "contains": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "Substring";
            var1Input.placeholder = "some";
            var1Div.appendChild(var1Input);

            let notSelect = document.createElement("select");
            notSelect.name = "var2";
            notSelect.id = "var2Input";
            notSelect.classList.add("form-control");
            let no = document.createElement("option");
            no.value = "false"
            no.innerHTML = "No";
            notSelect.appendChild(no);
            let yes = document.createElement("option");
            yes.value = "true"
            yes.innerHTML = "Yes";
            notSelect.appendChild(yes);
            var2Div.appendChild(notSelect);
            var2Label.innerHTML = "Invert";
            if (var2Value == "true" || var2Value == "false"){
                notSelect.value = var2Value;
            } else {
                notSelect.value = "false";
            }

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "math": {
            let mathSelect = document.createElement("select");
            mathSelect.name = "var1";
            mathSelect.id = "var1Input";
            mathSelect.classList.add("form-control");
            let mathOptionSum = document.createElement("option");
            mathOptionSum.value = "sum"
            mathOptionSum.innerHTML = "Sum";
            mathSelect.appendChild(mathOptionSum);
            let mathOptionMin = document.createElement("option");
            mathOptionMin.value = "min"
            mathOptionMin.innerHTML = "Min";
            mathSelect.appendChild(mathOptionMin);
            let mathOptionMax = document.createElement("option")
            mathOptionMax.value = "max";
            mathOptionMax.innerHTML = "Max";
            mathSelect.appendChild(mathOptionMax);
            let mathOptionAvg = document.createElement("option")
            mathOptionAvg.value = "average";
            mathOptionAvg.innerHTML = "Average";
            mathSelect.appendChild(mathOptionAvg);
            let mathOptionCount = document.createElement("option")
            mathOptionCount.value = "count";
            mathOptionCount.innerHTML = "Count";
            mathSelect.appendChild(mathOptionCount);
            let mathOptionRound = document.createElement("option")
            mathOptionRound.value = "round";
            mathOptionRound.innerHTML = "Round";
            mathSelect.appendChild(mathOptionRound);
            var1Label.innerHTML = "Function";
            var1Div.appendChild(mathSelect);
            if (var1Value == ""){
                mathSelect.value = "min";
            } else {
                mathSelect.value = var1Value;
            }
            mathSelect.onchange = function() {onMathChange(node)}
            
            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            if (mathSelect.value == "round"){
                var2Input.disabled = false;
                var2Label.innerHTML = "Decimals";
            } else {
                var2Input.placeholder = ""
                var2Input.disabled = true;
                var2Label.innerHTML = "-";
            }
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Input.placeholder = ""
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "store": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Input.disabled = true;
            var1Label.innerHTML = "-";
            var1Input.placeholder = "";
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "unique": {
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Input.disabled = true;
            var1Label.innerHTML = "-";
            var1Input.placeholder = "";
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "condition":{
            let conditionSelect = document.createElement("select");
            conditionSelect.name = "var1";
            conditionSelect.id = "var1Input";
            conditionSelect.classList.add("form-control");
            let differentThenLast = document.createElement("option");
            differentThenLast.value = "diff"
            differentThenLast.innerHTML = "Different Then Last";
            conditionSelect.appendChild(differentThenLast);
            let lowerThenLast = document.createElement("option");
            lowerThenLast.value = "lowerl"
            lowerThenLast.innerHTML = "Lower Then Last";
            conditionSelect.appendChild(lowerThenLast);
            let lowestEver = document.createElement("option");
            lowestEver.value = "lowest"
            lowestEver.innerHTML = "Lowest Ever";
            conditionSelect.appendChild(lowestEver);
            let lowerThan = document.createElement("option");
            lowerThan.value = "lowert"
            lowerThan.innerHTML = "Lower Than";
            conditionSelect.appendChild(lowerThan);
            let higherThenLast = document.createElement("option");
            higherThenLast.value = "higherl"
            higherThenLast.innerHTML = "Higher Then Last";
            conditionSelect.appendChild(higherThenLast);
            let highestEver = document.createElement("option");
            highestEver.value = "highest"
            highestEver.innerHTML = "Highest Ever";
            conditionSelect.appendChild(highestEver);
            let higherThan = document.createElement("option");
            higherThan.value = "highert"
            higherThan.innerHTML = "Higher Than";
            conditionSelect.appendChild(higherThan);
            conditionSelect.onchange = function() {onConditionChange()}
            var1Div.appendChild(conditionSelect);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            onConditionChange(node);
            break;
        }
        case "notify":{
            let var1Input = document.createElement("textarea");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "Template";
            var1Input.placeholder = "{{ .WatchName }} new lowest price: {{ .Price }}!";
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("select");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.classList.add("form-control")
            // @ts-ignore
            for (let notifier of notifiers){
                if (notifier == ""){
                    continue;
                }
                let option = document.createElement("option");
                option.value = notifier;
                option.innerHTML = notifier;
                var2Input.appendChild(option);
            }
            if (var2Value == ""){
                var2Input.value = "All"
            } else {
                var2Input.value = var2Value;
            }
            var2Div.appendChild(var2Input);
            var2Label.innerHTML = "Notify";

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "cron":{
            let var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "CRON";
            var1Input.placeholder = "30 3-6,20-23 * * *";
            var1Div.appendChild(var1Input);

            let enabledSelect = document.createElement("select");
            enabledSelect.name = "var2";
            enabledSelect.id = "var2Input";
            enabledSelect.classList.add("form-control");
            let enabledOption = document.createElement("option");
            enabledOption.value = "yes"
            enabledOption.innerHTML = "Enabled";
            enabledSelect.appendChild(enabledOption);
            let disabledOption = document.createElement("option");
            disabledOption.value = "no"
            disabledOption.innerHTML = "Disabled";
            enabledSelect.appendChild(disabledOption);
            if (var2Value == ""){
                enabledSelect.value = "yes";
            } else {
                enabledSelect.value = var2Value;
            }
            var2Div.appendChild(enabledSelect);
            var2Label.innerHTML = "Enabled"

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "brow":{
            let browserlessSelect = document.createElement("select");
            browserlessSelect.name = "var1";
            browserlessSelect.id = "var1Input";
            browserlessSelect.classList.add("form-control");
            var1Label.innerHTML = "Function";

            let gurlOption = document.createElement("option");
            gurlOption.value = "gurl"
            gurlOption.innerHTML = "Get URL";
            browserlessSelect.appendChild(gurlOption);

            let gurlsOption = document.createElement("option");
            gurlsOption.value = "gurls"
            gurlsOption.innerHTML = "Get URLs";
            browserlessSelect.appendChild(gurlsOption);
            
            let funcOption = document.createElement("option");
            funcOption.value = "func"
            funcOption.innerHTML = "Function";
            browserlessSelect.appendChild(funcOption);
            
            let funcsOption = document.createElement("option");
            funcsOption.value = "funcs"
            funcsOption.innerHTML = "Function on results";
            browserlessSelect.appendChild(funcsOption);
            
            var1Div.appendChild(browserlessSelect);

            if (var1Value == ""){
                browserlessSelect.value = "gurl";
            } else {
                browserlessSelect.value = var1Value;
            }
            browserlessSelect.onchange = function() {onBrowserlessChange()}
            onBrowserlessChange(node)
            break
        }
        case "lua":{
            let var1Input = document.createElement("textarea");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control")
            var1Label.innerHTML = "Template";
            var1Input.placeholder = "for i,input in pairs(inputs) do\n\ttable.insert(outputs, input)\nend";
            if (var1Value == "") {
                var1Input.value = "for i,input in pairs(inputs) do\n\ttable.insert(outputs, input)\nend";
            }
            var1Input.rows = 10;
            var1Div.appendChild(var1Input);

            // dev copy link
            let devCopyA = document.createElement('a');
            let results = node == null ? [] : node.results;
            let luaScript = `inputs = {"${results.join('","')}"}\noutputs = {}\n\n ${var1Input.value}`;
            devCopyA.setAttribute('href', '#')
            devCopyA.classList.add("btn", "btn-primary", "btn-sm");
            devCopyA.innerHTML = "Copy script with inputs";
            devCopyA.onclick = function() {
                if (navigator.clipboard){
                    navigator.clipboard.writeText(luaScript);
                    devCopyA.innerHTML = "Script Copied!";
                } else {
                    alert("Could not copy script, no secure origin?");
                }
            }
            var1Div.appendChild(devCopyA);

            let luaSnippets: Map<string, string> = new Map([
                ["HTTP GET", `local http = require("http")
local client = http.client()

local request = http.request("GET", "https://api.ipify.org")
local result, err = client:do_request(request)
if err then
    table.insert(logs, err)
    error(err)
end
if not (result.code == 200) then
    table.insert(logs, err)
    error(err)
end

table.insert(outputs, result.body)
                `],
                ["HTTP POST", `local http = require("http")
local client = http.client()

local request = http.request("POST", "http://httpbin.org/post", "{}")
local result, err = client:do_request(request)
if err then
    table.insert(logs, err)
    error(err)
end
if not (result.code == 200) then
    table.insert(logs, err)
    error(err)
end

table.insert(outputs, result.body)
                `],
                ["HTTP Auth", `local http = require("http")
local client = http.client()

local request = http.request("GET", "http://httpbin.org/basic-auth/gowatch/gowatch123")
request:set_basic_auth("gowatch", "gowatch123")
local result, err = client:do_request(request)
if err then
    table.insert(logs, err)
    error(err)
end
if not (result.code == 200) then
    table.insert(logs, err)
    error(err)
end

table.insert(outputs, result.body)
                `],
                ["HTTP Bearer", `local http = require("http")
local client = http.client()

local request = http.request("GET", "http://httpbin.org/bearer")
local token = "gowatch123"
request:header_set("Authorization", "Bearer " .. token)
local result, err = client:do_request(request)
if err then
    table.insert(logs, err)
    error(err)
end
if not (result.code == 200) then
    table.insert(logs, err)
    error(err)
end

table.insert(outputs, result.body)
                `],
                ["HTTP User Agent", `local http = require("http")
local agent = "GoWatch Crawler"
local client = http.client({user_agent = agent})

local request = http.request("GET", "http://httpbin.org/headers")
request:header_set("User-Agent", agent)
local result, err = client:do_request(request)
if err then
    table.insert(logs, err)
    error(err)
end
if not (result.code == 200) then
    table.insert(logs, err)
    error(err)
end

table.insert(outputs, result.body)
                `],
                ["HTTP Proxy", `local http = require("http")
local client = http.client({ proxy = "http://login:password@hostname.com" })

local request = http.request("GET", "https://api.ipify.org")
local result, err = client:do_request(request)
if err then
    table.insert(logs, err)
    error(err)
end
if not (result.code == 200) then
    table.insert(logs, err)
    error(err)
end

table.insert(outputs, result.body)
                `],
                ["HTTP Browserless", `local http = require("http")
local client = http.client()

# note " for keys/values
local data = '{"url": "https://api.ipify.org"}'
local request = http.request("POST", "http://browserless:3000/content", data)
request:header_set("Content-Type", "application/json")
local result, err = client:do_request(request)
if err then
    table.insert(logs, err)
    error(err)
end
if not (result.code == 200) then
    table.insert(logs, "Response != 200 - " .. result.code)
end

table.insert(outputs, result.body)        
                `],
                ["XPath", `local xmlpath = require("xmlpath")

local query = "//td[@class='price']"
local query, err = xmlpath.compile(query)
if err then
    table.insert(logs, err)
    error(err)
end

for i,input in pairs(inputs) do
    local node, err = xmlpath.load(input)
    if err then
        table.insert(logs, err)
        error(err)
    end
    local iter = query:iter(node)
    for k, v in pairs(iter) do
        table.insert(outputs, v:string())
    end
end
                `],
                ["strings", `local strings = require("strings")
for i,input in pairs(inputs) do
    table.insert(outputs, strings.trim_space(input))
end
                `],
                ["filter", `for i,input in pairs(inputs) do
    number = tonumber(input)
    if number % 2 == 0 then
        table.insert(outputs, input)
    end
end
                `],
            ]);
            let snippetDiv = document.createElement("div");
            snippetDiv.classList.add("d-flex", "flex-wrap")
            // add snippets
            for (let [name, snippet] of luaSnippets) {
                let link = document.createElement('a');
                link.setAttribute("href", "#");
                link.classList.add("btn", "btn-outline-secondary", "btn-sm", "xs");
                link.innerHTML = name;
                link.onclick = function() {var1Input.value = snippet;};
                snippetDiv.appendChild(link)

                let gap = document.createElement("div");
                gap.classList.add("m-1", "xs");
                snippetDiv.appendChild(gap);
            }
            var2Label.innerHTML = "Snippets";
            var2Div.appendChild(snippetDiv);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);

            let var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
    }
}

function onMathChange(node: DiagramNode | null = null){
    let var1Input = document.getElementById("var1Input") as HTMLSelectElement;
    let var1Label = document.getElementById("var1Label") as HTMLLabelElement;
    let var2Input = document.getElementById("var2Input") as HTMLInputElement;
    let var2Label = document.getElementById("var2Label") as HTMLLabelElement;
    let var3Input = document.getElementById("var3Input") as HTMLInputElement;
    let var3Label = document.getElementById("var3Label") as HTMLLabelElement;

    let var2Value = "";
    let var3Value = "";
    if (node != null){
        // @ts-ignore
        var2Value = node.meta.var2;
        // @ts-ignore
        var3Value = node.meta.var3;
    }

    if (var1Input.value == "round"){
        var2Input.disabled = false;
        var2Input.type = "number";
        var2Input.value = var2Value;
        var2Label.innerHTML = "Decimals";
    } else {
        var2Input.disabled = true;
        var2Input.type = "text";
        var2Input.value = "";
        var2Label.innerHTML = "-";
    }
}

function onConditionChange(node: DiagramNode | null = null){
    let var1Input = document.getElementById("var1Input") as HTMLSelectElement;
    let var1Label = document.getElementById("var1Label") as HTMLLabelElement;
    let var1Div = document.getElementById("var1Div") as HTMLDivElement;
    let var2Input = document.getElementById("var2Input") as HTMLInputElement;
    let var2Label = document.getElementById("var2Label") as HTMLLabelElement;
    let var2Div = document.getElementById("var2Div") as HTMLDivElement;
    var2Div.innerHTML = "";
    let var3Input = document.getElementById("var3Input") as HTMLInputElement;
    let var3Label = document.getElementById("var3Label") as HTMLLabelElement;
    let var3Div = document.getElementById("var3Div") as HTMLDivElement;

    let var1Value = "";
    let var2Value = "";
    let var3Value = "";
    if (node != null){
        // @ts-ignore
        var1Value = node.meta.var1;
        var1Input.value = var1Value;
        // @ts-ignore
        var2Value = node.meta.var2;
        // @ts-ignore
        var3Value = node.meta.var3;
    } else {
        var1Value = var1Input.value;
    }

    var1Label.innerHTML = "Condition";

    switch(var1Value) {
        case "lowert": {
            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.type = "number";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Label.innerHTML = "Threshold";
            var2Div.appendChild(var2Input)
            break;
        }
        case "highert": {
            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.type = "number";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Label.innerHTML = "Threshold";
            var2Div.appendChild(var2Input)
            break;
        }
        default: {
            let filterSelect = document.createElement("select");
            filterSelect.name = "var2";
            filterSelect.id = "var2Input";
            filterSelect.disabled = false;
            filterSelect.classList.add("form-control");
            
            for(let node of _diagram.nodes.values()) {
                if (node.type != "store"){
                    continue;
                }
                let nodeOption = document.createElement("option");
                nodeOption.value = node.label;
                nodeOption.innerHTML = node.label;
                filterSelect.appendChild(nodeOption);
            }
            if (var2Value != ""){
                filterSelect.value = var2Value;
            }
            var2Label.innerHTML = "Filter";
            var2Div.appendChild(filterSelect);
            break
        }
    }
}

function onBrowserlessChange(node: DiagramNode | null = null){
    let var1Input = document.getElementById("var1Input") as HTMLSelectElement;
    let var1Label = document.getElementById("var1Label") as HTMLLabelElement;
    let var1Div = document.getElementById("var1Div") as HTMLDivElement;
    let var2Input = document.getElementById("var2Input") as HTMLInputElement;
    let var2Label = document.getElementById("var2Label") as HTMLLabelElement;
    let var2Div = document.getElementById("var2Div") as HTMLDivElement;
    var2Div.innerHTML = "";
    let var3Input = document.getElementById("var3Input") as HTMLInputElement;
    let var3Label = document.getElementById("var3Label") as HTMLLabelElement;
    let var3Div = document.getElementById("var3Div") as HTMLDivElement;
    var3Div.innerHTML = "";

    let var1Value = "";
    let var2Value = "";
    let var3Value = "";
    if (node != null){
        // @ts-ignore
        var1Value = node.meta.var1;
        var1Input.value = var1Value;
        // @ts-ignore
        var2Value = node.meta.var2;
        // @ts-ignore
        var3Value = node.meta.var3;
    } else {
        var1Value = var1Input.value;
    }

    var1Label.innerHTML = "Function";

    switch(var1Value) {
        case "gurl": {
            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Label.innerHTML = "URL";
            var2Div.appendChild(var2Input);
            break;
        }
        case "gurls": {
            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = "";
            var2Input.disabled = true;
            var2Input.classList.add("form-control")
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            break;
        }
        case "func": {
            let var2Input = document.createElement("textarea");
            var2Input.name = "var2Input";
            var2Input.id = "var2Input";
            var2Input.value = `module.exports = async ({ page, context }) => {
  await page.goto("https://192.168.178.254:8000");

  const data = await page.content();

  return {
    data,
    type: 'text/plain', // 'application/html' 'application/json'
  };
};`;
            var2Input.classList.add("form-control");
            var2Input.rows = 15;
            var2Label.innerHTML = "Code";
            var2Div.appendChild(var2Input);

            if (var2Value != ""){
                var2Input.value = var2Value;
            }

            let var3Input = document.createElement("input");
            var3Input.type = "hidden";
            var3Input.id = "var3Input";
            var3Input.name = "var3Input";
            var3Div.appendChild(var3Input);

            var3Label.innerHTML = "Help";
            let helpLink1 = document.createElement("a");
            helpLink1.href = "https://www.browserless.io/docs/function";
            helpLink1.innerHTML = "Browserless /Funcion";
            helpLink1.target = "_blank";
            helpLink1.classList.add("btn", "btn-outline-secondary", "btn-sm", "xs");
            var3Div.appendChild(helpLink1);

            let helpLink2 = document.createElement("a");
            helpLink2.href = "https://pptr.dev/api/puppeteer.page";
            helpLink2.innerHTML = "Puppeteer Page";
            helpLink2.target = "_blank";
            helpLink2.classList.add("btn", "btn-outline-secondary", "btn-sm", "xs");
            var3Div.appendChild(helpLink2);
            break;
        }
        case "funcs":{
            let var2Input = document.createElement("textarea");
            var2Input.name = "var2Input";
            var2Input.id = "var2Input";
            var2Input.value = `module.exports = async ({ page, context }) => {
  const { result } = context;
  await page.goto(result);

  const data = await page.content();

  return {
    data,
    type: 'text/plain', // 'application/html' 'application/json'
  };
};`;
            var2Input.classList.add("form-control");
            var2Input.rows = 15;
            var2Label.innerHTML = "Code";
            var2Div.appendChild(var2Input);

            if (var2Value != ""){
                var2Input.value = var2Value;
            }

            let var3Input = document.createElement("input");
            var3Input.type = "hidden";
            var3Input.id = "var3Input";
            var3Input.name = "var3Input";
            var3Div.appendChild(var3Input);

            var3Label.innerHTML = "Help";
            let helpLink1 = document.createElement("a");
            helpLink1.href = "https://www.browserless.io/docs/function";
            helpLink1.innerHTML = "Browserless /Funcion";
            helpLink1.target = "_blank";
            helpLink1.classList.add("btn", "btn-outline-secondary", "btn-sm", "xs");
            var3Div.appendChild(helpLink1);

            let helpLink2 = document.createElement("a");
            helpLink2.href = "https://pptr.dev/api/puppeteer.page";
            helpLink2.innerHTML = "Puppeteer Page";
            helpLink2.target = "_blank";
            helpLink2.classList.add("btn", "btn-outline-secondary", "btn-sm", "xs");
            var3Div.appendChild(helpLink2);
            break;
        }
    }
}

function onSubmitNewFilter(){
    let nameInput = document.getElementById("nameInput") as HTMLInputElement;
    let name = nameInput.value;
    let selectType = document.getElementById("typeInput") as HTMLSelectElement;
    let type = selectType.value;
    
    let var1Input = document.getElementById("var1Input") as HTMLInputElement;
    
    let var2Input = document.getElementById("var2Input") as HTMLInputElement;

    let var3Input = document.getElementById("var3Input") as HTMLInputElement;

    _diagram.addNode(
        Math.max(...Array.from(_diagram.nodes.values()).map(n => n.id), 0) + 1,
        _diagram.canvas.width / 2 - _diagram.mouseState.offset.x, 
        _diagram.canvas.height / 2 - _diagram.mouseState.offset.y, 
        name, {
            type: type,
            var1: var1Input.value,
            var2: var2Input.value,
            var3: var3Input.value,
    })
}

function editNode(node: DiagramNode){
    let addFilterButton = document.getElementById("filterButton") as HTMLButtonElement;
    addFilterButton.click();

    let name = node.label;
    // @ts-ignore
    let type = node.meta.type;
    // @ts-ignore
    let var1 = node.meta.var1;
    // @ts-ignore
    let var2 = node.meta.var2;
    if (var2 === undefined){
        var2 = "";
    }
    // @ts-ignore
    let var3 = node.meta.var3;
    if (var3 === undefined){
        var3 = "";
    }

    let nameInput = document.getElementById("nameInput") as HTMLInputElement;
    nameInput.value = name;

    let selectType = document.getElementById("typeInput") as HTMLSelectElement;
    selectType.value = type;

    onTypeChange(node);
    let submitButton = document.getElementById("submitFilterButton") as HTMLButtonElement;
    submitButton.innerHTML = "Save";

    let logHeader = document.getElementById("logHeader") as HTMLHeadElement;
    if (node.logs.length > 0){
        logHeader.innerHTML = "Logs";
    } else {
        logHeader.innerHTML = "";
    }
    let logBody = document.getElementById("logTableBody") as HTMLElement;
    logBody.innerHTML = "";
    for (let log of node.logs){
        let row = document.createElement("tr");
        let cell = document.createElement("td");
        let code = document.createElement("code");
        code.innerHTML = log;
        cell.appendChild(code);
        row.appendChild(cell);
        logBody.appendChild(row);
    }

    let filterModalFooter = document.getElementById("filterResultsDiv") as HTMLDivElement;
    filterModalFooter.innerHTML = "";
    for (let result of node.results){
        let cardDiv = document.createElement("div");
        cardDiv.classList.add("card", "my-2");
        let cardBody = document.createElement("div");
        cardBody.classList.add("card-body", "text-center");
        let pre = document.createElement("pre");
        let code = document.createElement("code");
        if (result.length > 1500){
            code.innerHTML = `String of length >1500`;
        } else {
            result = result.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
            code.innerHTML = `'${result}'`;
        }
        cardDiv.appendChild(cardBody);
        pre.appendChild(code);
        cardBody.appendChild(pre);
        filterModalFooter.appendChild(cardDiv);
    }
    submitButton.onclick = function() {submitEditNode(node);}
}

function deleteNode(node: DiagramNode){
    _diagram.nodes.delete(node.id)
    for (let i = 0; i < _diagram.connections.length; i++){
        let connection = _diagram.connections[i];
        let output = connection.output;
        let input = connection.input;
        if (node.id == output.id || node.id == input.id){
            _diagram.connections.splice(i, 1)
            i--;
        }
    }
}

function submitEditNode(node: DiagramNode){
    let nameInput = document.getElementById("nameInput") as HTMLInputElement;
    node.label = nameInput.value;

    let selectType = document.getElementById("typeInput") as HTMLSelectElement;
    // @ts-ignore
    node.meta.type = selectType.value
    
    let var1Input = document.getElementById("var1Input") as HTMLInputElement;
    // @ts-ignore
    node.meta.var1 = var1Input.value;
    
    let var2Input = document.getElementById("var2Input") as HTMLInputElement;
    // @ts-ignore
    node.meta.var2 = var2Input.value;

    let var3Input = document.getElementById("var3Input") as HTMLInputElement;
    // @ts-ignore
    node.meta.var3 = var3Input.value;

    node.fixType();
    node.resize(_diagram.ctx);
    let saveWatchButton = document.getElementById("saveButtonMain") as HTMLButtonElement;
    saveWatchButton.classList.remove("btn-primary");
    saveWatchButton.classList.add("btn-danger");
}

function saveWatch(){
    let watchIdInput = document.getElementById("watch_id") as HTMLInputElement;
    let watchId = Number(watchIdInput.value);
    let filters = new Array<Object>();
    for (let filter of _diagram.nodes.values()){
        filters.push({
            filter_watch_id: watchId,
            filter_id: filter.id,
            filter_name: filter.label,
            x: Math.round(filter.x),
            y: Math.round(filter.y),
            // @ts-ignore
            filter_type: filter.meta.type,
            // @ts-ignore
            var1: filter.meta.var1,
            // @ts-ignore
            var2: filter.meta.var2,
            // @ts-ignore
            var3: filter.meta.var3,
        })
    }
    let filtersInput = document.getElementById("filtersInput") as HTMLInputElement;
    filtersInput.value = JSON.stringify(filters);

    let connections = new Array<Object>();
    for (let connection of _diagram.connections){
        let output = connection.output;
        let input = connection.input;
        connections.push({
            connection_watch_id: watchId,
            filter_output_id: output.id,
            filter_input_id: input.id,
        })
    }
    let connectionsInput = document.getElementById("connectionsInput") as HTMLInputElement;
    connectionsInput.value = JSON.stringify(connections);

    let saveWatchForm = document.getElementById("saveWatchForm") as HTMLFormElement;
    saveWatchForm.submit();
}

function addFilterButtonClicked(){
    let submitButton = document.getElementById("submitFilterButton") as HTMLButtonElement;
    submitButton.onclick = onSubmitNewFilter
    submitButton.innerHTML = "Add Filter"
    
    let filterModalFooter = document.getElementById("filterResultsDiv") as HTMLDivElement;
    filterModalFooter.innerHTML = "";
    
    let var1Input = document.getElementById("typeInput") as HTMLInputElement;
    var1Input.value = "xpath";
    onTypeChange();
}

function pageInit(){
    let select = document.getElementById("typeInput") as HTMLSelectElement;
    select.onchange = function () {onTypeChange()};

    let addFilterButton = document.getElementById("filterButton") as HTMLButtonElement;
    addFilterButton.onclick = addFilterButtonClicked

    let saveButtonModal = document.getElementById("saveButtonModal") as HTMLButtonElement;
    saveButtonModal.onclick = saveWatch;
    
    let saveButtonMain = document.getElementById("saveButtonMain") as HTMLButtonElement;
    saveButtonMain.onclick = saveWatch;
}

document.addEventListener('DOMContentLoaded', pageInit, false);

function clearCache(){
    let confirmed = confirm("Do you want to clear the URL cache?");
    if (!confirmed){
        return // do nothing
    }
    let data = new URLSearchParams();
    fetch("/cache/clear", {
            method: "POST"
    }).then((response) => {
        if(response.ok){
            window.location.reload()
        } else {
            alert("Could not clear cache");
        }
    });
}
function clearCacheButtonInit(){
    let clearCacheButton = document.getElementById("clearCacheButton") as HTMLElement;
    clearCacheButton.onclick = clearCache;
}
document.addEventListener('DOMContentLoaded', clearCacheButtonInit, false);