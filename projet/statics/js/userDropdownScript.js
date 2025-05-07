document.addEventListener("DOMContentLoaded", function () {

    /**
     * Allow a dropdown to be shown or hidden when the mouse hovers over the dropdown or its trigger.
     * @param {string} triggerId - The ID of the html element used as a trigger for the function.
     * @param {string} dropdownId - The idea od the dropdown element.
     **/

    function setupDropdown(triggerId, dropdownId) {
        const trigger = document.getElementById(triggerId);
        const dropdown = document.getElementById(dropdownId);
        let timeout;

        /**
         * Display the dropdown element.
         **/
        function showDropdown() {
            clearTimeout(timeout);
            dropdown.classList.add("active");
        }

        /**
         * Hides the dropdown element when the cursor doesn't hover over it or the trigger element.
         **/
        function hideDropdown() {
            // Add a small delay to allow for mouse transitions between trigger and dropdown
            timeout = setTimeout(() => {
                // Only hide if the mouse is not over either element
                if (!trigger.matches(':hover') && !dropdown.matches(':hover')) {
                    dropdown.classList.remove("active");
                }
            }, 300);
        }

        // Attach event listeners to the trigger
        trigger.addEventListener('mouseenter', showDropdown);
        trigger.addEventListener('mouseleave', hideDropdown);

        // Attach event listeners to the dropdown itself
        dropdown.addEventListener('mouseenter', showDropdown);
        dropdown.addEventListener('mouseleave', hideDropdown);
    }

    setupDropdown("user-profile-picture", "user-dropdown");

});