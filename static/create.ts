function urlOnChange(){
    let urlInput = document.getElementById("url") as HTMLInputElement;
    if (urlInput.value.length > 0){
        let urlInputRadio = document.getElementById("urlRadio") as HTMLInputElement;
        urlInputRadio.checked = true;
    }
}

function fileOnChange(){
    let fileInput = document.getElementById("file") as HTMLInputElement;
    if (fileInput.files !== null){
        let fileInputRadio = document.getElementById("fileRadio") as HTMLInputElement;
        fileInputRadio.checked = true;
    }
}

function initOnChange(){
    let urlInput = document.getElementById("url") as HTMLInputElement;
    urlInput.onchange = urlOnChange;

    let fileInput = document.getElementById("file") as HTMLInputElement;
    fileInput.onchange = fileOnChange;
}

document.addEventListener('DOMContentLoaded', initOnChange, false);