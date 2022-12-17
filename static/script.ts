function newWatch(){
    let response = prompt("Name of new Watch", "");
    if (response == null || response == ""){
        return // do nothing
    }
    let data = new URLSearchParams();
    data.append("watch_name", response);
    fetch("/watch/create", {
            method: "POST",
            body: data,
    }).then((response) => {
        if(!response.ok){
            alert("Could not create watch");
            return;
        }
        window.location.href = response.url;
    });
}

function newWatchLinkInit(){
    let newWatchLink = document.getElementById("newWatchLink") as HTMLElement;
    newWatchLink.onclick = newWatch;
}

document.addEventListener('DOMContentLoaded', newWatchLinkInit, false);