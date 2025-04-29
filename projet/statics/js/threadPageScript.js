
function initPost(){
    const postBox = document.querySelectorAll('.post-box');
    const threadName = getCurrentThreadName();

    postBox.forEach(post => {
        const downvote = document.querySelector('.downvote');
        const upvote = document.querySelector('.upvote');
        const voteNumber = document.querySelector('.vote-number');

        downvote.addEventListener('click', () => {
            downvoteMessage(getCurrentThreadName(), );
        })
        upvote.addEventListener('click', () => {
            upvoteMessage(getCurrentThreadName(), );
        })

    })
}

function loadMoreMessage(offset, order) {

}

document.addEventListener("DOMContentLoaded", function () {

    const leaveButton = document.getElementById("LeaveThreadButton")
    const joinButton = document.getElementById("JoinThreadButton")
    let offset= 0;
    let order = "asc"; // TODO : search in order dropdown

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

    /* getMessage(getCurrentThreadName(), )
    initPost(); */
});
