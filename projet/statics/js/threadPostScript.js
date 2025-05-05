document.addEventListener("DOMContentLoaded", function (){
    const loadMoreCommentsButton = document.getElementById("load-more-posts-button");

    const threadName = getCurrentThreadNameOrPostID()
    let userIsAuthenticated = document.getElementById("isAuthenticated").textContent === "true";
    let offset= 0;
    let hasReachedEnd = false;
    const commentsContainer = document.getElementById("comments-container");

    /**
     * Load more comments from the server.
     * @description This function sends a request to get more comments from the server.
     * @description Then it adds the comments to the comments container.
     */
    function loadMoreComments() {
        if (hasReachedEnd) {
            return;
        }
        let res = getMessage(threadName, offset);
        res.then(async (response) => {
            if (response.ok) {
                const data = await response.json();
                console.log(data);
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

        let currentVoteState = data.vote_state;
        let currentVoteCount = data.up_votes - data.down_votes;

        container.classList.add("comment-box", "win95-border");

        commentHeader.classList.add("comment-header", "win95-header");
        container.appendChild(commentHeader);

        commentAuthor.classList.add("comment-profile");
        commentHeader.appendChild(commentAuthor);

        authorPfp.src = `upload/${data.user_pfp_address}`;
        authorPfp.alt = "Author profile picture";
        authorPfp.classList.add("comment-profile-picture", "unselectable");
        authorPfp.draggable = false;
        commentAuthor.appendChild(authorPfp);

        author.innerText = `${data.user_name}`;
        commentAuthor.appendChild(author);

        option.classList.add();
        commentHeader.appendChild(option);

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

        commentContent.classList.add("comment-content");
        container.appendChild(commentContent);

        commentMedia.classList.add("comment-media");
        commentMedia.innerText = data.comment_content;
        commentContent.appendChild(commentMedia);

        commentVote.classList.add("comment-vote");
        container.appendChild(commentVote);

        voteField.classList.add("comment-vote-field");
        commentVote.appendChild(voteField);

        upvoteButton.type = "button";
        upvoteButton.classList.add("win95-button", "post-vote-button");
        if (userIsAuthenticated) {
            upvoteButton.addEventListener("click", function () {
                const commentId = data.comment_id.toString();
                upvoteComment(threadName, messageId, commentId)
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
                const commentId = data.comment_id.toString();
                downvoteComment(threadName, messageId, commentId)
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

        // Make sure to update the vote visual when the post is created so that if the user is already upvoting or downvoting the post, the visual is correct
        updateVoteVisual(currentVoteState, upvoteImg, downvoteImg);

        downvoteImg.src = `/img/downvote_empty.png`;
        downvoteImg.alt = "Downvote image";
        downvoteImg.classList.add("post-vote-image", "unselectable");
        downvoteImg.draggable = false;
        downvoteButton.appendChild(downvoteImg);

        return container;
    }
})

