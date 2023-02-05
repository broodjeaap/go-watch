// @ts-ignore
let urlPrefix = getURLPrefix();
function testSubmit() {
    let form = document.getElementById("uploadForm") as HTMLFormElement;
    form.action = urlPrefix + "backup/test";
    form.submit();
}

function restoreSubmit() {
    let form = document.getElementById("uploadForm") as HTMLFormElement;
    form.action = urlPrefix + "backup/restore";
    form.submit();
}

function initUploadSubmit(){
    let testSubmitInput = document.getElementById("testSubmit") as HTMLInputElement;
    testSubmitInput.onclick = testSubmit;
    
    let restoreSubmitInput = document.getElementById("restoreSubmit") as HTMLInputElement;
    restoreSubmitInput.onclick = restoreSubmit;
}

document.addEventListener('DOMContentLoaded', initUploadSubmit, false);