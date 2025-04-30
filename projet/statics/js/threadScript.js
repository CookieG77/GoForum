/**
 * Get the current thread name from the URL.
 * @description This function extracts the thread name from the URL path. Only work if you are on the thread page. (e.g. /t/123)
 * @returns {string}
 */
function getCurrentThreadName() {
    const path = window.location.pathname;
    const segments = path.split("/");
    const threadName = segments[2];
    // TODO: Remove the print statement
    console.log("Thread ID:", threadName);
    return threadName;
}

/**
 * Send a message to the current thread.
 * @description This function sends a message to the current thread. It does not handle the response.
 * @description But a success response means that the message has been sent.
 * @param threadName {string} - The name of the thread to send the message to.
 * @param messageTitle {string} - The title of the message.
 * @param messageContent {string} - The content of the message.
 * @param messageMedias {string[]} - The media files to attach to the message.
 * @param messageTags {string[]} - The tags of the message.
 * @returns {Promise<Response>} - The response from the server.
 */
function sendMessage(threadName, messageTitle, messageContent, messageMedias, messageTags) {
    return fetch(`/api/thread/${threadName}/sendMessage`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            title: messageTitle,
            content: messageContent,
            medias: messageMedias,
            tags: messageTags
        })
    });
}

/**
 * Delete a message from the current thread.
 * @description This function sends a request to delete a message from the current thread. It does not handle the response.
 * @description But a success response means that the message has been deleted.
 * @param threadName {string} - The name of the thread to delete the message from.
 * @param messageId {string} - The ID of the message to delete.
 * @returns {Promise<Response>} - The response from the server.
 */
function deleteMessage(threadName, messageId) {
    return fetch( `/api/thread/${threadName}/deleteMessage`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            messageId: messageId
        })
    });
}

/**
 * Remove a media from the current thread.
 * @description This function sends a request to remove a media from the current message. It does not handle the response.
 * @description But a success response means that the media has been removed.
 * @param threadName {string} - The name of the thread to remove the media from.
 * @param messageId {string} - The ID of the message to remove the media from.
 * @param mediaId {string} - The ID of the media to remove.
 * @returns {Promise<Response>} - The response from the server.
 */
function removeMedia(threadName, messageId, mediaId) {
    return fetch( `/api/thread/${threadName}/removeMedia`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            messageId: messageId,
            mediaId: mediaId
        })
    });
}

/**
 * Edit the message with the given id in the given thread .
 * @description This function sends a request to edit a message in the current thread. It does not handle the response.
 * @description But a success response means that the message has been edited.
 * @param threadName {string} - The name of the thread to edit the message in.
 * @param messageId {string} - The ID of the message to edit.
 * @param newMessageTitle {string} - The new title of the message.
 * @param newMessageContent {string} - The new content of the message.
 * @returns {Promise<Response>} - The response from the server.
 */
function editMessage(threadName, messageId, newMessageTitle, newMessageContent) {
    return fetch(`/api/thread/${threadName}/editMessage`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            messageId: messageId,
            title: newMessageTitle,
            content: newMessageContent
        })
    });
}

/**
 * Report the message with the given id in the given thread.
 * @description This function sends a request to report a message in the current thread. It does not handle the response.
 * @description But a success response means that the message has been reported.
 * @param threadName {string} - The name of the thread to report the message in.
 * @param messageId {string} - The ID of the message to report.
 * @returns {Promise<Response>} - The response from the server.
 */
function reportMessage(threadName, messageId) {
    // TODO : Not implemented yet in the backend
    return fetch( `/api/thread/${threadName}/reportMessage`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            messageId: messageId
        })
    });
}

/**
 * Upvote the message with the given id in the given thread.
 * @description This function sends a request to upvote a message in the current thread. It does not handle the response.
 * @description But a success response means that the message has been upvoted.
 * @param threadName {string} - The name of the thread to upvote the message in.
 * @param messageId {string} - The ID of the message to upvote.
 * @returns {Promise<Response>} - The response from the server.
 */
function upvoteMessage(threadName, messageId) {
    return fetch( `/api/thread/${threadName}/upvoteMessage`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            messageId: messageId
        })
    });
}

/**
 * Downvote the message with the given id in the given thread.
 * @description This function sends a request to downvote a message in the current thread. It does not handle the response.
 * @description But a success response means that the message has been downvoted.
 * @param threadName {string} - The name of the thread to downvote the message in.
 * @param messageId {string} - The ID of the message to downvote.
 * @returns {Promise<Response>} - The response from the server.
 */
function downvoteMessage(threadName, messageId) {
    return fetch( `/api/thread/${threadName}/downvoteMessage`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            messageId: messageId
        })
    });
}

/**
 * Join the thread with the given name.
 * @description This function sends a request to join a thread. It does not handle the response.
 * @description But a success response means that the thread has been joined.
 * @description A BadRequest response means that the user already joined the thread.
 * @param threadName
 * @returns {Promise<Response>}
 */
function joinThread(threadName) {
    return fetch( `/api/thread/${threadName}/joinThread`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        }
    });
}

/**
 * Leave the thread with the given name.
 * @description This function sends a request to leave a thread. It does not handle the response.
 * @description But a success response means that the thread has been left.
 * @description A BadRequest response means that the user is already not in the thread.
 * @param threadName
 * @returns {Promise<Response>}
 */
function leaveThread(threadName) {
    return fetch( `/api/thread/${threadName}/leaveThread`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        }
    });
}

function getMessage(threadName, offset, order) {
    return fetch( `/api/messages?thread=${threadName}&offset=${offset}&order=${order}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        }
    });
}

