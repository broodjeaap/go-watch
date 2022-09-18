function onTypeChange(){
    let select = document.getElementById("typeInput") as HTMLSelectElement;
    let type = select.value;
    
    let var1Input = document.getElementById("var1Input") as HTMLInputElement;
    let var1Label = document.getElementById("var1Label") as HTMLLabelElement;
    
    let var2Input = document.getElementById("var2Input") as HTMLInputElement;
    let var2Label = document.getElementById("var2Label") as HTMLLabelElement;

    let var3Input = document.getElementById("var3Input") as HTMLInputElement;
    let var3Label = document.getElementById("var3Label") as HTMLLabelElement;
    
    switch(type){
        case "xpath": {
            var1Label.innerHTML = "XPath";
            var1Input.placeholder = "//a[@class='price";
            var2Input.disabled = true;
            var2Input.placeholder = ""
            var2Label.innerHTML = "-";
            var3Input.disabled = true;
            var3Input.placeholder = ""
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

function onSubmitNewFilter(){
    let nameInput = document.getElementById("nameInput") as HTMLInputElement;
    let name = nameInput.value;
    let selectType = document.getElementById("typeInput") as HTMLSelectElement;
    let type = selectType.value;
    
    let var1Input = document.getElementById("var1Input") as HTMLInputElement;
    
    let var2Input = document.getElementById("var2Input") as HTMLInputElement;

    let var3Input = document.getElementById("var3Input") as HTMLInputElement;

    _diagram.addNode(
        Math.max(..._diagram.nodes.map(n => n.id)) + 1,
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
    
    let var1Input = document.getElementById("var1Input") as HTMLInputElement;
    var1Input.value = var1;
    
    let var2Input = document.getElementById("var2Input") as HTMLInputElement;
    var2Input.value = var2;

    let var3Input = document.getElementById("var3Input") as HTMLInputElement;
    var3Input.value = var3;

    onTypeChange();
    let submitButton = document.getElementById("submitFilterButton") as HTMLButtonElement;
    submitButton.innerHTML = "Save";
    submitButton.onclick = function() {submitEditNode(node);}
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

function addFilterButtonClicked(){
    let submitButton = document.getElementById("submitFilterButton") as HTMLButtonElement;
    submitButton.onclick = onSubmitNewFilter
    submitButton.innerHTML = "Add Filter"
}

function newFilterInit(){
    let select = document.getElementById("typeInput") as HTMLSelectElement;
    select.onchange = onTypeChange;

    let addFilterButton = document.getElementById("filterButton") as HTMLButtonElement;
    addFilterButton.onclick = addFilterButtonClicked
}
document.addEventListener('DOMContentLoaded', newFilterInit, false);