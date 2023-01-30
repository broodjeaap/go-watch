function testSubmit() {
    let form = document.getElementById("uploadForm") as HTMLFormElement;
    form.action = "/backup/test";
    form.submit();
}

function restoreSubmit() {
    let form = document.getElementById("uploadForm") as HTMLFormElement;
    form.action = "/backup/restore";
    form.submit();
}

function initUploadSubmit(){
    let testSubmitInput = document.getElementById("testSubmit") as HTMLInputElement;
    testSubmitInput.onclick = testSubmit;
    
    let restoreSubmitInput = document.getElementById("restoreSubmit") as HTMLInputElement;
    restoreSubmitInput.onclick = restoreSubmit;
}

document.addEventListener('DOMContentLoaded', initUploadSubmit, false);