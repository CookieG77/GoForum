document.addEventListener("DOMContentLoaded", function () {
    const threadSelect = document.getElementById("threadSelect");
    const messageTitle = document.getElementById("messageTitle");
    const messageContent = document.getElementById("messageContent");

    const afterMessageSendOptionContainer = document.getElementById("afterMessageSendOptionContainer");

    const messageThreadContainer = document.getElementById("messageThreadContainer");
    const messageIDContainer = document.getElementById("messageIDContainer");
    const updatedMessageTitle = document.getElementById("updatedMessageTitle");
    const updatedMessageContent = document.getElementById("updatedMessageContent");

    document.getElementById("sendMessageButton").addEventListener("click", function() {
        const threadName = threadSelect.value
        sendMessage(threadName, messageTitle.value, messageContent.value)
            .then(r => {
                if (r.ok) {
                    return r.json();
                } else {
                    throw new Error("Error while sending message");
                }
            })
            .then(data => {
                console.log("Message sent successfully", data);
                // Optionally, you can clear the input fields after sending the message
                updatedMessageTitle.value = messageTitle.value;
                updatedMessageContent.value = messageContent.value;
                messageTitle.value = "";
                messageContent.value = "";

                // We fill the messageIDContainer with the new message ID
                messageThreadContainer.textContent = threadName;
                messageIDContainer.textContent = data.messageId;
                afterMessageSendOptionContainer.classList.remove("hidden");

            })
            .catch(error => {
                    console.error("Error:", error);
                });
    });

    document.getElementById("deleteMessageButton").addEventListener("click", function() {
    const threadName = messageThreadContainer.textContent
        deleteMessage(threadName, messageIDContainer.textContent)
            .then(r => {
                if (r.ok) {
                    return r.json();
                } else {
                    throw new Error("Error while deleting message");
                }
            })
            .then(data => {
                console.log("Message deleted successfully", data);
                // Optionally, you can clear the input fields after deleting the message
                messageThreadContainer.textContent = "";
                messageIDContainer.textContent = "";
                updatedMessageTitle.value = "";
                updatedMessageContent.value = "";
                afterMessageSendOptionContainer.classList.add("hidden");
            })
            .catch(error => {
                    console.error("Error:", error);
                });
    });

    document.getElementById("upvoteButton").addEventListener("click", function () {
        const threadName = messageThreadContainer.textContent
        const messageId = messageIDContainer.textContent
        upvoteMessage(threadName, messageId)
            .then(r => {
                if (r.ok) {
                    return r.json();
                } else {
                    throw new Error("Error while upvoting message");
                }
            })
            .then(data => {
                console.log("Message upvoted successfully", data);
            })
            .catch(error => {
                    console.error("Error:", error);
                });
    });

    document.getElementById("downvoteButton").addEventListener("click", function () {
        const threadName = messageThreadContainer.textContent
        const messageId = messageIDContainer.textContent
        downvoteMessage(threadName, messageId)
            .then(r => {
                if (r.ok) {
                    return r.json();
                } else {
                    throw new Error("Error while downvoting message");
                }
            })
            .then(data => {
                console.log("Message downvoted successfully", data);
            })
            .catch(error => {
                    console.error("Error:", error);
                });
    });

    document.getElementById("updateMessageButton").addEventListener("click", function () {
        const threadName = messageThreadContainer.textContent
        const messageId = messageIDContainer.textContent
        editMessage(threadName, messageId, updatedMessageTitle.value, updatedMessageContent.value)
            .then(r => {
                if (r.ok) {
                    return r.json();
                } else {
                    throw new Error("Error while updating message");
                }
            })
            .then(data => {
                console.log("Message updated successfully", data);
            })
            .catch(error => {
                    console.error("Error:", error);
                });
    });

});