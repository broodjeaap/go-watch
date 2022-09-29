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
            var1Input.placeholder = "//a[@class='price";
            var1Div.appendChild(var1Input);

            let var2Input = document.createElement("input");
            var2Input.name = "var2";
            var2Input.id = "var2Input";
            var2Input.value = var2Value;
            var2Input.classList.add("form-control")
            var2Input.disabled = true;
            var2Input.placeholder = "";
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

function onSubmitNewFilter(){
    console.log("TEST")
    let nameInput = document.getElementById("nameInput") as HTMLInputElement;
    let name = nameInput.value;
    let selectType = document.getElementById("typeInput") as HTMLSelectElement;
    let type = selectType.value;
    
    let var1Input = document.getElementById("var1Input") as HTMLInputElement;
    
    let var2Input = document.getElementById("var2Input") as HTMLInputElement;

    let var3Input = document.getElementById("var3Input") as HTMLInputElement;

    _diagram.addNode(
        Math.max(...Array.from(_diagram.nodes.values()).map(n => n.id)) + 1,
        _diagram.canvas.width / 2 - _diagram.cameraX, 
        _diagram.canvas.height / 2 - _diagram.cameraY, 
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
    submitButton.onclick = function() {submitEditNode(node);}
}

function deleteNode(node: DiagramNode){
    _diagram.nodes.delete(node.id)
    for (let i = 0; i < _diagram.connections.length; i++){
        let [output, input] = _diagram.connections[i];
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

    node.resize();
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
    for (let [output, input] of _diagram.connections){
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