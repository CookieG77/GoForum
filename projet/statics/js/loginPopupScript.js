document.addEventListener("DOMContentLoaded", function () {

    function setupPopup(triggerId, popupId, popupBackgroundId, popupCloseButtonId) {
        const trigger = document.getElementById(triggerId);
        const popup = document.getElementById(popupId);
        const popupBackground = document.getElementById(popupBackgroundId);
        const popupCloseButton = document.getElementById(popupCloseButtonId);

        function showPopup() {
            popup.classList.remove("hidden");
        }

        function hidePopup() {
            popup.classList.add("hidden");
        }

        trigger.addEventListener("click", showPopup);
        popupBackground.addEventListener("click", hidePopup);
        popupCloseButton.addEventListener("click", hidePopup);
    }

    setupPopup("login-popup-opener","login-popup", "login-popup-background", "login-popup-close-button");
    console.log("loginPopupScript.js loaded");
});