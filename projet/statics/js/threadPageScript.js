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
        const messageID = document.createElement("p");
        messageID.innerText = `Message ID : ${data.message_id}`;
        container.appendChild(messageID);
        const title = document.createElement("h2");
        title.innerText = `Message title : ${data.message_title}`;
        container.appendChild(title);
        const message = document.createElement("p");
        message.innerText = `Message content : ${data.message_content}`;
        container.appendChild(message);
        const wasEdited = document.createElement("p");
        wasEdited.innerText = `Was edited : ${data.was_edited}`;
        container.appendChild(wasEdited);
        const date = document.createElement("p");
        date.innerText = `Date : ${data.creation_date}`;
        container.appendChild(date);
        const author = document.createElement("p");
        author.innerText = `Author : ${data.user_name}`;
        container.appendChild(author);
        const authorPfp = document.createElement("img");
        authorPfp.src = `/upload/${data.user_pfp_address}`;
        authorPfp.alt = "Author profile picture";
        container.appendChild(authorPfp);
        const upvotes = document.createElement("p");
        upvotes.innerText = `Upvotes : ${data.up_votes}`;
        container.appendChild(upvotes);
        const downvotes = document.createElement("p");
        downvotes.innerText = `Downvotes : ${data.down_votes}`;
        container.appendChild(downvotes);
        const medias = document.createElement("p");
        if (data.media_links != null) {
            for (let i = 0; i < data.media_links.length; i++) {
                const media = data.media_links[i];
                var mediaElement = document.createElement("img");
                mediaElement.src = `/upload/${media}`;
                mediaElement.alt = `Media[${media}]`;
                medias.appendChild(mediaElement);
            }
        }
        container.appendChild(medias);
        const tags = document.createElement("p");
        if (data.message_tags != null) {
            for (let i = 0; i < data.message_tags; i++) {
                const tag = data.message_tags[i];
                var tagElement = document.createElement("span");
                tagElement.innerText = `Tag : ${tag}`;
                tags.appendChild(tagElement);
            }
        }
        container.appendChild(tags);
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
