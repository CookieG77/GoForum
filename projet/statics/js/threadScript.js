const path = window.location.pathname;
const segments = path.split("/");
const threadName = segments[2];
// TODO: Remove the print statement
console.log("Thread ID:", threadName);

function sendMessage() {
    const titleText = document.getElementById("new_message_title").value;
    const contentText = document.getElementById("new_message_content").value;

    fetch(`/api/thread/${threadName}/sendMessage`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            title: titleText,
            content: contentText
        })
    })
    .then(response => {
        if (response.ok) {
            // TODO: Add a success message
            console.log("Message sent successfully");
        } else {
            // TODO: Add an error message
            console.error("Error sending message");
            console.log(response);
        }
    });
}

function deleteMessage(messageId) {
    // TODO: Make the delete message function
}

function editMessage(messageId, title, content) {
    // TODO: Make the edit message function
}

function upvoteMessage(messageId) {
    // TODO: Make the upvote message function
}

function downvoteMessage(messageId) {
    // TODO: Make the downvote message function
}