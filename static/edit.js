var __values = (this && this.__values) || function(o) {
    var s = typeof Symbol === "function" && Symbol.iterator, m = s && o[s], i = 0;
    if (m) return m.call(o);
    if (o && typeof o.length === "number") return {
        next: function () {
            if (o && i >= o.length) o = void 0;
            return { value: o && o[i++], done: !o };
        }
    };
    throw new TypeError(s ? "Object is not iterable." : "Symbol.iterator is not defined.");
};
var __read = (this && this.__read) || function (o, n) {
    var m = typeof Symbol === "function" && o[Symbol.iterator];
    if (!m) return o;
    var i = m.call(o), r, ar = [], e;
    try {
        while ((n === void 0 || n-- > 0) && !(r = i.next()).done) ar.push(r.value);
    }
    catch (error) { e = { error: error }; }
    finally {
        try {
            if (r && !r.done && (m = i["return"])) m.call(i);
        }
        finally { if (e) throw e.error; }
    }
    return ar;
};
var __spread = (this && this.__spread) || function () {
    for (var ar = [], i = 0; i < arguments.length; i++) ar = ar.concat(__read(arguments[i]));
    return ar;
};
function onTypeChange(node) {
    if (node === void 0) { node = null; }
    var select = document.getElementById("typeInput");
    var type = select.value;
    var var1Div = document.getElementById("var1Div");
    var1Div.innerHTML = "";
    var var1Label = document.getElementById("var1Label");
    var var2Div = document.getElementById("var2Div");
    var2Div.innerHTML = "";
    var var2Label = document.getElementById("var2Label");
    var var3Div = document.getElementById("var3Div");
    var3Div.innerHTML = "";
    var var3Label = document.getElementById("var3Label");
    var var1Value = "";
    var var2Value = "";
    var var3Value = "";
    if (node != null) {
        // @ts-ignore
        var1Value = node.meta.var1;
        // @ts-ignore
        var2Value = node.meta.var2;
        // @ts-ignore
        var3Value = node.meta.var3;
    }
    switch (type) {
        case "gurl": {
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "URL";
            var1Input.placeholder = "https://shopping.website.com";
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Input.placeholder = "";
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Input.placeholder = "";
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "gurls": {
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "-";
            var1Input.placeholder = "From parents";
            var1Input.disabled = true;
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Input.placeholder = "";
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Input.placeholder = "";
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "xpath": {
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "XPath";
            var1Input.placeholder = "//a[@class='price]";
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Input.placeholder = "";
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Input.placeholder = "";
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "json": {
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "JSON";
            var1Input.placeholder = "products.#.price";
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
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
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "Selector";
            var1Input.placeholder = ".price";
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
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
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "Regex";
            var1Input.placeholder = "So[mM]e(thing|where)";
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = false;
            var2Label.innerHTML = "With";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
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
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "Regex";
            var1Input.placeholder = "So[mM]e(thing|where)";
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
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
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "Substring";
            var1Input.placeholder = ":20,25-40,45,47,49,-20:";
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
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
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "Substring";
            var1Input.placeholder = "some";
            var1Div.appendChild(var1Input);
            var notSelect = document.createElement("select");
            notSelect.name = "var2";
            notSelect.id = "var2Input";
            notSelect.classList.add("form-control");
            var no = document.createElement("option");
            no.value = "false";
            no.innerHTML = "No";
            notSelect.appendChild(no);
            var yes = document.createElement("option");
            yes.value = "true";
            yes.innerHTML = "Yes";
            notSelect.appendChild(yes);
            var2Div.appendChild(notSelect);
            var2Label.innerHTML = "Invert";
            if (var2Value == "true" || var2Value == "false") {
                notSelect.value = var2Value;
            }
            else {
                notSelect.value = "false";
            }
            var var3Input = document.createElement("input");
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
            var mathSelect = document.createElement("select");
            mathSelect.name = "var1";
            mathSelect.id = "var1Input";
            mathSelect.classList.add("form-control");
            var mathOptionSum = document.createElement("option");
            mathOptionSum.value = "sum";
            mathOptionSum.innerHTML = "Sum";
            mathSelect.appendChild(mathOptionSum);
            var mathOptionMin = document.createElement("option");
            mathOptionMin.value = "min";
            mathOptionMin.innerHTML = "Min";
            mathSelect.appendChild(mathOptionMin);
            var mathOptionMax = document.createElement("option");
            mathOptionMax.value = "max";
            mathOptionMax.innerHTML = "Max";
            mathSelect.appendChild(mathOptionMax);
            var mathOptionAvg = document.createElement("option");
            mathOptionAvg.value = "average";
            mathOptionAvg.innerHTML = "Average";
            mathSelect.appendChild(mathOptionAvg);
            var mathOptionCount = document.createElement("option");
            mathOptionCount.value = "count";
            mathOptionCount.innerHTML = "Count";
            mathSelect.appendChild(mathOptionCount);
            var mathOptionRound = document.createElement("option");
            mathOptionRound.value = "round";
            mathOptionRound.innerHTML = "Round";
            mathSelect.appendChild(mathOptionRound);
            var1Label.innerHTML = "Function";
            var1Div.appendChild(mathSelect);
            if (var1Value == "") {
                mathSelect.value = "min";
            }
            else {
                mathSelect.value = var1Value;
            }
            mathSelect.onchange = function () { onMathChange(node); };
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            if (mathSelect.value == "round") {
                var2Input.disabled = false;
                var2Label.innerHTML = "Decimals";
            }
            else {
                var2Input.placeholder = "";
                var2Input.disabled = true;
                var2Label.innerHTML = "-";
            }
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Input.placeholder = "";
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "store": {
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Input.disabled = true;
            var1Label.innerHTML = "-";
            var1Input.placeholder = "";
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "condition": {
            var conditionSelect = document.createElement("select");
            conditionSelect.name = "var1";
            conditionSelect.id = "var1Input";
            conditionSelect.classList.add("form-control");
            var differentThenLast = document.createElement("option");
            differentThenLast.value = "diff";
            differentThenLast.innerHTML = "Different Then Last";
            conditionSelect.appendChild(differentThenLast);
            var lowerThenLast = document.createElement("option");
            lowerThenLast.value = "lowerl";
            lowerThenLast.innerHTML = "Lower Then Last";
            conditionSelect.appendChild(lowerThenLast);
            var lowestEver = document.createElement("option");
            lowestEver.value = "lowest";
            lowestEver.innerHTML = "Lowest Ever";
            conditionSelect.appendChild(lowestEver);
            var lowerThan = document.createElement("option");
            lowerThan.value = "lowert";
            lowerThan.innerHTML = "Lower Than";
            conditionSelect.appendChild(lowerThan);
            var higherThenLast = document.createElement("option");
            higherThenLast.value = "higherl";
            higherThenLast.innerHTML = "Higher Then Last";
            conditionSelect.appendChild(higherThenLast);
            var highestEver = document.createElement("option");
            highestEver.value = "highest";
            highestEver.innerHTML = "Highest Ever";
            conditionSelect.appendChild(highestEver);
            var higherThan = document.createElement("option");
            higherThan.value = "highert";
            higherThan.innerHTML = "Higher Than";
            conditionSelect.appendChild(higherThan);
            conditionSelect.onchange = function () { onConditionChange(); };
            var1Div.appendChild(conditionSelect);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
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
        case "notify": {
            var var1Input = document.createElement("textarea");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "Template";
            var1Input.placeholder = "{{ .WatchName }} new lowest price: {{ .Price }}!";
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "cron": {
            var var1Input = document.createElement("input");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "CRON";
            var1Input.placeholder = "30 3-6,20-23 * * *";
            var1Div.appendChild(var1Input);
            var enabledSelect = document.createElement("select");
            enabledSelect.name = "var2";
            enabledSelect.id = "var2Input";
            enabledSelect.classList.add("form-control");
            var enabledOption = document.createElement("option");
            enabledOption.value = "yes";
            enabledOption.innerHTML = "Enabled";
            enabledSelect.appendChild(enabledOption);
            var disabledOption = document.createElement("option");
            disabledOption.value = "no";
            disabledOption.innerHTML = "Disabled";
            enabledSelect.appendChild(disabledOption);
            if (var2Value == "") {
                enabledSelect.value = "yes";
            }
            else {
                enabledSelect.value = var2Value;
            }
            var2Div.appendChild(enabledSelect);
            var2Label.innerHTML = "Enabled";
            var var3Input = document.createElement("input");
            var3Input.name = "var3";
            var3Input.id = "var3Input";
            var3Input.value = var3Value;
            var3Input.classList.add("form-control");
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            var3Div.appendChild(var3Input);
            break;
        }
        case "lua": {
            var var1Input = document.createElement("textarea");
            var1Input.name = "var1";
            var1Input.id = "var1Input";
            var1Input.value = var1Value;
            var1Input.classList.add("form-control");
            var1Label.innerHTML = "Template";
            var1Input.placeholder = "for i,input in pairs(inputs) do\n\ttable.insert(outputs, input)\nend";
            if (var1Value == "") {
                var1Input.value = "for i,input in pairs(inputs) do\n\ttable.insert(outputs, input)\nend";
            }
            var1Input.rows = 10;
            var1Div.appendChild(var1Input);
            var var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control");
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var2Div.appendChild(var2Input);
            var var3Input = document.createElement("input");
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
function onMathChange(node) {
    if (node === void 0) { node = null; }
    var var1Input = document.getElementById("var1Input");
    var var1Label = document.getElementById("var1Label");
    var var2Input = document.getElementById("var2Input");
    var var2Label = document.getElementById("var2Label");
    var var3Input = document.getElementById("var3Input");
    var var3Label = document.getElementById("var3Label");
    var var2Value = "";
    var var3Value = "";
    if (node != null) {
        // @ts-ignore
        var2Value = node.meta.var2;
        // @ts-ignore
        var3Value = node.meta.var3;
    }
    if (var1Input.value == "round") {
        var2Input.disabled = false;
        var2Input.type = "number";
        var2Input.value = var2Value;
        var2Label.innerHTML = "Decimals";
    }
    else {
        var2Input.disabled = true;
        var2Input.type = "text";
        var2Input.value = "";
        var2Label.innerHTML = "-";
    }
}
function onConditionChange(node) {
    var e_1, _a;
    if (node === void 0) { node = null; }
    var var1Input = document.getElementById("var1Input");
    var var1Label = document.getElementById("var1Label");
    var var1Div = document.getElementById("var1Div");
    var var2Input = document.getElementById("var2Input");
    var var2Label = document.getElementById("var2Label");
    var var2Div = document.getElementById("var2Div");
    var2Div.innerHTML = "";
    var var3Input = document.getElementById("var3Input");
    var var3Label = document.getElementById("var3Label");
    var var3Div = document.getElementById("var3Div");
    var var1Value = "";
    var var2Value = "";
    var var3Value = "";
    if (node != null) {
        // @ts-ignore
        var1Value = node.meta.var1;
        var1Input.value = var1Value;
        // @ts-ignore
        var2Value = node.meta.var2;
        // @ts-ignore
        var3Value = node.meta.var3;
    }
    else {
        var1Value = var1Input.value;
    }
    var1Label.innerHTML = "Condition";
    switch (var1Value) {
        case "lowert": {
            var var2Input_1 = document.createElement("input");
            var2Input_1.name = "var2";
            var2Input_1.id = "var2Input";
            var2Input_1.type = "number";
            var2Input_1.value = var2Value;
            var2Input_1.classList.add("form-control");
            var2Label.innerHTML = "Threshold";
            var2Div.appendChild(var2Input_1);
            break;
        }
        case "highert": {
            var var2Input_2 = document.createElement("input");
            var2Input_2.name = "var2";
            var2Input_2.id = "var2Input";
            var2Input_2.type = "number";
            var2Input_2.value = var2Value;
            var2Input_2.classList.add("form-control");
            var2Label.innerHTML = "Threshold";
            var2Div.appendChild(var2Input_2);
            break;
        }
        default: {
            var filterSelect = document.createElement("select");
            filterSelect.name = "var2";
            filterSelect.id = "var2Input";
            filterSelect.disabled = false;
            filterSelect.classList.add("form-control");
            try {
                for (var _b = __values(_diagram.nodes.values()), _c = _b.next(); !_c.done; _c = _b.next()) {
                    var node_1 = _c.value;
                    if (node_1.type != "store") {
                        continue;
                    }
                    var nodeOption = document.createElement("option");
                    nodeOption.value = node_1.label;
                    nodeOption.innerHTML = node_1.label;
                    filterSelect.appendChild(nodeOption);
                }
            }
            catch (e_1_1) { e_1 = { error: e_1_1 }; }
            finally {
                try {
                    if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
                }
                finally { if (e_1) throw e_1.error; }
            }
            if (var2Value != "") {
                filterSelect.value = var2Value;
            }
            var2Label.innerHTML = "Filter";
            var2Div.appendChild(filterSelect);
            break;
        }
    }
}
function onSubmitNewFilter() {
    var nameInput = document.getElementById("nameInput");
    var name = nameInput.value;
    var selectType = document.getElementById("typeInput");
    var type = selectType.value;
    var var1Input = document.getElementById("var1Input");
    var var2Input = document.getElementById("var2Input");
    var var3Input = document.getElementById("var3Input");
    _diagram.addNode(Math.max.apply(Math, __spread(Array.from(_diagram.nodes.values()).map(function (n) { return n.id; }))) + 1, _diagram.canvas.width / 2 - _diagram.mouseState.offset.x, _diagram.canvas.height / 2 - _diagram.mouseState.offset.y, name, {
        type: type,
        var1: var1Input.value,
        var2: var2Input.value,
        var3: var3Input.value
    });
}
function editNode(node) {
    var e_2, _a;
    var addFilterButton = document.getElementById("filterButton");
    addFilterButton.click();
    var name = node.label;
    // @ts-ignore
    var type = node.meta.type;
    // @ts-ignore
    var var1 = node.meta.var1;
    // @ts-ignore
    var var2 = node.meta.var2;
    if (var2 === undefined) {
        var2 = "";
    }
    // @ts-ignore
    var var3 = node.meta.var3;
    if (var3 === undefined) {
        var3 = "";
    }
    var nameInput = document.getElementById("nameInput");
    nameInput.value = name;
    var selectType = document.getElementById("typeInput");
    selectType.value = type;
    onTypeChange(node);
    var submitButton = document.getElementById("submitFilterButton");
    submitButton.innerHTML = "Save";
    var filterModalFooter = document.getElementById("filterResultsDiv");
    filterModalFooter.innerHTML = "";
    try {
        for (var _b = __values(node.results), _c = _b.next(); !_c.done; _c = _b.next()) {
            var result = _c.value;
            var cardDiv = document.createElement("div");
            cardDiv.classList.add("card", "my-2");
            var cardBody = document.createElement("div");
            cardBody.classList.add("card-body", "text-center");
            var pre = document.createElement("pre");
            var code = document.createElement("code");
            if (result.length > 1500) {
                code.innerHTML = "String of length >1500";
            }
            else {
                result = result.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
                code.innerHTML = "'" + result + "'";
            }
            cardDiv.appendChild(cardBody);
            pre.appendChild(code);
            cardBody.appendChild(pre);
            filterModalFooter.appendChild(cardDiv);
        }
    }
    catch (e_2_1) { e_2 = { error: e_2_1 }; }
    finally {
        try {
            if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
        }
        finally { if (e_2) throw e_2.error; }
    }
    submitButton.onclick = function () { submitEditNode(node); };
}
function logNode(node) {
    var e_3, _a;
    var logButton = document.getElementById("logButton");
    logButton.click();
    var logTitle = document.getElementById("logModalLabel");
    logTitle.innerHTML = node.label;
    var logBody = document.getElementById("logTableBody");
    logBody.innerHTML = "";
    var logTable = document.getElementById("logTable");
    try {
        for (var _b = __values(node.logs), _c = _b.next(); !_c.done; _c = _b.next()) {
            var log = _c.value;
            var row = document.createElement("tr");
            var cell = document.createElement("td");
            var code = document.createElement("code");
            code.innerHTML = log;
            cell.appendChild(code);
            row.appendChild(cell);
            logBody.appendChild(row);
        }
    }
    catch (e_3_1) { e_3 = { error: e_3_1 }; }
    finally {
        try {
            if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
        }
        finally { if (e_3) throw e_3.error; }
    }
}
function deleteNode(node) {
    _diagram.nodes["delete"](node.id);
    for (var i = 0; i < _diagram.connections.length; i++) {
        var connection = _diagram.connections[i];
        var output = connection.output;
        var input = connection.input;
        if (node.id == output.id || node.id == input.id) {
            _diagram.connections.splice(i, 1);
            i--;
        }
    }
}
function submitEditNode(node) {
    var nameInput = document.getElementById("nameInput");
    node.label = nameInput.value;
    var selectType = document.getElementById("typeInput");
    // @ts-ignore
    node.meta.type = selectType.value;
    var var1Input = document.getElementById("var1Input");
    // @ts-ignore
    node.meta.var1 = var1Input.value;
    var var2Input = document.getElementById("var2Input");
    // @ts-ignore
    node.meta.var2 = var2Input.value;
    var var3Input = document.getElementById("var3Input");
    // @ts-ignore
    node.meta.var3 = var3Input.value;
    node.fixType();
    node.resize(_diagram.ctx);
    var saveWatchButton = document.getElementById("saveButtonMain");
    saveWatchButton.classList.remove("btn-primary");
    saveWatchButton.classList.add("btn-danger");
}
function saveWatch() {
    var e_4, _a, e_5, _b;
    var watchIdInput = document.getElementById("watch_id");
    var watchId = Number(watchIdInput.value);
    var filters = new Array();
    try {
        for (var _c = __values(_diagram.nodes.values()), _d = _c.next(); !_d.done; _d = _c.next()) {
            var filter = _d.value;
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
                var3: filter.meta.var3
            });
        }
    }
    catch (e_4_1) { e_4 = { error: e_4_1 }; }
    finally {
        try {
            if (_d && !_d.done && (_a = _c["return"])) _a.call(_c);
        }
        finally { if (e_4) throw e_4.error; }
    }
    var filtersInput = document.getElementById("filtersInput");
    filtersInput.value = JSON.stringify(filters);
    var connections = new Array();
    try {
        for (var _e = __values(_diagram.connections), _f = _e.next(); !_f.done; _f = _e.next()) {
            var connection = _f.value;
            var output = connection.output;
            var input = connection.input;
            connections.push({
                connection_watch_id: watchId,
                filter_output_id: output.id,
                filter_input_id: input.id
            });
        }
    }
    catch (e_5_1) { e_5 = { error: e_5_1 }; }
    finally {
        try {
            if (_f && !_f.done && (_b = _e["return"])) _b.call(_e);
        }
        finally { if (e_5) throw e_5.error; }
    }
    var connectionsInput = document.getElementById("connectionsInput");
    connectionsInput.value = JSON.stringify(connections);
    var saveWatchForm = document.getElementById("saveWatchForm");
    saveWatchForm.submit();
}
function addFilterButtonClicked() {
    var submitButton = document.getElementById("submitFilterButton");
    submitButton.onclick = onSubmitNewFilter;
    submitButton.innerHTML = "Add Filter";
    var filterModalFooter = document.getElementById("filterResultsDiv");
    filterModalFooter.innerHTML = "";
    var var1Input = document.getElementById("typeInput");
    var1Input.value = "xpath";
    onTypeChange();
}
function pageInit() {
    var select = document.getElementById("typeInput");
    select.onchange = function () { onTypeChange(); };
    var addFilterButton = document.getElementById("filterButton");
    addFilterButton.onclick = addFilterButtonClicked;
    var saveButtonModal = document.getElementById("saveButtonModal");
    saveButtonModal.onclick = saveWatch;
    var saveButtonMain = document.getElementById("saveButtonMain");
    saveButtonMain.onclick = saveWatch;
}
document.addEventListener('DOMContentLoaded', pageInit, false);
