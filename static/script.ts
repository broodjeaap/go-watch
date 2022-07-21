function submitNewUrlForm(){
    let nameFrom = document.getElementById("urlNameFrom") as HTMLInputElement;
    let nameTo = document.getElementById("urlNameTo") as HTMLInputElement;
    nameTo.value = nameFrom.value;

    let urlFrom = document.getElementById("urlFrom") as HTMLInputElement;
    let urlTo = document.getElementById("urlTo") as HTMLInputElement;
    urlTo.value = urlFrom.value;

    let newUrlForm = document.getElementById("newUrlForm") as HTMLFormElement;
    newUrlForm.submit();
}

function submitNewQueryForm() {
    let queryNameFrom = document.getElementById("queryNameFrom") as HTMLInputElement;
    let queryNameTo = document.getElementById("queryNameTo") as HTMLInputElement;
    queryNameTo.value = queryNameFrom.value;
    
    let queryFrom = document.getElementById("queryFrom") as HTMLInputElement;
    let queryTo = document.getElementById("queryTo") as HTMLInputElement;
    queryTo.value = queryFrom.value;

    let newQueryForm = document.getElementById("newQueryForm") as HTMLFormElement;
    newQueryForm.submit();
}