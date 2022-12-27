function clearCache(){
    let confirmed = confirm("Do you want to clear the URL cache?");
    if (!confirmed){
        return // do nothing
    }
    let data = new URLSearchParams();
    fetch("/cache/clear", {
            method: "POST"
    }).then((response) => {
        if(response.ok){
            window.location.reload()
        } else {
            alert("Could not clear cache");
        }
    });
}
function clearCacheButtonInit(){
    let clearCacheButton = document.getElementById("clearCacheButton") as HTMLElement;
    clearCacheButton.onclick = clearCache;
}
document.addEventListener('DOMContentLoaded', clearCacheButtonInit, false);