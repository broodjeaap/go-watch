function newWatch() {
    var response = prompt("Name of new Watch", "");
    if (response == null || response == "") {
        return; // do nothing
    }
    var data = new URLSearchParams();
    data.append("watch_name", response);
    fetch("/watch/create", {
        method: "POST",
        body: data
    }).then(function (response) {
        if (!response.ok) {
            alert("Could not create watch");
            return;
        }
        window.location.href = response.url;
    });
}
function newWatchLinkInit() {
    var newWatchLink = document.getElementById("newWatchLink");
    newWatchLink.onclick = newWatch;
}
document.addEventListener('DOMContentLoaded', newWatchLinkInit, false);
