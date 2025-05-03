package apiPageHandlers

import (
	f "GoForum/functions"
	"encoding/json"
	"net/http"
	"strconv"
)

func MessageCommentGetter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	threadName := query.Get("thread")
	messageId := query.Get("message")
	offset := query.Get("offset")

	// Check if the thread name is empty or does not exist
	if threadName == "" || !f.CheckIfThreadNameExists(threadName) {
		f.DebugPrintf("Thread \"%s\" does not exist\n", threadName)
		http.Error(w, "Thread does not exist or was not specified !", http.StatusNotFound)
		return
	}
	thread := f.GetThreadFromName(threadName)
	threadConfig := f.GetThreadConfigFromThread(thread)

	// Check if the message MessageID is empty or not a number
	if messageId == "" {
		f.DebugPrintf("Message MessageID is empty\n")
		http.Error(w, "Message MessageID is empty", http.StatusBadRequest)
		return
	}
	messageIdInt, err := strconv.Atoi(messageId)
	if err != nil {
		f.ErrorPrintf("Error parsing message MessageID: %s\n", err)
		http.Error(w, "Message MessageID is not a number", http.StatusBadRequest)
		return
	}

	// Check if the message MessageID exists
	if !f.MessageExistsInThread(thread, messageIdInt) {
		f.DebugPrintf("Message MessageID \"%s\" does not exist\n", messageId)
		http.Error(w, "Message MessageID does not exist", http.StatusNotFound)
		return
	}

	// Check if the offset is empty or not a number
	if offset == "" {
		f.DebugPrintf("Offset is empty\n")
		http.Error(w, "Offset is empty", http.StatusBadRequest)
		return
	}

	// Convert the offset to an int
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		f.ErrorPrintf("Error parsing offset: %s\n", err)
		http.Error(w, "Offset is not a number", http.StatusBadRequest)
		return
	}
	user := f.GetUser(r)

	// Check if the user is banned from the thread the message is in
	userRank := f.GetUserRankInThread(thread, user)
	if userRank < 0 { // If the user is banned from the thread we show him the YOU ARE BANNED page
		f.DebugPrintf("User is banned from the thread he's trying to access\n")
		http.Error(w, "User is banned from the thread", http.StatusForbidden)
		return
	}

	// Check if the user is in the thread
	if !f.IsUserInThread(thread, user) && !threadConfig.IsOpenToNonMembers {
		f.DebugPrintf("User is not in the thread he's trying to access and the thread forbbid non member to access it\n")
		http.Error(w, "User is not in the thread", http.StatusForbidden)
		return
	}

	var Comments []f.FormattedMessageComment
	if user != (f.User{}) {
		Comments, err = f.GetCommentsFromMessageWithPOV(messageIdInt, offsetInt, user)
	} else {
		Comments, err = f.GetCommentsFromMessage(messageIdInt, offsetInt)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Comments)
	if err != nil {
		f.ErrorPrintf("Error encoding comments to JSON: %s\n", err)
		http.Error(w, "Error encoding comments to JSON", http.StatusInternalServerError)
		return
	}
}
