function testSubmit() {
    var form = document.getElementById("uploadForm");
    form.action = "/backup/test";
    form.submit();
}
function restoreSubmit() {
    var form = document.getElementById("uploadForm");
    form.action = "/backup/restore";
    form.submit();
}
function initUploadSubmit() {
    var testSubmitInput = document.getElementById("testSubmit");
    testSubmitInput.onclick = testSubmit;
    var restoreSubmitInput = document.getElementById("restoreSubmit");
    restoreSubmitInput.onclick = restoreSubmit;
}
document.addEventListener('DOMContentLoaded', initUploadSubmit, false);
