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
function onTypeChange() {
    var select = document.getElementById("typeInput");
    var type = select.value;
    var var1Input = document.getElementById("var1Input");
    var var1Label = document.getElementById("var1Label");
    var var2Input = document.getElementById("var2Input");
    var var2Label = document.getElementById("var2Label");
    var var3Input = document.getElementById("var3Input");
    var var3Label = document.getElementById("var3Label");
    switch (type) {
        case "xpath": {
            var1Label.innerHTML = "XPath";
            var1Input.placeholder = "//a[@class='price";
            var2Input.disabled = true;
            var2Input.placeholder = "";
            var2Label.innerHTML = "-";
            var3Input.disabled = true;
            var3Input.placeholder = "";
            var3Label.innerHTML = "-";
            break;
        }
        case "json": {
            var1Label.innerHTML = "JSON";
            var1Input.placeholder = "products.#.price";
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            break;
        }
        case "css": {
            var1Label.innerHTML = "Selector";
            var1Input.placeholder = ".price";
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            break;
        }
        case "replace": {
            var1Label.innerHTML = "Regex";
            var1Input.placeholder = "So[mM]e(thing|where)";
            var2Input.disabled = false;
            var2Label.innerHTML = "With";
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            break;
        }
        case "match": {
            var1Label.innerHTML = "Regex";
            var1Input.placeholder = "So[mM]e(thing|where)";
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
            break;
        }
        case "substring": {
            var1Label.innerHTML = "Substring";
            var1Input.placeholder = ":20,25-40,45,47,49,-20:";
            var2Input.disabled = true;
            var2Label.innerHTML = "-";
            var3Input.disabled = true;
            var3Label.innerHTML = "-";
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
    _diagram.addNode(Math.max.apply(Math, __spread(Array.from(_diagram.nodes.values()).map(function (n) { return n.id; }))) + 1, _diagram.canvas.width / 2 - _diagram.cameraX, _diagram.canvas.height / 2 - _diagram.cameraY, name, {
        type: type,
        var1: var1Input.value,
        var2: var2Input.value,
        var3: var3Input.value
    });
}
function editNode(node) {
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
    var var1Input = document.getElementById("var1Input");
    var1Input.value = var1;
    var var2Input = document.getElementById("var2Input");
    var2Input.value = var2;
    var var3Input = document.getElementById("var3Input");
    var3Input.value = var3;
    onTypeChange();
    var submitButton = document.getElementById("submitFilterButton");
    submitButton.innerHTML = "Save";
    submitButton.onclick = function () { submitEditNode(node); };
}
function deleteNode(node) {
    _diagram.nodes["delete"](node.id);
    for (var i = 0; i < _diagram.connections.length; i++) {
        var _a = __read(_diagram.connections[i], 2), output = _a[0], input = _a[1];
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
    node.resize();
}
function saveWatch() {
    var e_1, _a, e_2, _b;
    var watchIdInput = document.getElementById("watch_id");
    var watchId = Number(watchIdInput.value);
    var filters = new Array();
    try {
        for (var _c = __values(_diagram.nodes.values()), _d = _c.next(); !_d.done; _d = _c.next()) {
            var filter = _d.value;
            filters.push({
                WatchID: watchId,
                id: filter.id,
                filter_name: filter.label,
                x: filter.x,
                y: filter.y,
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
    catch (e_1_1) { e_1 = { error: e_1_1 }; }
    finally {
        try {
            if (_d && !_d.done && (_a = _c["return"])) _a.call(_c);
        }
        finally { if (e_1) throw e_1.error; }
    }
    var filtersInput = document.getElementById("filtersInput");
    filtersInput.value = JSON.stringify(filters);
    var connections = new Array();
    try {
        for (var _e = __values(_diagram.connections), _f = _e.next(); !_f.done; _f = _e.next()) {
            var _g = __read(_f.value, 2), output = _g[0], input = _g[1];
            connections.push({
                WatchID: watchId,
                filter_output_id: output.id,
                filter_input_id: input.id
            });
        }
    }
    catch (e_2_1) { e_2 = { error: e_2_1 }; }
    finally {
        try {
            if (_f && !_f.done && (_b = _e["return"])) _b.call(_e);
        }
        finally { if (e_2) throw e_2.error; }
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
}
function pageInit() {
    var select = document.getElementById("typeInput");
    select.onchange = onTypeChange;
    var addFilterButton = document.getElementById("filterButton");
    addFilterButton.onclick = addFilterButtonClicked;
    var saveButton = document.getElementById("saveButton");
    saveButton.onclick = saveWatch;
}
document.addEventListener('DOMContentLoaded', pageInit, false);
