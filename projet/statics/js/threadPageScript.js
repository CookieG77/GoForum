document.addEventListener("DOMContentLoaded", function () {

    const leaveButton = document.getElementById("LeaveThreadButton")
    const joinButton = document.getElementById("JoinThreadButton")
    const editThreadButton = document.getElementById("EditThreadButton")

    // Debug button
    const loadMorePostsButton = document.getElementById("load-more-posts-button");

    const threadName = getCurrentThreadName()
    let userIsAuthenticated = document.getElementById("isAuthenticated").textContent === "true";
    let offset= 0;
    let hasReachedEnd = false;
    let selectedTags = [];
    let orderSelect = document.getElementById("order")
    const postsContainer = document.getElementById("posts-container");

    const tagListContainer = document.getElementById('tagList');
    const noTagsMessage = document.getElementById('noTagsMessage');

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
     * @returns {void}
     */
    function renderTags(tags, tagContainer, selectable = false) {
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
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.id = `tag-${tag.tag_id}`;
                checkbox.value = tag.tag_name;
                checkbox.dataset.tagId = tag.tag_id;
                tagItem.appendChild(checkbox);
                tagItem.addEventListener('click', () => {
                    checkbox.checked = !checkbox.checked;
                    tagItem.classList.toggle('selected');
                    if (checkbox.checked) { // If the tag is selected we add it to the list
                        selectedTags.push(tag.tag_name);
                    } else { // If the tag is deselected we remove it from the list
                        selectedTags.splice(selectedTags.indexOf(tag.tag_name), 1);
                    }
                    postsContainer.innerHTML = "";
                    offset = 0;
                    // Reload more messages
                    loadMorePosts();
                });
            }

            const tagText = document.createElement('span');
            tagText.textContent = tag.tag_name;
            tagText.classList.add('unselectable');
            tagItem.appendChild(tagText);

            tagContainer.appendChild(tagItem);
        });
    }

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
        const optionMenu = document.createElement("div")
        const title = document.createElement("span");
        const tags = document.createElement("p");
        const postContent = document.createElement("section");
        const mediaContainer = document.createElement("div")
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
        option.appendChild(optionButton);

        optionMenu.classList.add("option-menu", "win95-border");
        option.appendChild(optionMenu);

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
        wasEdited.innerText = `Was edited : ${data.was_edited}`;
        container.appendChild(wasEdited);
        date.innerText = `Date : ${data.creation_date}`;
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
                    console.log("You have joined the thread");
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
                    console.log("You have left the thread");
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
                }
            })
            .catch(err => {
                console.error('Failed to load tags:', err);
                noTagsMessage.style.display = 'block';
            });
    }

    // Load the tags when the page is loaded
    loadTags(threadName);
    loadMorePosts();
});
