document.addEventListener("DOMContentLoaded", function () {

    const leaveButton = document.getElementById("LeaveThreadButton")
    const joinButton = document.getElementById("JoinThreadButton")

    leaveButton.addEventListener("click", function() {
        const result = leaveThread(getCurrentThreadName());
        result.then(async (response) => {
            if (response.ok) {
                joinButton.classList.remove("hidden");
                leaveButton.classList.add("hidden");
                console.log("You have joined the thread");
            } else {
                console.error(response);
            }
        });
    });

    joinButton.addEventListener("click", function() {
        const result = joinThread(getCurrentThreadName())
        result.then(async (response) => {
            if (response.ok) {
                joinButton.classList.add("hidden");
                leaveButton.classList.remove("hidden");
                console.log("You have left the thread");
            } else {
                console.error(response);
            }
        });
    });

});
