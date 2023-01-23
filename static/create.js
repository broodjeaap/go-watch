function urlOnChange() {
    var urlInput = document.getElementById("url");
    if (urlInput.value.length > 0) {
        var urlInputRadio = document.getElementById("urlRadio");
        urlInputRadio.checked = true;
    }
}
function fileOnChange() {
    var fileInput = document.getElementById("file");
    if (fileInput.files !== null) {
        var fileInputRadio = document.getElementById("fileRadio");
        fileInputRadio.checked = true;
    }
}
function initOnChange() {
    var urlInput = document.getElementById("url");
    urlInput.onchange = urlOnChange;
    var fileInput = document.getElementById("file");
    fileInput.onchange = fileOnChange;
}
document.addEventListener('DOMContentLoaded', initOnChange, false);
