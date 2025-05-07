document.addEventListener("DOMContentLoaded", function (){
    const loadMoreCommentsButton = document.getElementById("load-more-comments-button");
    const postContent = document.getElementById("t-post-content");
    const mediaContainer = document.createElement("div");
    const medias = document.createElement("div");
    const mediaLinks = [];
    const postVoteCountSpan = document.getElementById("t-vote-count");
    let postVoteCount = parseInt(getI18nText("data_upvotes"), 10) - parseInt(getI18nText("data_downvotes"), 10);
    const postVoteUpImage = document.getElementById("t-post-vote-up-image");
    const postVoteUpButton = document.getElementById("t-post-vote-up-button");
    const postVoteDownImage = document.getElementById("t-post-vote-down-image");
    const postVoteDownButton = document.getElementById("t-post-vote-down-button");
    let voteState = parseInt(getI18nText("data_postVoteState"), 10);
    const postDateContainer = document.getElementById("t-post-date");
    let postDate = getI18nText("data_postDate");
    for (const medialink of document.getElementById("data_MediaLinks").children){
        mediaLinks.push(medialink.textContent);
    }

    const threadName = getCurrentThreadName();
    const messageId = document.getElementById("data_postID").textContent;
    let userIsAuthenticated = document.getElementById("isAuthenticated").textContent === "true";
    let userIsAMember = document.getElementById("isAMember").textContent === "true";
    let userRank = parseInt(document.getElementById("userRank").textContent,10);
    let userIsModerator = userRank >= 1;
    let userIsAdmin = userRank >= 2;
    let userIsThreadOwner = userRank >= 3;
    let offset= 0;
    let hasReachedEnd = false;
    const commentsContainer = document.getElementById("comments-container");

    let newCommentButton;
    let newCommentContent;
    let newCommentContentCharCountValue;

    const scrollbar = document.getElementsByClassName("custom-scrollbar")[0];

    // Report menu elements
    const reportMenu = document.getElementById("report-button-menu");
    const reportMenuBackground = reportMenu.getElementsByClassName("full-screens-menu-background")[0];
    const reportMenuCloseButton = document.getElementById("close-report-menu");
    const reportMenuSendButton = document.getElementById("send-report-button");
    const reportReason = document.getElementById("report-reason");
    const reportContent = document.getElementById("report-content");
    const reportContentCharCountValue = document.getElementById("report-content-char-count-value");
    const reportMenuSuccessMessage = document.getElementById("report-success-message");
    const reportMenuErrorMessage = document.getElementById("report-error-message");
    let commentToReport = null;

    // Edit comment menu elements
    const editCommentMenu = document.getElementById("edit-comment-button-menu");
    const editCommentMenuBackground = editCommentMenu.getElementsByClassName("full-screens-menu-background")[0];
    const editCommentMenuCloseButton = document.getElementById("close-comment-post-menu");
    const editCommentMenuNewContentField = document.getElementById("edit-comment-content");
    const editCommentMenuNewContentCharCountValue = document.getElementById("edit-comment-content-char-count-value");
    const editCommentMenuSendButton = document.getElementById("edit-comment-send-button");
    let editedCommentID = null;

    /**
     * Show the report menu for a message.
     * @description This function displays the report menu and sets the message ID to report.
     * @param messageID {number} - The ID of the message to report
     * @param isPost {boolean} - Whether the message is a post or a comment (true for post, false for comment)
     */
    function showReportMenu(messageID) {
        reportMenu.classList.remove("hidden");
        scrollbar.classList.add("hidden");
        commentToReport = messageID;
        reportMenuSendButton.disabled = true;
    }

    /**
     * Hide the report menu.
     * @description This function hides the report menu and resets the fields.
     */
    function hideReportMenu() {
        reportMenu.classList.add("hidden");
        scrollbar.classList.remove("hidden");
        reportMenuSuccessMessage.classList.add("hidden");
        reportMenuErrorMessage.classList.add("hidden");
        reportContent.value = "";
        commentToReport = null;
        reportMenuSendButton.disabled = true;
    }

    function showEditMenu(commentID, content) {
        editCommentMenu.classList.remove("hidden");
        scrollbar.classList.add("hidden");
        editCommentMenuNewContentField.value = content;
        editedCommentID = commentID;
        editCommentMenuSendButton.disabled = true;

        // Update the character count
        const charCount = editCommentMenuNewContentField.value.length;
        editCommentMenuNewContentCharCountValue.innerText = `${charCount}`;
        if (charCount > 500) {
            editCommentMenuNewContentField.value = editCommentMenuNewContentField.value.substring(0, 500);
            editCommentMenuNewContentCharCountValue.innerText = `500`;
        }
        editCommentMenuSendButton.disabled = (charCount < 5 || charCount > 500);
    }

    function hideEditMenu() {
        editCommentMenu.classList.add("hidden");
        scrollbar.classList.remove("hidden");
        editCommentMenuNewContentField.value = "";
        editedCommentID = null;
        editCommentMenuSendButton.disabled = true;
    }

    if (userIsAuthenticated && userIsAMember) {
        newCommentButton = document.getElementById("new-comment-send-button");
        newCommentContent = document.getElementById("new-comment-content");
        newCommentContentCharCountValue = document.getElementById("new-comment-content-char-count-value");
    }

    if (voteState === 1) {
        postVoteUpImage.src = `/img/upvote.png`;
        postVoteDownImage.src = `/img/downvote_empty.png`;
    }
    else if (voteState === -1) {
        postVoteUpImage.src = `/img/upvote_empty.png`;
        postVoteDownImage.src = `/img/downvote.png`;
    } else {
        postVoteUpImage.src = `/img/upvote_empty.png`;
        postVoteDownImage.src = `/img/downvote_empty.png`;
    }

    if (!userIsAuthenticated) {
        postVoteUpButton.disabled = true;
        postVoteDownButton.disabled = true;
    }

    mediaContainer.classList.add("post-media-container");
    mediaContainer.id = "t-post-media-container";
    postContent.appendChild(mediaContainer);

    medias.classList.add("post-medias");
    if (mediaLinks.length !== 0) {
        mediaContainer.classList.add("win95-border-indent");
        for (let i = 0; i < mediaLinks.length; i++) {
            const media = mediaLinks[i];
            var mediaElement = document.createElement("img");
            mediaElement.src = `/upload/${media}`;
            mediaElement.alt = `Media[${media}]`;
            mediaElement.draggable = false;
            mediaElement.classList.add("post-picture", "unselectable");
            mediaElement.loading ="lazy";
            mediaElement.style.display = "none";
            medias.appendChild(mediaElement);
        }

        mediaContainer.appendChild(medias);

        if (mediaLinks.length > 1){
            const prev = document.createElement("button");
            const prevImg = document.createElement("img");
            prev.classList.add("prev-button", "win95-button");
            prevImg.src = `/img/prev_arrow.png`;
            prevImg.alt = "prev";
            prevImg.draggable = false;
            prevImg.classList.add("prev-img", "unselectable");
            prev.appendChild(prevImg);
            medias.appendChild(prev);

            const next = document.createElement("button");
            const nextImg = document.createElement("img");
            next.classList.add("next-button", "win95-button");
            nextImg.src = `/img/next_arrow.png`;
            nextImg.alt = "next";
            nextImg.draggable = false;
            nextImg.classList.add("next-img", "unselectable");
            next.appendChild(nextImg);
            medias.appendChild(next);

            const dots = document.createElement("div");
            dots.classList.add("dots");
            mediaLinks.forEach((_, i) => {
                const dot = document.createElement("img");
                dot.src = `/img/dot_inactive.png`
                dot.classList.add("dot");
                dot.dataset.index = i;
                dots.appendChild(dot);
            });
            mediaContainer.appendChild(dots);

            let currentSlide = 0;
            const slides = Array.from(medias.querySelectorAll(".post-picture"));
            const allDots = Array.from(dots.children);

            function show(n) {
                if (n < 0) n = slides.length - 1;
                if (n >= slides.length) n = 0;
                currentSlide = n;
                slides.forEach((s, i) => s.style.display = i === n ? "block" : "none");
                allDots.forEach((d, i) => {
                    d.src = i === n
                        ? `/img/dot_active.png`
                        : `/img/dot_inactive.png`;
                });
            }

            prev.addEventListener("click", () => show(currentSlide - 1));
            next.addEventListener("click", () => show(currentSlide + 1));
            dots.addEventListener("click", e => {
                if (e.target.matches(".dot")) {
                    show(+e.target.dataset.index);
                }
            });

            show(0);
        } else if (mediaLinks.length === 1){
            const onlySlide = medias.querySelector('.post-picture');
            onlySlide.style.display = 'block';
        }
    }

    /**
     * Load more comments from the server.
     * @description This function sends a request to get more comments from the server.
     * @description Then it adds the comments to the comments container.
     */
    function loadMoreComments() {
        if (hasReachedEnd) {
            return;
        }
        let res = getComment(threadName, offset, messageId);
        res.then(async (response) => {
            if (response.ok) {
                const data = await response.json();
                if (data == null) {
                    hasReachedEnd = true;
                    loadMoreCommentsButton.disabled = true;
                    loadMoreCommentsButton.innerText = "No more comments";
                    return;
                }
                for (let i = 0; i < data.length; i++) {
                    const comment = data[i];
                    const postElement = createNewComment(comment);
                    commentsContainer.appendChild(postElement);
                }
                offset += data.length;
            } else {
                console.error(response);
            }
        });
    }

    /**
     * Create a new comment element.
     * @description This function creates a new comment element from the given data.
     * @param data {object} - The data of the comment to create.
     * @returns {HTMLElement} - The new comment element.
     */
    function createNewComment(data){
        const container = document.createElement("div");
        const commentHeader = document.createElement("section");
        const commentAuthor = document.createElement("div");
        const authorPfp = document.createElement("img");
        const authorAndTime = document.createElement("div");
        const author = document.createElement("span");
        const option = document.createElement("div");
        const optionButton = document.createElement("button");
        const optionMenu = document.createElement("div");
        const commentContent = document.createElement("section");
        const commentMedia = document.createElement("p");
        const commentVote = document.createElement("section");
        const voteField = document.createElement("div");
        const upvoteButton = document.createElement("button");
        const upvoteImg = document.createElement("img");
        const vote = document.createElement("span");
        const downvoteButton = document.createElement("button");
        const downvoteImg = document.createElement("img");
        const dateAndEditContainer = document.createElement("section");
        const dateSpan = document.createElement("span");
        const isEditedSpan = document.createElement("span");

        let currentVoteState = data.vote_state;
        let currentVoteCount = data.up_votes - data.down_votes;
        let isCommentOwner = data.user_name === document.getElementById("username").textContent;
        console.log(data.user_name);
        console.log(document.getElementById("username").textContent);

        container.classList.add("comment-box", "win95-border");

        commentHeader.classList.add("comment-header", "win95-header");
        container.appendChild(commentHeader);

        commentAuthor.classList.add("comment-profile");
        commentHeader.appendChild(commentAuthor);

        authorPfp.src = `/upload/${data.user_pfp_address}`;
        authorPfp.alt = "Author profile picture";
        authorPfp.classList.add("comment-profile-picture");
        authorPfp.draggable = false;
        authorPfp.onclick = function (){
            window.location.href = `/profile/${data.user_name}`
        }
        commentAuthor.appendChild(authorPfp);

        authorAndTime.classList.add("author-and-time");
        commentAuthor.appendChild(authorAndTime);

        author.classList.add("author-pseudo");
        author.innerText = `${data.user_name}`;
        author.onclick = function (){
            window.location.href = `/profile/${data.user_name}`
        }
        authorAndTime.appendChild(author);

        option.classList.add();
        commentHeader.appendChild(option);

        optionButton.innerText = "...";
        optionButton.type = "button";
        optionButton.classList.add("win95-button");
        option.appendChild(optionButton);

        optionMenu.classList.add("option-menu", "win95-border");

        let optionMenuReportButtonHTML = `
            <li class="win95-menu-button message-report menu-button" id="comment-report-button-p${data.comment_id}">
                <img src="/img/report.png" alt="" class="win95-minor-logo unselectable" draggable="false">
                <span>${getI18nText("option-menu-report-button-text")}</span>
            </li>`

        let optionMenuEditButtonHTML = `
            <li class="win95-menu-button message-edit menu-button" id="comment-edit-button-p${data.comment_id}">
                <img src="/img/edit.png" alt="" class="win95-minor-logo  unselectable" draggable="false">
                <span>${getI18nText("option-menu-edit-button-text")}</span>
            </li>`

        let optionMenuDeleteButtonHTML = `
            <li class="win95-menu-button message-delete menu-button" id="comment-delete-button-p${data.comment_id}">
                <img src="/img/delete.png" alt="" class="win95-minor-logo  unselectable" draggable="false">
                <span>${getI18nText("option-menu-delete-button-text")}</span>
            </li>`

        let optionMenuBanButtonHTML = `
            <li class="win95-menu-button message-ban menu-button" id="comment-ban-button-p${data.comment_id}">
                <img src="/img/ban.png" alt="" class="win95-minor-logo unselectable" draggable="false">
                <span>${getI18nText("option-menu-ban-button-text")}</span>
            </li>`
        let additionalButtonsHTML = "";
        let showReportButton = false;
        let showEditButton = false;
        let showDeleteButton = false;
        let showBanButton = false;

        if (userIsAuthenticated) { // If the user is authenticated he can see the option menu
            if (!isCommentOwner) { // If the user is authenticated he can report a post (exept his posts)
                additionalButtonsHTML += optionMenuReportButtonHTML;
                showReportButton = true;
            }
            if (isCommentOwner) { // If the user is the owner of the post he can edit it
                additionalButtonsHTML += optionMenuEditButtonHTML;
                showEditButton = true;
            }
            if (isCommentOwner || userIsModerator || userIsAdmin || userIsAdmin) { // If the user is the owner of the post or his rank is moderator or higher he can delete it
                additionalButtonsHTML += optionMenuDeleteButtonHTML;
                showDeleteButton = true;
            }
            if ((userIsThreadOwner || userIsAdmin) && !isCommentOwner) { // If the user rank is admin or higher he can ban the user (exept himself)
                additionalButtonsHTML += optionMenuBanButtonHTML;
                showBanButton = true;
            }
        }

        // All user can report a post
        optionMenu.innerHTML =`
        <ul>
            ${additionalButtonsHTML}
        </ul>
        
        `
        option.appendChild(optionMenu);

        // Add the event listener to the edit button
        if (showEditButton) {
            const editButton = optionMenu.querySelector(`#comment-edit-button-p${data.comment_id}`);
            editButton.addEventListener("click", function() {
                console.log(`Edit button clicked for comment ${data.comment_id}`);
                showEditMenu(data.comment_id, data.comment_content);
            });
        }

        // Add the event listener to the delete button
        if (showDeleteButton) {
            const deleteButton = optionMenu.querySelector(`#comment-delete-button-p${data.comment_id}`);
            deleteButton.addEventListener("click", function() {
                console.log(`Delete button clicked for post ${data.comment_id}`);
                const result = deleteComment(threadName, messageId, data.comment_id);
                result.then(async (response) => {
                    if (response.ok) {
                        container.remove();
                        console.log("Comment deleted successfully");
                    } else {
                        alert("Error while deleting comment : " + response.statusText);
                        console.error(response);
                    }
                });
            });
        }

        // Add the event listener to the report button
        if (showReportButton) {
            const reportButton = optionMenu.querySelector(`#comment-report-button-p${data.comment_id}`);
            reportButton.addEventListener("click", function() {
                console.log(`Report button clicked for post ${data.comment_id}`);
                showReportMenu(data.comment_id);
            });
        }

        // Add the event listener to the ban button
        if (showBanButton) {
            const banButton = optionMenu.querySelector(`#comment-ban-button-p${data.comment_id}`);
            banButton.addEventListener("click", function() {
                console.log(`Ban button clicked for post ${data.comment_id}`);
                const result = banUser(threadName, data.user_name);
                result.then(async (response) => {
                    if (response.ok) {
                        alert("User banned successfully");
                        console.log("User banned successfully");
                    } else {
                        alert("Error while banning user : " + response.statusText);
                        console.error(response);
                    }
                });
            });
        }

        function toggleOptionMenu() {
            optionMenu.classList.add("active")
        }

        optionButton.addEventListener("click", (e) => {
            if (!optionMenu.classList.contains("active")){
                const options = document.querySelectorAll(".option-menu");
                e.stopPropagation();
                options.forEach((opt)=>{
                    opt.classList.remove("active");
                })
                toggleOptionMenu();
            }
        })

        window.addEventListener("click", (e) => {
            if (!optionMenu.contains(e.target)){
                optionMenu.classList.remove("active");
            }
        })

        commentContent.classList.add("comment-content");
        container.appendChild(commentContent);

        commentMedia.classList.add("comment-media","win95-border-indent");
        commentMedia.innerText = data.comment_content;
        commentContent.appendChild(commentMedia);

        commentVote.classList.add("comment-vote");
        container.appendChild(commentVote);

        voteField.classList.add("comment-vote-field");
        commentVote.appendChild(voteField);

        upvoteButton.type = "button";
        upvoteButton.classList.add("win95-button", "comment-vote-button");
        if (userIsAuthenticated) {
            upvoteButton.addEventListener("click", function () {
                upvoteComment(threadName, messageId, data.comment_id)
                    .then(r => {
                        if (r.ok) {
                            return r.json();
                        } else {
                            throw new Error("Error while upvoting comment");
                        }
                    })
                    .then(data => {
                        console.log("Comment upvoted successfully", data);
                    })
                    .catch(error => {
                        console.error("Error:", error);
                    });
                const {state, count} = updateVoteState(currentVoteState, currentVoteCount, true, upvoteImg, downvoteImg);
                currentVoteState = state;
                currentVoteCount = count;
                vote.innerText = currentVoteCount;
            });
        } else {
            upvoteButton.disabled = true;
        }
        voteField.appendChild(upvoteButton);

        upvoteImg.src = `/img/upvote_empty.png`;
        upvoteImg.alt = "Upvote image";
        upvoteImg.classList.add("comment-vote-image", "unselectable");
        upvoteImg.draggable = false;
        upvoteButton.appendChild(upvoteImg);

        vote.classList.add("comment-vote-value");
        vote.innerText = `${data.up_votes - data.down_votes}`;
        voteField.appendChild(vote);

        downvoteButton.type = "button";
        downvoteButton.classList.add("win95-button", "comment-vote-button");
        if (userIsAuthenticated) {
            downvoteButton.addEventListener("click", function () {
                downvoteComment(threadName, messageId, data.comment_id)
                    .then(r => {
                        if (r.ok) {
                            return r.json();
                        } else {
                            throw new Error("Error while downvoting comment");
                        }
                    })
                    .then(data => {
                        console.log("Comment downvoted successfully", data);
                    })
                    .catch(error => {
                        console.error("Error:", error);
                    });
                const {state, count} = updateVoteState(currentVoteState, currentVoteCount, false, upvoteImg, downvoteImg);
                currentVoteState = state;
                currentVoteCount = count;
                vote.innerText = currentVoteCount;
            });
        } else {
            downvoteButton.disabled = true;
        }
        voteField.appendChild(downvoteButton);

        downvoteImg.src = `/img/downvote_empty.png`;
        downvoteImg.alt = "Downvote image";
        downvoteImg.classList.add("comment-vote-image", "unselectable");
        downvoteImg.draggable = false;
        downvoteButton.appendChild(downvoteImg);
        // Make sure to update the vote visual when the comment is created so that if the user is already upvoting or downvoting the post, the visual is correct
        updateVoteVisual(currentVoteState, upvoteImg, downvoteImg);


        container.appendChild(dateAndEditContainer);
        dateAndEditContainer.appendChild(dateSpan);
        dateAndEditContainer.appendChild(isEditedSpan);
        console.log(dateAndEditContainer);
        dateAndEditContainer.classList.add("comment-date-edit-container");
        dateSpan.classList.add("comment-date");
        dateSpan.innerText = timeAgo(data.creation_date);
        isEditedSpan.classList.add("comment-edited");
        if (data.was_edited) {
            isEditedSpan.innerText = getI18nText("was-edited");
        } else {
            isEditedSpan.innerText = "";
        }
        return container;
    }

    loadMoreCommentsButton.addEventListener('click', function() {
        loadMoreComments();
    })

    if (userIsAuthenticated && userIsAMember) {
        postVoteUpButton.addEventListener('click', function() {
            postId = parseInt(messageId, 10);
            upvoteMessage(threadName, postId)
                .then(r => {
                    if (r.ok) {
                        return r.json();
                    } else {
                        throw new Error("Error while upvoting post");
                    }
                })
                .then(data => {
                    console.log("Post upvoted successfully", data);
                })
                .catch(error => {
                    console.error("Error:", error);
                });
            const {state, count} = updateVoteState(voteState, postVoteCount, true, postVoteUpImage, postVoteDownImage);
            voteState = state;
            postVoteCount = count;
            postVoteCountSpan.innerText = `${postVoteCount}`;
        });

        postVoteDownButton.addEventListener('click', function() {
            postId = parseInt(messageId, 10);
            downvoteMessage(threadName, postId)
                .then(r => {
                    if (r.ok) {
                        return r.json();
                    } else {
                        throw new Error("Error while downvoting post");
                    }
                })
                .then(data => {
                    console.log("Post downvoted successfully", data);
                })
                .catch(error => {
                    console.error("Error:", error);
                });
            const {state, count} = updateVoteState(voteState, postVoteCount, false, postVoteUpImage, postVoteDownImage);
            voteState = state;
            postVoteCount = count;
            postVoteCountSpan.innerText = `${postVoteCount}`;
        });

        newCommentContent.addEventListener('input' , function() {
            const charCount = newCommentContent.value.length;
            newCommentContentCharCountValue.innerText = `${charCount}`;
            if (charCount > 500) {
                newCommentContent.value = newCommentContent.value.substring(0, 500);
            }
            newCommentButton.disabled = (charCount < 5 || charCount > 500);
        });

        newCommentButton.addEventListener("click", function() {
            const commentContent = newCommentContent.value;
            const postId = messageId.toString();
            const threadName = getCurrentThreadName();
            sendComment(threadName, postId, commentContent)
                .then(r => {
                    if (r.ok) {
                        return r.json();
                    } else {
                        throw new Error("Error while creating comment");
                    }
                })
                .then(data => {
                    console.log("Comment created successfully", data);
                    const postElement = createNewComment(data);
                    commentsContainer.appendChild(postElement);
                    newCommentContent.value = "";
                    newCommentContentCharCountValue.innerText = "0";

                    // reload the comments
                    offset = 0;
                    hasReachedEnd = false;
                    loadMoreCommentsButton.disabled = false;
                    loadMoreCommentsButton.innerText = "Load more comments";
                    commentsContainer.innerHTML = "";
                    loadMoreComments();
                })
                .catch(error => {
                    console.error("Error:", error);
                });
        });
    }

    // Close the report menu when the close button is clicked
    reportMenuBackground.addEventListener('click', hideReportMenu);
    // Close the report menu when the close button is clicked
    reportMenuCloseButton.addEventListener('click', hideReportMenu);
    // Send the report when the send button is clicked
    reportMenuSendButton.addEventListener("click", function() {
        if (reportMenuSendButton.disabled) {
            return;
        }
        reportComment(threadName, messageId, commentToReport, reportReason.value, reportContent.value)
            .then(r => {
                if (r.ok) {
                    reportMenuSuccessMessage.classList.remove("hidden");
                    reportMenuSendButton.disabled = true;
                    reportReason.disabled = true;
                    reportContent.disabled = true;

                } else {
                    reportMenuErrorMessage.classList.remove("hidden");
                    console.error(r);
                }
            });
    });
    // Update the report content character count
    reportContent.addEventListener("input", function() {
        const charCount = reportContent.value.length;
        reportContentCharCountValue.innerText = `${charCount}`;
        if (charCount > 500) {
            reportContent.value = reportContent.value.substring(0, 500);
            reportContentCharCountValue.innerText = `500`;
        }
        reportMenuSendButton.disabled = (charCount < 20 || charCount > 500);
    });

    // Close the edit menu when the close button is clicked
    editCommentMenuBackground.addEventListener('click', hideEditMenu);
    // Close the edit menu when the close button is clicked
    editCommentMenuCloseButton.addEventListener('click', hideEditMenu);
    // Send the edit when the send button is clicked
    editCommentMenuSendButton.addEventListener("click", function() {
        if (editCommentMenuSendButton.disabled) {
            return;
        }
        editComment(threadName, messageId, editedCommentID, editCommentMenuNewContentField.value)
            .then(r => {
                if (r.ok) {
                    hideEditMenu();
                    commentsContainer.innerHTML = "";
                    hasReachedEnd = false;
                    offset = 0;
                    // Reload the comments
                    loadMoreComments();
                } else {
                    console.error(r);
                }
            });
    });
    // Update the edit content character count
    editCommentMenuNewContentField.addEventListener("input", function() {
        const charCount = editCommentMenuNewContentField.value.length;
        editCommentMenuNewContentCharCountValue.innerText = `${charCount}`;
        if (charCount > 500) {
            editCommentMenuNewContentField.value = editCommentMenuNewContentField.value.substring(0, 500);
            editCommentMenuNewContentCharCountValue.innerText = `500`;
        }
        editCommentMenuSendButton.disabled = (charCount < 5 || charCount > 500);
    });

    loadMoreComments()
    // Update the vote count
    postVoteCountSpan.innerText = `${postVoteCount}`;
    // Update the post date
    postDateContainer.innerText = timeAgo(postDate);
});

