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

    /**
     * Create a new post element.
     * @description This function creates a new post element from the given data.
     * @param data {object} - The data of the post to create.
     * @returns {HTMLElement} - The new post element.
     */
    function createNewPost(data) {

        const container = document.createElement("div");
        container.classList.add("post-box", "win95-border");

        const postHeader = document.createElement("section");
        postHeader.classList.add("post-header");
        container.appendChild(postHeader);

        const postAuthor = document.createElement("div");
        postAuthor.classList.add('post-profile');
        postHeader.appendChild(postAuthor);

        const authorPfp = document.createElement("img");
        authorPfp.src = `/upload/${data.user_pfp_address}`;
        authorPfp.alt = "Author profile picture";
        authorPfp.classList.add("post-profile-picture", "unselectable");
        authorPfp.draggable = false;
        postAuthor.appendChild(authorPfp);

        const author = document.createElement("span");
        author.innerText = `${data.user_name}`;
        postAuthor.appendChild(author);

        const optionButton = document.createElement("button");
        optionButton.innerText = "...";
        optionButton.type = "button";
        optionButton.classList.add("win95-button");
        postHeader.appendChild(optionButton);

        const title = document.createElement("span");
        title.innerText = `${data.message_title}`;
        title.classList.add("post-title");
        postHeader.appendChild(title);

        const tags = document.createElement("p");
        if (data.message_tags != null) {
            for (let i = 0; i < data.message_tags; i++) {
                const tag = data.message_tags[i];
                var tagElement = document.createElement("span");
                tagElement.innerText = `Tag : ${tag}`;
                tags.appendChild(tagElement);
            }
        }
        postHeader.appendChild(tags);

        const postContent = document.createElement("section");
        postContent.classList.add("post-content");
        container.appendChild(postContent);

        const medias = document.createElement("div");
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

        const postVote = document.createElement("section");
        postVote.classList.add("post-vote");
        container.appendChild(postVote);

        const upvote = document.createElement("div");
        upvote.classList.add("post-vote-field");
        postVote.appendChild(upvote);

        const upvoteButton = document.createElement("button");
        upvoteButton.type = "button";
        upvoteButton.classList.add("win95-button", "post-vote-button");
        upvote.appendChild(upvoteButton);

        const upvotes = document.createElement("span");
        upvotes.innerText = `${data.up_votes}`;
        upvotes.classList.add("post-vote-value")
        upvote.appendChild(upvotes);

        const upvoteImg = document.createElement("img");
        upvoteImg.src = `/img/upvote_empty.png`;
        upvoteImg.alt = "Upvote image";
        upvoteImg.classList.add("post-vote-image", "unselectable");
        upvoteImg.draggable = false;
        upvoteButton.appendChild(upvoteImg);

        const downvote = document.createElement("div");
        downvote.classList.add("post-vote-field");
        postVote.appendChild(downvote);

        const downvoteButton = document.createElement("button");
        downvoteButton.type = "button";
        downvoteButton.classList.add("win95-button", "post-vote-button");
        downvote.appendChild(downvoteButton);

        const downvotes = document.createElement("span");
        downvotes.innerText = `${data.down_votes}`;
        downvotes.classList.add("post-vote-value")
        downvote.appendChild(downvotes);

        const downvoteImg = document.createElement("img");
        downvoteImg.src = `/img/downvote_empty.png`;
        downvoteImg.alt = "Downvote image";
        downvoteImg.classList.add("post-vote-image", "unselectable");
        downvoteImg.draggable = false;
        downvoteButton.appendChild(downvoteImg);

        const messageID = document.createElement("p");
        messageID.innerText = `Message ID : ${data.message_id}`;
        container.appendChild(messageID);

        const message = document.createElement("p");
        message.innerText = `Message content : ${data.message_content}`;
        container.appendChild(message);
        const wasEdited = document.createElement("p");
        wasEdited.innerText = `Was edited : ${data.was_edited}`;
        container.appendChild(wasEdited);
        const date = document.createElement("p");
        date.innerText = `Date : ${data.creation_date}`;
        container.appendChild(date);

        const voteState = document.createElement("p");
        voteState.innerText = `Vote state : ${data.vote_state}`;
        container.appendChild(voteState);
        const br = document.createElement("br");
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
