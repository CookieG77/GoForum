document.addEventListener("DOMContentLoaded", function () {

    const leaveButton = document.getElementById("LeaveThreadButton");
    const joinButton = document.getElementById("JoinThreadButton");
    const editThreadButton = document.getElementById("EditThreadButton");

    // Debug button
    const loadMorePostsButton = document.getElementById("load-more-posts-button");

    const threadName = getCurrentThreadName()
    let userIsAuthenticated = document.getElementById("isAuthenticated").textContent === "true";
    let userIsMember = document.getElementById("isAMember").textContent === "true";
    let userRank = parseInt(document.getElementById("userRank").textContent,10);
    let userIsModerator = userRank >= 1;
    let userIsAdmin = userRank >= 2;
    let userIsThreadOwner = userRank >= 3;
    let offset= 0;
    let hasReachedEnd = false;
    let orderSelect = document.getElementById("order")
    const postsContainer = document.getElementById("posts-container");
    let selectedTags = [];

    const tagListContainer = document.getElementById('tagList');
    const noTagsMessage = document.getElementById('noTagsMessage');

    const newPostContainer = document.getElementById("new-post-box");
    const newPostTitle = document.getElementById("new-post-title");
    const newPostContent = document.getElementById("new-post-content");
    const newPostButton = document.getElementById("new-post-send-button");
    const newPostContentCharCountValue = document.getElementById("new-post-content-char-count-value");
    const newPostfileInput = document.getElementById("new-post-file-input")
    const newPostImagesPreview = document.getElementById("new-post-medias-container")
    const newPostTagsContainer = document.getElementById("new-post-tags-container");
    let MediaIDs = [];
    let titleValid = false;
    let contentValid = false;
    let newPostTags = [];

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
    let messageToReport = null;

    // Edit post menu elements
    const editMenu = document.getElementById("edit-post-button-menu");
    const editMenuBackground = editMenu.getElementsByClassName("full-screens-menu-background")[0];
    const editMenuCloseButton = document.getElementById("close-edit-post-menu");
    const editMenuNewTitleField = document.getElementById("edit-post-title");
    const editMenuNewContentField = document.getElementById("edit-post-content");
    const editMenuNewContentCharCountValue = document.getElementById("edit-post-content-char-count-value");
    const editMenuSendButton = document.getElementById("edit-post-send-button");
    let editedPostMedias = [];
    let editedPostID = null;
    const editMenuMediasPreview = document.getElementById("edit-post-medias-container");

    /**
     * Show the report menu for a message.
     * @description This function displays the report menu and sets the message ID to report.
     * @param messageID {number} - The ID of the message to report
     */
    function showReportMenu(messageID) {
        reportMenu.classList.remove("hidden");
        scrollbar.classList.add("hidden");
        messageToReport = messageID;
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
        messageToReport = null;
        reportMenuSendButton.disabled = false;
    }

    function showEditMenu(messageID, title, content, medias) {
        editMenu.classList.remove("hidden");
        scrollbar.classList.add("hidden");
        editMenuNewTitleField.value = title;
        editMenuNewContentField.value = content;
        editedPostID = messageID;
        editedPostMedias = medias;
    }

    function hideEditMenu() {
        editMenu.classList.add("hidden");
        scrollbar.classList.remove("hidden");
        editMenuNewTitleField.value = "";
        editMenuNewContentField.value = "";
        editedPostID = null;
        editedPostMedias = [];
    }

    /**
     * Load more posts from the server.
     * @description This function sends a request to get more posts from the server.
     * @description Then it adds the posts to the posts container.
     */
    function loadMorePosts() {
        if (hasReachedEnd) {
            return;
        }
        let res = getMessage(threadName, offset, orderSelect.value, selectedTags);
        res.then(async (response) => {
            if (response.ok) {
                const data = await response.json();
                console.log(data);
                if (data == null) {
                    hasReachedEnd = true;
                    loadMorePostsButton.disabled = true;
                    loadMorePostsButton.innerText = "No more posts";
                    return;
                }
                for (let i = 0; i < data.length; i++) {
                    const post = data[i];
                    const postElement = createNewPost(post);
                    postsContainer.appendChild(postElement);
                }
                offset += data.length;
            } else {
                console.error(response);
            }
        });
    }

    /**
     * Render the tags in the tag list container.
     * @description This function takes an array of tags and creates a visual representation of them.
     * @param tags {Array} - The array of tags to render.
     * @param tagContainer {HTMLElement} - The container to render the tags in.
     * @param selectable {boolean} - Whether the tags should be selectable or not. Used for the tag selection (default is false)
     * @param isNewPostTagSelect {boolean} - Whether the tags should be selectable for the new post or not (default is false)
     * @returns {void}
     */
    function renderTags(tags, tagContainer, selectable = false, isNewPostTagSelect = false) {
        tags.forEach(tag => {
            const tagItem = document.createElement('div');
            tagItem.classList.add('tag-item');
            tagItem.style.backgroundColor = tag.tag_color;

            if (selectable) {
                if (tags.length === 0) {
                    noTagsMessage.style.display = 'block';
                    return;
                }
                noTagsMessage.style.display = 'none';

                tagItem.classList.add('clickable-tag-item');
                tagItem.style.outlineColor = tag.tag_color;
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.id = `tag-${tag.tag_id}`;
                checkbox.value = tag.tag_name;
                checkbox.dataset.tagId = tag.tag_id;
                tagItem.appendChild(checkbox);
                if (isNewPostTagSelect) {
                    checkbox.classList.add("new-post-tag");
                    tagItem.classList.add("new-post-tag-item");
                    tagItem.addEventListener('click', () => {
                        checkbox.checked = !checkbox.checked;tagItem.classList.toggle('selected');
                        if (checkbox.checked) { // If the tag is selected we add it to the list
                            newPostTags.push(`${tag.tag_id}`);
                        } else { // If the tag is deselected we remove it from the list
                            newPostTags.splice(newPostTags.indexOf(`${tag.tag_id}`), 1);
                        }
                    });
                } else
                    tagItem.addEventListener('click', () => {
                        checkbox.checked = !checkbox.checked;
                        tagItem.classList.toggle('selected');
                        if (checkbox.checked) { // If the tag is selected we add it to the list
                            selectedTags.push(tag.tag_name);
                        } else { // If the tag is deselected we remove it from the list
                            selectedTags.splice(selectedTags.indexOf(tag.tag_name), 1);
                        }
                        postsContainer.innerHTML = "";
                        hasReachedEnd = false;
                        offset = 0;
                        // Reload more messages
                        loadMorePosts();
                    });
            } else {
                tagItem.classList.add("unclickable-tag-item");
            }

            const tagText = document.createElement('span');
            tagText.textContent = tag.tag_name;
            tagText.classList.add('unselectable');
            tagItem.appendChild(tagText);

            tagContainer.appendChild(tagItem);
        });
    }


    /**
     * Create a new post element.
     * @description This function creates a new post element from the given data.
     * @param data {object} - The data of the post to create.
     * @returns {HTMLElement} - The new post element.
     */
    function createNewPost(data) {

        const container = document.createElement("div");
        const postHeader = document.createElement("section");
        const postAuthor = document.createElement("div");
        const authorPfp = document.createElement("img");
        const author = document.createElement("span");
        const option = document.createElement("div");
        const optionButton = document.createElement("button");
        const optionMenu = document.createElement("div");
        const title = document.createElement("span");
        const tags = document.createElement("section");
        const postContent = document.createElement("section");
        const mediaContainer = document.createElement("div");
        const medias = document.createElement("div");
        const postVote = document.createElement("section");
        const voteField = document.createElement("div");
        const upvoteButton = document.createElement("button");
        const upvoteImg = document.createElement("img");
        const vote = document.createElement("span");
        const downvoteButton = document.createElement("button");
        const downvoteImg = document.createElement("img");
        const messageID = document.createElement("p");
        const message = document.createElement("p");
        const wasEdited = document.createElement("p");
        const date = document.createElement("p");
        const voteState = document.createElement("p");
        const br = document.createElement("br");

        let currentVoteState = data.vote_state;
        let currentVoteCount = data.up_votes - data.down_votes;
        let isPostOwner = data.user_name === document.getElementById("username").textContent;

        container.classList.add("post-box", "win95-border");

        postHeader.classList.add("post-header", "win95-header");
        container.appendChild(postHeader);

        postAuthor.classList.add('post-profile');
        postHeader.appendChild(postAuthor);

        authorPfp.src = `/upload/${data.user_pfp_address}`;
        authorPfp.alt = "Author profile picture";
        authorPfp.classList.add("post-profile-picture", "unselectable");
        authorPfp.draggable = false;
        postAuthor.appendChild(authorPfp);

        author.innerText = `${data.user_name}`;
        postAuthor.appendChild(author);

        option.classList.add();
        postHeader.appendChild(option)

        optionButton.innerText = "...";
        optionButton.type = "button";
        optionButton.classList.add("win95-button");
        if (!userIsAuthenticated) {
            optionButton.disabled = true;
        }
        option.appendChild(optionButton);

        optionMenu.classList.add("option-menu", "win95-border");

        let optionMenuReportButtonHTML = `
            <li class="win95-menu-button message-report menu-button" id="post-report-button-p${data.message_id}">
                <img src="/img/report.png" alt="" class="win95-minor-logo">
                <span>${getI18nText("option-menu-report-button-text")}</span>
            </li>`

        let optionMenuEditButtonHTML = `
            <li class="win95-menu-button message-edit menu-button" id="post-edit-button-p${data.message_id}">
                <img src="/img/edit.png" alt="" class="win95-minor-logo">
                <span>${getI18nText("option-menu-edit-button-text")}</span>
            </li>`

        let optionMenuDeleteButtonHTML = `
            <li class="win95-menu-button message-delete menu-button" id="post-delete-button-p${data.message_id}">
                <img src="/img/delete.png" alt="" class="win95-minor-logo">
                <span>${getI18nText("option-menu-delete-button-text")}</span>
            </li>`

        let optionMenuBanButtonHTML = `
            <li class="win95-menu-button message-ban menu-button" id="post-ban-button-p${data.message_id}">
                <img src="/img/ban.png" alt="" class="win95-minor-logo">
                <span>${getI18nText("option-menu-ban-button-text")}</span>
            </li>`
        let additionalButtonsHTML = "";
        let showReportButton = false;
        let showEditButton = false;
        let showDeleteButton = false;
        let showBanButton = false;

        if (userIsAuthenticated) { // If the user is authenticated he can see the option menu
            if (!isPostOwner) { // If the user is authenticated he can report a post (exept his posts)
                additionalButtonsHTML += optionMenuReportButtonHTML;
                showReportButton = true;
            }
            if (isPostOwner) { // If the user is the owner of the post he can edit it
                additionalButtonsHTML += optionMenuEditButtonHTML;
                showEditButton = true;
            }
            if (isPostOwner || userIsModerator || userIsAdmin || userIsAdmin) { // If the user is the owner of the post or his rank is moderator or higher he can delete it
                additionalButtonsHTML += optionMenuDeleteButtonHTML;
                showDeleteButton = true;
            }
            if ((userIsThreadOwner || userIsAdmin) && !isPostOwner) { // If the user rank is admin or higher he can ban the user (exept himself)
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

        console.log(optionMenu.children);

        // Add the event listener to the edit button
        if (showEditButton) {
            const editButton = optionMenu.querySelector(`#post-edit-button-p${data.message_id}`);
            editButton.addEventListener("click", function() {
                console.log(`Edit button clicked for post ${data.message_id}`);
                showEditMenu(data.message_id, data.message_title, data.message_content, data.media_links);
            });
        }

        // Add the event listener to the delete button
        if (showDeleteButton) {
            const deleteButton = optionMenu.querySelector(`#post-delete-button-p${data.message_id}`);
            deleteButton.addEventListener("click", function() {
                console.log(`Delete button clicked for post ${data.message_id}`);
                const result = deleteMessage(threadName, data.message_id);
                result.then(async (response) => {
                    if (response.ok) {
                        container.remove();
                        console.log("Post deleted successfully");
                    } else {
                        alert("Error while deleting post : " + response.statusText);
                        console.error(response);
                    }
                });
            });
        }

        // Add the event listener to the report button
        if (showReportButton) {
            const reportButton = optionMenu.querySelector(`#post-report-button-p${data.message_id}`);
            reportButton.addEventListener("click", function() {
                console.log(`Report button clicked for post ${data.message_id}`);
                showReportMenu(data.message_id);
            });
        }

        // Add the event listener to the ban button
        if (showBanButton) {
            const banButton = optionMenu.querySelector(`#post-ban-button-p${data.message_id}`);
            banButton.addEventListener("click", function() {
                console.log(`Ban button clicked for post ${data.message_id}`);
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

        title.innerText = `${data.message_title}`;
        title.classList.add("post-title");
        postHeader.appendChild(title);

        tags.classList.add("tag-container");
        container.appendChild(tags);

        postContent.classList.add("post-content");
        container.appendChild(postContent);

        mediaContainer.classList.add("post-media-container");
        postContent.appendChild(mediaContainer);

        medias.classList.add("post-medias");
        if (data.media_links != null) {
            postContent.classList.add("win95-border-indent");
            for (let i = 0; i < data.media_links.length; i++) {
                const media = data.media_links[i];
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

            if (data.media_links.length > 1){
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
                data.media_links.forEach((_, i) => {
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
            } else {
                const onlySlide = medias.querySelector('.post-picture');
                onlySlide.style.display = 'block';
            }
        }

        postVote.classList.add("post-vote");
        container.appendChild(postVote);

        voteField.classList.add("post-vote-field");
        postVote.appendChild(voteField);

        upvoteButton.type = "button";
        upvoteButton.classList.add("win95-button", "post-vote-button");
        if (userIsAuthenticated) {
            upvoteButton.addEventListener("click", function () {
                const messageId = data.message_id.toString();
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
        upvoteImg.classList.add("post-vote-image", "unselectable");
        upvoteImg.draggable = false;
        upvoteButton.appendChild(upvoteImg);

        vote.classList.add("post-vote-value");
        vote.innerText = `${data.up_votes - data.down_votes}`;
        voteField.appendChild(vote);

        downvoteButton.type = "button";
        downvoteButton.classList.add("win95-button", "post-vote-button");
        if (userIsAuthenticated) {
            downvoteButton.addEventListener("click", function () {
                const messageId = data.message_id.toString();
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
        downvoteImg.classList.add("post-vote-image", "unselectable");
        downvoteImg.draggable = false;
        downvoteButton.appendChild(downvoteImg);

        // Make sure to update the vote visual when the post is created so that if the user is already upvoting or downvoting the post, the visual is correct
        updateVoteVisual(currentVoteState, upvoteImg, downvoteImg);
        if (data.message_tags != null) {
            renderTags(data.message_tags, tags, false);
        }
        container.appendChild(tags);

        messageID.classList.add("hidden", "messageId");
        messageID.innerText = `${data.message_id}`;
        container.appendChild(messageID);

        message.innerText = `Message content : ${data.message_content}`;
        container.appendChild(message);
        wasEdited.innerText = data.was_edited ? getI18nText('was-edited') : "";
        container.appendChild(wasEdited);
        //date.innerText = `${timeAgo(data.creation_date)}`;
        container.appendChild(date);

        voteState.innerText = `Vote state : ${data.vote_state}`;
        container.appendChild(voteState);
        container.appendChild(br);

        return container;
    }

    // Add each button its event listener if it exists
    if (leaveButton) {
        leaveButton.addEventListener("click", function() {
            const result = leaveThread(getCurrentThreadName());
            result.then(async (response) => {
                if (response.ok) {
                    joinButton.classList.remove("hidden");
                    leaveButton.classList.add("hidden");
                    newPostContainer.classList.add("hidden");
                    console.log("You have left the thread");
                } else {
                    console.error(response);
                }
            });
        });
    }
    if (joinButton) {
        joinButton.addEventListener("click", function() {
            const result = joinThread(getCurrentThreadName())
            result.then(async (response) => {
                if (response.ok) {
                    joinButton.classList.add("hidden");
                    leaveButton.classList.remove("hidden");
                    newPostContainer.classList.remove("hidden");
                    console.log("You have joined the thread");
                } else {
                    console.error(response);
                }
            });
        });
    }
    if (editThreadButton) {
        editThreadButton.addEventListener("click", function() {
            window.location = `/t/${getCurrentThreadName()}/edit`;
        });
    }



    loadMorePostsButton.addEventListener('click' , function() {
        loadMorePosts();
    });

    orderSelect.addEventListener('change' , function() {
       // Empty the posts container
        postsContainer.innerHTML = "";
        hasReachedEnd = false;
        offset = 0;
        // Reload more messages
        loadMorePosts();
    });

    /**
     * Load tags from the server.
     * @description This function sends a request to get the tags for the current thread.
     * @param threadName {string} - The name of the thread to get tags for.
     * @returns {void}
     */
    function loadTags(threadName) {
        getThreadTags(threadName)
            .then(response => {
                if (!response.ok) throw new Error('Failed to load tags.');
                return response.json();
            })
            .then(data => {
                if (!data || data.length === 0) {
                    renderTags([]);
                } else {
                    renderTags(data, tagListContainer, true);
                    renderTags(data, newPostTagsContainer, true, true)
                    pageTags = data;
                }
            })
            .catch(err => {
                console.error('Failed to load tags:', err);
                noTagsMessage.style.display = 'block';
            });
    }

    // New post buttons functions

    function updateNewPostButton() {
        newPostButton.disabled = !(titleValid && contentValid);
    }

    // Display the maximum number of characters for the content
    // And add the validity check
    newPostContent.addEventListener("input", function() {
        const charCount = newPostContent.value.length;
        newPostContentCharCountValue.innerText = `${charCount}`;
        if (charCount > 500) {
            newPostContent.value = newPostContent.value.substring(0, 500);
            newPostContentCharCountValue.innerText = `500`;
        }
        contentValid = !(charCount < 20 || charCount > 500);
        updateNewPostButton();
    });

    newPostTitle.addEventListener("input" , function() {
        const charCount = newPostTitle.value.length;
        if (charCount > 50) {
            newPostTitle.value = newPostTitle.value.substring(0, 50);
        }
        titleValid = !(charCount < 5 || charCount > 50);
        updateNewPostButton();
    });

    function removeSelfAndChildren(element) {
        element.remove();
    }

    newPostfileInput.addEventListener("change", async (e) => {
        e.preventDefault();

        const res = await UploadImages(newPostfileInput, "message_picture")

        for (const [url, id] of res.results) {
            if (url !== null) {
                const wrapper = document.createElement('div');
                wrapper.classList.add("image-preview");
                wrapper.addEventListener('click', () => { // Add click event to remove image
                    removeSelfAndChildren(wrapper);
                    MediaIDs.splice(MediaIDs.indexOf(id), 1);
                });
                const img = document.createElement('img');
                img.src = "/upload/" + url;
                img.alt = "Image preview";
                img.draggable = false;
                img.classList.add("unselectable");
                MediaIDs.push(id);
                wrapper.appendChild(img);
                const imgContainer = newPostImagesPreview.appendChild(wrapper);
            }
        }
    });

    newPostButton.addEventListener("click", function() {
        if (!userIsAuthenticated) {
            alert("You must be logged in to create a post.");
            return;
        }
        if (!userIsMember) {
            alert("You must be a member of the thread to create a post.");
            return;
        }
        sendMessage(threadName, newPostTitle.value, newPostContent.value, MediaIDs, newPostTags)
            .then(r => {
                if (r.ok) {
                    return r.json();
                } else {
                    throw new Error("Error while sending message");
                }
            })
            .then(
                () => {
                    // Clear the new post fields
                    newPostTitle.value = "";
                    newPostContent.value = "";
                    newPostTags = [];
                    let tagCheckbox;
                    for (const tag of document.getElementsByClassName("new-post-tag-item")) {
                        tagCheckbox = tag.querySelector("input[type='checkbox']");
                        if (tagCheckbox.checked) {
                            tagCheckbox.checked = false;
                            tag.classList.toggle('selected');
                        }
                    }
                    MediaIDs = [];
                    while (newPostImagesPreview.firstChild) {
                        newPostImagesPreview.removeChild(newPostImagesPreview.firstChild);
                    }
                    // Clear the new displayed post to avoid duplicates
                    newPostImagesPreview.innerHTML = "";
                    postsContainer.innerHTML = "";
                    hasReachedEnd = false;
                    offset = 0;
                    // Reload more messages
                    loadMorePosts();
                }
            ).catch(error => {
                alert("An error occurred while sending message");
                console.error("Error:", error);
        });
    });

    // Close the report menu when clicking outside of it
    reportMenuBackground.addEventListener("click", hideReportMenu);
    // Close the report menu when clicking on the close button
    reportMenuCloseButton.addEventListener("click", hideReportMenu);
    // Send the report when clicking on the send button
    reportMenuSendButton.addEventListener("click", function() {
        if (reportMenuSendButton.disabled) {
            return;
        }
        reportMessage(threadName, messageToReport, reportReason.value, reportContent.value)
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

    // Close the edit menu when clicking outside of it
    editMenuBackground.addEventListener("click", hideEditMenu);
    // Close the edit menu when clicking on the close button
    editMenuCloseButton.addEventListener("click", hideEditMenu);
    // Send the edit when clicking on the send button
    editMenuSendButton.addEventListener("click", function() {
        if (editMenuSendButton.disabled) {
            return;
        }
        editMessage(threadName, editedPostID, editMenuNewTitleField.value, editMenuNewContentField.value)
            .then(r => {
                if (r.ok) {
                    hideEditMenu();
                    postsContainer.innerHTML = "";
                    hasReachedEnd = false;
                    offset = 0;
                    // Reload more messages
                    loadMorePosts();
                } else {
                    console.error(r);
                }
            });
    });
    // Update the edit content character count
    editMenuNewContentField.addEventListener("input", function() {
        const charCount = editMenuNewContentField.value.length;
        editMenuNewContentCharCountValue.innerText = `${charCount}`;
        if (charCount > 500) {
            editMenuNewContentField.value = editMenuNewContentField.value.substring(0, 500);
            editMenuNewContentCharCountValue.innerText = `500`;
        }
        editMenuSendButton.disabled = (charCount < 20 || charCount > 500);
    });


    // Init the page
    loadTags(threadName);
    loadMorePosts();
    updateNewPostButton();
    if (!userIsAuthenticated || !userIsMember) {
        newPostContainer.classList.add("hidden");
    }
});
