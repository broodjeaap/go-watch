function clearCache() {
    var confirmed = confirm("Do you want to clear the URL cache?");
    if (!confirmed) {
        return; // do nothing
    }
    var data = new URLSearchParams();
    fetch("/cache/clear", {
        method: "POST"
    }).then(function (response) {
        if (response.ok) {
            window.location.reload();
        }
        else {
            alert("Could not clear cache");
        }
    });
}
function clearCacheButtonInit() {
    var clearCacheButton = document.getElementById("clearCacheButton");
    clearCacheButton.onclick = clearCache;
}
document.addEventListener('DOMContentLoaded', clearCacheButtonInit, false);
