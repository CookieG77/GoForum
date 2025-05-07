/**
 * Get the current thread name from the URL.
 * @description This function extracts the thread name from the URL path. Only work if you are on the thread page. (e.g. /t/123)
 * @returns {string}
 */
function getCurrentThreadName() {
    const path = window.location.pathname;
    const segments = path.split("/");
    return segments[2];
}

/**
 * Send a message to the current thread.
 * @description This function sends a message to the current thread. It does not handle the response.
 * @description But a success response means that the message has been sent.
 * @param threadName {string} - The name of the thread to send the message to.
 * @param messageTitle {string} - The title of the message.
 * @param messageContent {string} - The content of the message.
 * @param messageMedias {string[]} - The media files to attach to the message.
 * @param messageTags {int[]} - The tags of the message.
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
 * @param reportReason {string} - The reason for reporting the message.
 * @param reportDescription {string} - The description of the report.
 * @returns {Promise<Response>} - The response from the server.
 */
function reportMessage(threadName, messageId, reportReason, reportDescription) {
    return fetch( `/api/thread/${threadName}/reportMessage`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            contentToReportID: messageId,
            reportType: reportReason,
            content: reportDescription
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
 * @param threadName {string} - The name of the thread to join.
 * @returns {Promise<Response>} - The response from the server.
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
 * @param threadName {string} - The name of the thread to leave.
 * @returns {Promise<Response>} - The response from the server.
 */
function leaveThread(threadName) {
    return fetch( `/api/thread/${threadName}/leaveThread`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        }
    });
}

/**
 * Get the messages from the current thread.
 * @description This function sends a request to get the messages from the current thread. It does not handle the response.
 * @description But a success response means that the messages have been retrieved.
 * @param threadName {string} - The name of the thread to get the messages from.
 * @param offset {number} - The offset to start getting the messages from.
 * @param order {string} - The order to get the messages in.
 * @param tags {string[]} - The tags to filter the messages by.
 * @returns {Promise<Response>} - The response from the server.
 */
function getMessage(threadName, offset, order, tags = []) {
    if (tags.length > 0) {
        return fetch( `/api/messages?thread=${threadName}&offset=${offset}&order=${order}&tags=${encodeURIComponent(JSON.stringify(tags))}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            }
        });
    }
    return fetch( `/api/messages?thread=${threadName}&offset=${offset}&order=${order}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        }
    });
}

/**
 * Send a comment to the message with the given id in the given thread.
 * @description This function sends a request to send a comment to a message in the current thread. It does not handle the response.
 * @param threadName {string} - The name of the thread to send the comment to.
 * @param messageId {string} - The ID of the message to send the comment to.
 * @param commentContent {string} - The content of the comment.
 * @returns {Promise<Response>} - The response from the server.
 */
function sendComment(threadName, messageId, commentContent) {
    return fetch( `/api/thread/${threadName}/sendComment`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            messageId: messageId,
            content: commentContent
        })
    });
}

/**
 * Delete a comment from the message with the given id in the given thread.
 * @description This function sends a request to delete a comment from a message in the current thread. It does not handle the response.
 * @param threadName {string} - The name of the thread to delete the comment from.
 * @param messageId {string} - The ID of the message to delete the comment from.
 * @param commentId {string} - The ID of the comment to delete.
 * @returns {Promise<Response>} - The response from the server.
 */
function deleteComment(threadName, messageId, commentId) {
    return fetch( `/api/thread/${threadName}/deleteComment`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            commentId: commentId,
            messageId: messageId
        })
    });
}

/**
 * Edit the comment with the given id in the given thread.
 * @description This function sends a request to edit a comment in the current thread. It does not handle the response.
 * @param threadName {string} - The name of the thread to edit the comment in.
 * @param messageId {string} - The ID of the message to edit the comment in.
 * @param commentId {string} - The ID of the comment to edit.
 * @param newCommentContent {string} - The new content of the comment.
 * @returns {Promise<Response>} - The response from the server.
 */
function editComment(threadName, messageId, commentId, newCommentContent) {
    return fetch( `/api/thread/${threadName}/editComment`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            commentId: commentId,
            messageId: messageId,
            content: newCommentContent
        })
    });
}

/**
 * Report the comment with the given id in the given thread.
 * @description This function sends a request to report a comment in the current thread. It does not handle the response.
 * @param threadName {string} - The name of the thread to report the comment in.
 * @param messageId {string} - The ID of the message to report the comment in.
 * @param commentId {string} - The ID of the comment to report.
 * @param reportReason {string} - The reason for reporting the comment.
 * @param reportDescription {string} - The description of the report.
 * @returns {Promise<Response>} - The response from the server.
 */
function reportComment(threadName, messageId, commentId, reportReason, reportDescription) {
    return fetch( `/api/thread/${threadName}/reportComment`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            contentToReportID: commentId,
            messageId: messageId,
            reportType: reportReason,
            content: reportDescription
        })
    });
}

/**
 * Upvote the comment with the given id in the given thread.
 * @description This function sends a request to upvote a comment in the current thread. It does not handle the response.
 * @param threadName {string} - The name of the thread to upvote the comment in.
 * @param messageId {string} - The ID of the message to upvote the comment in.
 * @param commentId {string} - The ID of the comment to upvote.
 * @returns {Promise<Response>} - The response from the server.
 */
function upvoteComment(threadName, messageId, commentId) {
    return fetch( `/api/thread/${threadName}/upvoteComment`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            commentId: commentId,
            messageId: messageId
        })
    });
}

/**
 * Downvote the comment with the given id in the given thread.
 * @description This function sends a request to downvote a comment in the current thread. It does not handle the response.
 * @param threadName {string} - The name of the thread to downvote the comment in.
 * @param messageId {string} - The ID of the message to downvote the comment in.
 * @param commentId {string} - The ID of the comment to downvote.
 * @returns {Promise<Response>} - The response from the server.
 */

function downvoteComment(threadName, messageId, commentId) {
    return fetch( `/api/thread/${threadName}/downvoteComment`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            commentId: commentId,
            messageId: messageId
        })
    });
}

/**
 * Ban the user with the given username from the given thread.
 * @description This function sends a request to ban a user from the current thread. It does not handle the response.
 * @description But a success response means that the user has been banned.
 * @description If you're not a dev looking at this code, you won't be able to use this function, the server double checks if the user has the right to do so. (●'◡'●)
 * @param threadName {string} - The name of the thread to ban the user from.
 * @param username {string} - The username of the user to ban.
 * @returns {Promise<Response>} - The response from the server.
 */
function banUser(threadName, username) {
    return fetch( `/api/thread/${threadName}/banUser`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            username: username
        })
    });
}

/**
 * Set a report to resolved.
 * @description This function sends a request to set a report to resolved. It does not handle the response.
 * @description But a success response means that the report has been set to resolved.
 * @param threadName
 * @param reportId
 * @returns {Promise<Response>}
 */
function setReportToResolved(threadName, reportId) {
    return fetch( `/api/thread/${threadName}/setReportToResolved`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            reportId: reportId
        })
    });
}

/**
 * Get the comments from the message with the given id in the given thread.
 * @description This function sends a request to get the comments from a message in the current thread. It does not handle the response.
 * @description But a success response means that the comments have been retrieved.
 * @param threadName {string} - The name of the thread to get the comments from.
 * @param offset {number} - The offset to start getting the comments from.
 * @param messageId {string} - The ID of the message to get the comments from.
 * @returns {Promise<Response>} - The response from the server.
 */
function getComment(threadName, offset, messageId) {
    return fetch( `/api/comments?thread=${threadName}&offset=${offset}&message=${messageId}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        }
    });
}

/**
 * Create a tag in the given thread.
 * @description This function sends a request to create a tag in the current thread. It does not handle the response.
 * @description But a success response means that the tag has been created.
 * @param threadName {string} - The name of the thread to create the tag in.
 * @param tagName {string} - The name of the tag to create.
 * @param tagColor {string} - The color of the tag to create.
 * @returns {Promise<Response>} - The response from the server.
 */
function createThreadTag(threadName, tagName, tagColor) {
    console.log("Creating tag with name: " + tagName + " and color: " + tagColor);
    return fetch( `/api/thread/${threadName}/createThreadTag`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            tagName: tagName,
            tagColor: tagColor
        })
    });
}

/**
 * Edit the tag with the given id in the given thread.
 * @description This function sends a request to edit a tag in the current thread. It does not handle the response.
 * @description But a success response means that the tag has been edited.
 * @param threadName {string} - The name of the thread to edit the tag in.
 * @param tagId {string} - The ID of the tag to edit.
 * @param tagName {string} - The new name of the tag.
 * @param tagColor {string} - The new color of the tag.
 * @returns {Promise<Response>} - The response from the server.
 */
function editThreadTag(threadName, tagId, tagName, tagColor) {
    return fetch( `/api/thread/${threadName}/editThreadTag`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            tagId: tagId,
            tagName: tagName,
            tagColor: tagColor
        })
    });
}

/**
 * Delete the tag with the given id in the given thread.
 * @description This function sends a request to delete a tag in the current thread. It does not handle the response.
 * @description But a success response means that the tag has been deleted.
 * @param threadName {string} - The name of the thread to delete the tag from.
 * @param tagId {string} - The ID of the tag to delete.
 * @returns {Promise<Response>} - The response from the server.
 */
function deleteThreadTag(threadName, tagId) {
    return fetch( `/api/thread/${threadName}/deleteThreadTag`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            tagId: tagId
        })
    });
}

/**
 * Get the tags from the given thread.
 * @description This function sends a request to get the tags from the current thread. It does not handle the response.
 * @description But a success response means that the tags have been retrieved.
 * @param threadName {string} - The name of the thread to get the tags from.
 * @returns {Promise<Response>} - The response from the server.
 */
function getThreadTags(threadName) {
    return fetch( `/api/threadTags?thread=${threadName}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        }
    });
}

/**
 * Promote the user with the given username in the given thread.
 * @description This function sends a request to promote a user in the current thread. It does not handle the response.
 * @description But a success response means that the user has been promoted.
 * @param threadName {string} - The name of the thread to promote the user in.
 * @param username {string} - The username of the user to promote.
 * @returns {Promise<Response>} - The response from the server.
 */
function promoteUser(threadName, username) {
    return fetch( `/api/thread/${threadName}/promoteUser`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            username: username
        })
    });
}

/**
 * Demote the user with the given username in the given thread.
 * @description This function sends a request to demote a user in the current thread. It does not handle the response.
 * @description But a success response means that the user has been demoted.
 * @param threadName {string} - The name of the thread to demote the user in.
 * @param username {string} - The username of the user to demote.
 * @returns {Promise<Response>} - The response from the server.
 */
function demoteUser(threadName, username) {
    return fetch( `/api/thread/${threadName}/demoteUser`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            username: username
        })
    });
}


/**
 * Get the i18n text for the given key.
 * @param key {string} - The key to get the text for.
 * @param value {string|null} - The value to replace in the text. If null, no replacement is done.
 * @returns {string} - The i18n text for the given key.
 */
function getI18nText(key, value = null) {
    const el = document.querySelector(`#i18n [data-key="${key}"]`);
    if (!el) return '';
    return value !== null ? el.textContent.replace('{n}', value) : el.textContent;
}

/**
 * Get the time ago string for the given date.
 * @param dateString {string} - The date string to convert to a time ago string.
 * @returns {string} - The time ago string for the given date.
 */
function timeAgo(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;

    const seconds = Math.floor(diffMs / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours   = Math.floor(minutes / 60);
    const days    = Math.floor(hours / 24);
    const months  = Math.floor(days / 30);

    if (minutes < 1) {
        return getI18nText('ago-seconds');
    } else if (minutes < 60) {
        return getI18nText(minutes === 1 ? 'ago-minute' : 'ago-minutes', minutes);
    } else if (hours < 24) {
        return getI18nText(hours === 1 ? 'ago-hour' : 'ago-hours', hours);
    } else if (days < 30) {
        return getI18nText(days === 1 ? 'ago-day' : 'ago-days', days);
    } else {
        const options = { month: 'long', year: 'numeric' };
        return date.toLocaleDateString(undefined, options); // utilise la locale du navigateur
    }
}

/**
 * Update the vote visual based on the current vote state.
 * @description This function updates the vote visual based on the current vote state.
 * @description It changes the image source of the upvote and downvote images based on the current vote state.
 * @param currentVoteState {string} - The current vote state. 1 for upvote, 0 for no vote, -1 for downvote.
 * @param upvoteImg {HTMLImageElement} - The image element for the upvote button.
 * @param downvoteImg {HTMLImageElement} - The image element for the downvote button.
 * @returns {void} - This function does not return anything.
 */
function updateVoteVisual(currentVoteState, upvoteImg, downvoteImg) {
    let state = parseInt(currentVoteState);
    if (state === 1){
        upvoteImg.src = `/img/upvote.png`
        downvoteImg.src = `/img/downvote_empty.png`
    } else if(state === 0){
        upvoteImg.src = `/img/upvote_empty.png`
        downvoteImg.src = `/img/downvote_empty.png`
    } else if(state === -1){
        upvoteImg.src = `/img/upvote_empty.png`
        downvoteImg.src = `/img/downvote.png`
    }
}

/**
 * Update the vote state based on the current vote state and the action taken.
 * @description This function updates the vote state based on the current vote state and the action taken.
 * @description It changes the vote state and the vote count based on the action taken.
 * @param currentVoteState {string} - The current vote state. 1 for upvote, 0 for no vote, -1 for downvote.
 * @param currentVoteCount {number} - The current vote count.
 * @param IsUpvoting {boolean} - Whether the user is upvoting or downvoting.
 * @param upvoteImg {HTMLImageElement} - The image element for the upvote button.
 * @param downvoteImg {HTMLImageElement} - The image element for the downvote button.
 * @returns {{state: number, count}} - The updated vote state and vote count.
 */
function updateVoteState(currentVoteState, currentVoteCount, IsUpvoting, upvoteImg, downvoteImg) {
    let state = parseInt(currentVoteState);
    let count = currentVoteCount;
    if (state === 1){
        if (IsUpvoting){ // upvoting when already upvoted so we remove the upvote
            state = 0;
            count -= 1;
        } else { // downvoting when already upvoted so we remove the upvote and add a downvote
            state = -1;
            count -= 2;
        }
    } else if(state === 0){
        if (IsUpvoting){ // upvoting when not voted so we add an upvote
            state = 1;
            count += 1;
        } else { // downvoting when not voted so we add a downvote
            state = -1;
            count -= 1;
        }
    } else if(state === -1){
        if (IsUpvoting){ // upvoting when already downvoted so we remove the downvote and add an upvote
            state = 1;
            count += 2;
        } else { // downvoting when already downvoted so we remove the downvote
            state = 0;
            count += 1;
        }
    }
    updateVoteVisual(state, upvoteImg, downvoteImg);
    console.log("Current Vote State: ", state);
    console.log("Current Vote Count: ", count);
    return {state, count};
}