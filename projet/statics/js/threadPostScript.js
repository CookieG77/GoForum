document.addEventListener("DOMContentLoaded", function (){
    const loadMoreCommentsButton = document.getElementById("load-more-comments-button");
    const postMenu = document.getElementById("t-post-menu");
    const optionMenu = document.getElementById("t-option-menu");
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
    let userIsThreadMember = document.getElementById("isAMember").textContent === "true";
    let offset= 0;
    let hasReachedEnd = false;
    const commentsContainer = document.getElementById("comments-container");

    let newCommentButton;
    let newCommentContent;
    let newCommentContentCharCountValue;

    if (userIsAuthenticated && userIsThreadMember) {
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

    if (postMenu && optionMenu) {
        postMenu.addEventListener("click", function(e) {
            if (!optionMenu.classList.contains("active")) {
                const options = document.querySelectorAll(".option-menu");
                e.stopPropagation();
                options.forEach(opt => opt.classList.remove("active"));
                optionMenu.classList.add("active");
            }
        });

        window.addEventListener("click", function(e) {
            if (!optionMenu.contains(e.target)) {
                optionMenu.classList.remove("active");
            }
        });
    } else {
        console.log("elements do not exist");
    }

    mediaContainer.classList.add("post-media-container");
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
        const dateAndEditContainer = document.createElement("div");
        const dateSpan = document.createElement("span");
        const isEditedSpan = document.createElement("span");

        let currentVoteState = data.vote_state;
        let currentVoteCount = data.up_votes - data.down_votes;

        container.classList.add("comment-box", "win95-border");

        commentHeader.classList.add("comment-header", "win95-header");
        container.appendChild(commentHeader);

        commentAuthor.classList.add("comment-profile");
        commentHeader.appendChild(commentAuthor);

        authorPfp.src = `/upload/${data.user_pfp_address}`;
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

    if (userIsAuthenticated && userIsThreadMember) {
        postVoteUpButton.addEventListener('click', function() {
            const postId = messageId.toString();
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
            const postId = messageId.toString();
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

    loadMoreComments()
    // Update the vote count
    postVoteCountSpan.innerText = `${postVoteCount}`;
    // Update the post date
    postDateContainer.innerText = timeAgo(postDate);
});

