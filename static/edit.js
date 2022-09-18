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
    _diagram.addNode(Math.max.apply(Math, _diagram.nodes.map(function (n) { return n.id; })) + 1, _diagram.canvas.width / 2, _diagram.canvas.height / 2, name, {
        type: type,
        var1: var1Input.value,
        var2: var2Input.value,
        var3: var3Input.value
    });
}
function newFilterInit() {
    var select = document.getElementById("typeInput");
    select.onchange = onTypeChange;
    var button = document.getElementById("submitFilterButton");
    button.onclick = onSubmitNewFilter;
}
document.addEventListener('DOMContentLoaded', newFilterInit, false);
