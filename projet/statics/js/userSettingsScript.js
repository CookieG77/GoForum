// On page load link buttons to the funcs
document.addEventListener("DOMContentLoaded", function () {

    let openSettingsConfigButton = document.getElementById("change-settings-button")
    let closeSettingsConfigButton = document.getElementById("close-change-settings-button")

    let currentSettings = document.getElementById("current-settings")
    let changeSettings = document.getElementById("change-settings")

    let pfpContainer = document.getElementById("pfp-container")
    let changePfpPopupBackground = document.getElementById("change-pfp-popup-bg")
    let changePfpPopup = document.getElementById("change-pfp-popup")

    let newPfpInput = document.getElementById("new-pfp-input")
    let newPfpDropzone = document.getElementById("new-pfp-dropzone")

    function openChangeSettingsMenu() {
        currentSettings.classList.add("hidden")
        changeSettings.classList.remove("hidden")
    }

    function closeChangeSettingsMenu() {
        currentSettings.classList.remove("hidden")
        changeSettings.classList.add("hidden")
    }

    openSettingsConfigButton.addEventListener('click', openChangeSettingsMenu)
    closeSettingsConfigButton.addEventListener('click', closeChangeSettingsMenu)

    pfpContainer.addEventListener('click' , function () {
        changePfpPopupBackground.classList.remove("hidden")
        changePfpPopup.classList.remove("hidden")
    });

    changePfpPopupBackground.addEventListener('click', function () {
        changePfpPopupBackground.classList.add("hidden")
        changePfpPopup.classList.add("hidden")
    });

    newPfpDropzone.addEventListener('click', () => newPfpInput.click());

    newPfpDropzone.addEventListener('drop', function (event) {
        event.preventDefault();
        const file = event.dataTransfer.files[0];
        if (file) {
            UploadPfp(file);
        }
    });

    newPfpInput.addEventListener('change', function () {
        const file = newPfpInput.files[0];
        if (file) {
            UploadPfp(file);
        }
    });

    function UploadPfp(file) {
        UploadImg(file, "pfp").then(r => {
            if (r == null) return;
            if (r.url) {
                document.getElementById("pfp").src = "/upload/" + r.url;
                document.getElementById("user-profile-picture").src = "/upload/" + r.url;
                changePfpPopupBackground.classList.add("hidden")
                changePfpPopup.classList.add("hidden")
            } else {
                alert("Server response error");
            }
        }).catch((error) => {
            console.error("Error uploading image:", error);
            alert("Error uploading image");
        });
    }
});