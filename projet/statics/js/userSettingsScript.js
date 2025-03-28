// On page load link buttons to the funcs
document.addEventListener("DOMContentLoaded", function () {
    let open_settings_config_button = document.getElementById("change-settings-button")
    let close_settings_config_button = document.getElementById("close-change-settings-button")

    let current_settings = document.getElementById("current-settings")
    let change_settings = document.getElementById("change-settings")

    function openChangeSettingsMenu() {
        current_settings.classList.add("hidden")
        change_settings.classList.remove("hidden")
    }

    function closeChangeSettingsMenu() {
        current_settings.classList.remove("hidden")
        change_settings.classList.add("hidden")
    }

    open_settings_config_button.addEventListener("click", openChangeSettingsMenu)
    close_settings_config_button.addEventListener("click", closeChangeSettingsMenu)
});