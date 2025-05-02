document.addEventListener("DOMContentLoaded", function () {

    const leaveButton = document.getElementById("LeaveThreadButton")
    const joinButton = document.getElementById("JoinThreadButton")
    const editThreadButton = document.getElementById("EditThreadButton")

    // Debug button
    const loadMorePostsButton = document.getElementById("load-more-posts-button");

    const threadName = getCurrentThreadName()
    let offset= 0;
    let hasReachedEnd = false;
    let orderSelect = document.getElementById("order")
    const examplePostContainer = document.getElementById("example-post-container")
    const postsContainer = document.getElementById("posts-container");

    /**
     * Load more posts from the server.
     * @description This function sends a request to get more posts from the server.
     * @description Then it adds the posts to the posts container.
     */
    function loadMorePosts() {
        if (hasReachedEnd) {
            return;
        }
        let res = getMessage(threadName, offset, orderSelect.value);
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

    function updateUpvote(currentVoteState, currentVoteCount, upvoteImg, downvoteImg) {
        let state = parseInt(currentVoteState);
        let count = currentVoteCount;
        if (state === 1){
            state = 0;
            count -= 1;
            upvoteImg.src = `/img/upvote_empty.png`
            downvoteImg.src = `/img/downvote_empty.png`
        } else if(state === 0 || state === -1){
            state = 1;
            count += (currentVoteState === -1 ? 2 : 1);
            upvoteImg.src = `/img/upvote.png`
            downvoteImg.src = `/img/downvote_empty.png`
        } else {
            state = 422;
        }
        console.log("Current Vote State: ", state);
        console.log("Current Vote Count: ", count);
        return {state, count};
    }

    function updateDownvote(currentVoteState, currentVoteCount, upvoteImg, downvoteImg) {
        let state = parseInt(currentVoteState);
        let count = currentVoteCount;
        if (state === -1){
            state = 0;
            count += 1;
            upvoteImg.src = `/img/upvote_empty.png`
            downvoteImg.src = `/img/downvote_empty.png`
        } else if(state === 0 || state === 1){
            state = -1;
            count -= (currentVoteState === 1 ? 2 : 1);
            upvoteImg.src = `/img/upvote_empty.png`
            downvoteImg.src = `/img/downvote.png`
        } else {
            state = 422;
        }
        console.log("Current Vote State: ", state);
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
        const optionButton = document.createElement("button");
        const title = document.createElement("span");
        const tags = document.createElement("p");
        const postContent = document.createElement("section");
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

        postHeader.classList.add("post-header");
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

        optionButton.innerText = "...";
        optionButton.type = "button";
        optionButton.classList.add("win95-button");
        postHeader.appendChild(optionButton);

        title.innerText = `${data.message_title}`;
        title.classList.add("post-title");
        postHeader.appendChild(title);

        if (data.message_tags != null) {
            for (let i = 0; i < data.message_tags; i++) {
                const tag = data.message_tags[i];
                var tagElement = document.createElement("span");
                tagElement.innerText = `Tag : ${tag}`;
                tags.appendChild(tagElement);
            }
        }
        postHeader.appendChild(tags);

        postContent.classList.add("post-content");
        container.appendChild(postContent);

        medias.classList.add("post-medias");
        if (data.media_links != null) {
            postContent.classList.add("win95-border-bulge");
            for (let i = 0; i < data.media_links.length; i++) {
                const media = data.media_links[i];
                var mediaElement = document.createElement("img");
                mediaElement.src = `/upload/${media}`;
                mediaElement.alt = `Media[${media}]`;
                mediaElement.draggable = false;
                mediaElement.classList.add("post-picture", "unselectable");
                medias.appendChild(mediaElement);
            }
        }
        postContent.appendChild(medias);

        postVote.classList.add("post-vote");
        container.appendChild(postVote);

        voteField.classList.add("post-vote-field");
        postVote.appendChild(voteField);

        upvoteButton.type = "button";
        upvoteButton.classList.add("win95-button", "post-vote-button");
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
            const {state, count} = updateUpvote(currentVoteState, currentVoteCount, upvoteImg, downvoteImg);
            currentVoteState = state;
            currentVoteCount = count;
            vote.innerText = currentVoteCount;
        });
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
            const {state, count} = updateDownvote(currentVoteState, currentVoteCount, upvoteImg, downvoteImg);
            currentVoteState = state;
            currentVoteCount = count;
            vote.innerText = currentVoteCount;
        });
        voteField.appendChild(downvoteButton);

        downvoteImg.src = `/img/downvote_empty.png`;
        downvoteImg.alt = "Downvote image";
        downvoteImg.classList.add("post-vote-image", "unselectable");
        downvoteImg.draggable = false;
        downvoteButton.appendChild(downvoteImg);

        messageID.classList.add("hidden", "messageId");
        messageID.innerText = `Message ID : ${data.message_id}`;
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
});
