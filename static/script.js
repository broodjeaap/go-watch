function submitNewUrlForm() {
    var nameFrom = document.getElementById("urlNameFrom");
    var nameTo = document.getElementById("urlNameTo");
    nameTo.value = nameFrom.value;
    var urlFrom = document.getElementById("urlFrom");
    var urlTo = document.getElementById("urlTo");
    urlTo.value = urlFrom.value;
    var newUrlForm = document.getElementById("newUrlForm");
    newUrlForm.submit();
}
function submitNewQueryForm() {
    var queryNameFrom = document.getElementById("queryNameFrom");
    var queryNameTo = document.getElementById("queryNameTo");
    queryNameTo.value = queryNameFrom.value;
    var queryFrom = document.getElementById("queryFrom");
    var queryTo = document.getElementById("queryTo");
    queryTo.value = queryFrom.value;
    var newQueryForm = document.getElementById("newQueryForm");
    newQueryForm.submit();
}
