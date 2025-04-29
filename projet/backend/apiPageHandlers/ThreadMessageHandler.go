package apiPageHandlers

import (
	f "GoForum/functions"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type jsonMessage struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    []int  `json:"tags,string"`
}

type jsonUpdateMessage struct {
	ID      int    `json:"messageId,string"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type jsonUpdateMessageMedia struct {
	ID      int `json:"messageId,string"`
	MediaID int `json:"mediaId,string"`
}
type jsonMessageDesignator struct {
	ID int `json:"messageId,string"`
}

// ThreadMessageHandler handles the thread message requests from ajax calls
// Its path is /api/thread/{thread}/m/{action}?id={id}
// The "thread" is the name of the thread
// The "action" can be "sendMessage", "deleteMessage", "editMessage", "reportMessage", "upvoteMessage", "downvoteMessage"
// The "id" is the id of the message to edit/delete/report
func ThreadMessageHandler(w http.ResponseWriter, r *http.Request) {
	f.DebugPrintln("ThreadMessageHandler called")

	vars := mux.Vars(r)
	threadName := vars["thread"]
	action := vars["action"]
	// Check if the thread name is empty or does not exist
	if threadName == "" || !f.CheckIfThreadNameExists(threadName) {
		f.DebugPrintf("Thread \"%s\" does not exist\n", threadName)
		http.Error(w, "Thread does not exist or was not specified !", http.StatusNotFound)
		return
	}

	// Check if the action is empty
	if action == "" {
		f.DebugPrintf("Action is empty\n")
		http.Error(w, "Action is empty !", http.StatusBadRequest)
		return
	}

	// Check if the action is a valid action
	if !(action == "sendMessage" ||
		action == "deleteMessage" ||
		action == "removeMedia" ||
		action == "editMessage" ||
		action == "reportMessage" ||
		action == "upvoteMessage" ||
		action == "downvoteMessage" ||
		action == "joinThread" ||
		action == "leaveThread") {

		f.DebugPrintf("Action \"%s\" does not exist\n", action)
		http.Error(w, "Action is empty or does not exist !", http.StatusNotFound)
		return
	}

	// Check if the user is authenticated
	if !f.IsAuthenticated(r) {
		f.DebugPrintf("User is not authenticated\n")
		http.Error(w, "User is not authenticated", http.StatusUnauthorized)
		return
	}

	// Check if the user is verified
	if !f.IsUserVerified(r) {
		f.DebugPrintf("User is not verified\n")
		http.Error(w, "User is not verified", http.StatusUnauthorized)
		return
	}

	thread := f.GetThreadFromName(threadName)
	user := f.GetUser(r)

	// Check if the user is banned from the thread
	userRank := f.GetUserRankInThread(thread, user)
	if userRank < 0 { // If the user is banned from the thread we show him the YOU ARE BANNED page
		f.DebugPrintf("User is banned from the thread he's trying to access\n")
		http.Error(w, "User is banned from the thread", http.StatusForbidden)
		return
	}

	// Execute the action
	switch action {
	case "sendMessage":
		sendMessage(w, r, thread, user)
		return
	case "deleteMessage":
		deleteMessage(w, r, thread, user)
		return
	case "removeMedia":
		removeMedia(w, r, thread, user)
		return
	case "editMessage":
		editMessage(w, r, thread, user)
		return
	case "reportMessage":
		reportMessage(w, r, thread, user)
		return
	case "upvoteMessage":
		upVoteMessage(w, r, thread, user)
		return
	case "downvoteMessage":
		downVoteMessage(w, r, thread, user)
		return
	case "joinThread":
		joinThread(w, r, thread, user)
		return
	case "leaveThread":
		leaveThread(w, r, thread, user)
		return
	default:
		f.DebugPrintf("Action \"%s\" does not exist\n", action)
		http.Error(w, "Action does not exist !", http.StatusNotFound)
		return
	}
}

// sendMessage handles the send message action
// This action is used to send a message to the thread
func sendMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	if r.Method != "POST" {
		f.DebugPrintf("Method is not POST\n")
		http.Error(w, "Method is not POST", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form
	err := r.ParseForm()
	if err != nil {
		f.ErrorPrintf("Error while parsing the form: %v\n", err)
		http.Error(w, "Error while parsing the form", http.StatusBadRequest)
		return
	}

	// Getting the form values
	var msg jsonMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the message title is empty
	if msg.Title == "" {
		f.DebugPrintf("Message title is empty\n")
		http.Error(w, "Message title is empty", http.StatusBadRequest)
		return
	}

	// Check if the message title is valid
	if !f.IsMessageTitleValid(msg.Title) {
		f.DebugPrintf("Message title is not valid\n")
		http.Error(w, "Message title is not valid", http.StatusBadRequest)
		return
	}

	// Check if the message content is empty
	if msg.Content == "" {
		f.DebugPrintf("Message content is empty\n")
		http.Error(w, "Message content is empty", http.StatusBadRequest)
		return
	}

	// Check if the message content is valid
	if !f.IsMessageContentValid(msg.Content) {
		f.DebugPrintf("Message content is not valid\n")
		http.Error(w, "Message content is not valid", http.StatusBadRequest)
		return
	}

	// Check if the user is allowed to send a message
	if !f.IsUserAllowedToSendMessageInThread(thread, user) {
		f.DebugPrintf("User is not allowed to send a message in this thread\n")
		http.Error(w, "User is not allowed to send a message in this thread", http.StatusForbidden)
		return
	}

	// Send the message
	messageID, err := f.AddMessageInThread(thread, msg.Title, msg.Content, user)
	if err != nil {
		f.ErrorPrintf("Error while sending the message: %v\n", err)
		http.Error(w, "Error while sending the message", http.StatusInternalServerError)
		return
	}
	f.DebugPrintf("Message sent with ID: %d\n", messageID)

	// Add the tags to the message
	if len(msg.Tags) > 0 {
		for _, tagID := range msg.Tags {
			isCorrect, err := f.IsTagIDAssociatedWithThread(thread, tagID)
			if !isCorrect {
				f.DebugPrintf("Tag ID %d is not associated with thread %s\n", tagID, thread.ThreadID)
				http.Error(w, "Tag ID is not associated with thread", http.StatusBadRequest)
			} else {
				err = f.AddTagToMessage(messageID, tagID)
				if err != nil {
					f.ErrorPrintf("Error while adding the tag to the message: %v\n", err)
					http.Error(w, "Error while adding the tag to the message", http.StatusInternalServerError)
					return
				}
			}
		}
		f.DebugPrintf("Tags added to message with ID: %d\n", messageID)
	} else {
		f.DebugPrintf("No tags to add to the message with ID: %d\n", messageID)
	}

	// Return the response with the message ID
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success", "messageId":` + strconv.Itoa(messageID) + `}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// deleteMessage handles the delete message action
// This action is used to delete a message
func deleteMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	msgID := _checkMessageApiCallValidity(w, r, thread)
	if msgID < 0 {
		return
	}

	// Check if the user is allowed to delete the message
	if !f.IsUserAllowedToDeleteMessage(thread, user, msgID) {
		f.DebugPrintf("User is not allowed to delete this message\n")
		http.Error(w, "User is not allowed to delete this message", http.StatusForbidden)
		return
	}

	// Delete the message
	err := f.RemoveMessageFromThread(thread, msgID)
	if err != nil {
		f.ErrorPrintf("Error while deleting the message: %v\n", err)
		http.Error(w, "Error while deleting the message", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Message with ID \"%d\" deleted by %s\n", msgID, user.Username)

	// Return the response
	// Return the message ID
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// removeMedia handles the remove media action
// This action is used to remove a media from a message
func removeMedia(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	if r.Method != "POST" {
		f.DebugPrintf("Method is not POST\n")
		http.Error(w, "Method is not POST", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form
	err := r.ParseForm()
	if err != nil {
		f.ErrorPrintf("Error while parsing the form: %v\n", err)
		http.Error(w, "Error while parsing the form", http.StatusBadRequest)
		return
	}

	// Getting the form values
	var msg jsonUpdateMessageMedia
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the message ID is empty
	if msg.ID < 1 {
		f.DebugPrintf("Message ID is empty\n")
		http.Error(w, "Message ID is empty", http.StatusBadRequest)
		return
	}

	// Check if the message ID is valid
	if !f.MessageExistsInThread(thread, msg.ID) {
		f.DebugPrintf("Message ID is not valid\n")
		http.Error(w, "Message ID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the media ID is empty
	if msg.MediaID < 1 {
		f.DebugPrintf("Media ID is empty\n")
		http.Error(w, "Media ID is empty", http.StatusBadRequest)
		return
	}

	// Check if the media ID is valid
	if !f.MediaExistsInMessage(msg.ID, msg.MediaID) {
		f.DebugPrintf("Media ID is not valid\n")
		http.Error(w, "Media ID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the user is allowed to remove the media
	if !f.IsUserAllowedToEditMessageInThread(thread, user, msg.ID) {
		f.DebugPrintf("User is not allowed to remove this media\n")
		http.Error(w, "User is not allowed to remove this media", http.StatusForbidden)
		return
	}

	// Remove the media
	err = f.RemoveMediaLinkFromMessage(msg.ID, msg.MediaID)
	if err != nil {
		f.ErrorPrintf("Error while removing the media: %v\n", err)
		http.Error(w, "Error while removing the media", http.StatusInternalServerError)
		return
	}
	f.DebugPrintf("Media with ID \"%d\" removed from message with ID \"%d\" by %s\n", msg.MediaID, msg.ID, user.Username)
	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// editMessage handles the edit message action
// This action is used to edit a message
func editMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	if r.Method != "POST" {
		f.DebugPrintf("Method is not POST\n")
		http.Error(w, "Method is not POST", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form
	err := r.ParseForm()
	if err != nil {
		f.ErrorPrintf("Error while parsing the form: %v\n", err)
		http.Error(w, "Error while parsing the form", http.StatusBadRequest)
		return
	}

	// Getting the form values
	var msg jsonUpdateMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the message title is empty
	if msg.Title == "" {
		f.DebugPrintf("Message title is empty\n")
		http.Error(w, "Message title is empty", http.StatusBadRequest)
		return
	}

	// Check if the message title is valid
	if !f.IsMessageTitleValid(msg.Title) {
		f.DebugPrintf("Message title is not valid\n")
		http.Error(w, "Message title is not valid", http.StatusBadRequest)
		return
	}

	// Check if the message content is empty
	if msg.Content == "" {
		f.DebugPrintf("Message content is empty\n")
		http.Error(w, "Message content is empty", http.StatusBadRequest)
		return
	}

	// Check if the message content is valid
	if !f.IsMessageContentValid(msg.Content) {
		f.DebugPrintf("Message content is not valid\n")
		http.Error(w, "Message content is not valid", http.StatusBadRequest)
		return
	}

	// Check if the message ID is empty
	if msg.ID < 1 {
		f.DebugPrintf("Message ID is empty\n")
		http.Error(w, "Message ID is empty", http.StatusBadRequest)
		return
	}

	// Check if the message ID is valid
	if !f.MessageExistsInThread(thread, msg.ID) {
		f.DebugPrintf("Message ID is not valid\n")
		http.Error(w, "Message ID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the user is allowed to update the message
	if !f.IsUserAllowedToEditMessageInThread(thread, user, msg.ID) {
		f.DebugPrintf("User is not allowed to update the message in this thread\n")
		http.Error(w, "User is not allowed to update the message in this thread", http.StatusForbidden)
		return
	}
	f.DebugPrintf("new msg data: %v", msg)

	// Send the message
	err = f.EditMessageFromThread(thread, msg.ID, msg.Title, msg.Content)
	if err != nil {
		f.ErrorPrintf("Error while sending the message: %v\n", err)
		http.Error(w, "Error while sending the message", http.StatusInternalServerError)
		return
	}
	f.DebugPrintf("Message with ID \"%d\" was edited by %s\n", msg.ID, user.Username)
	// Return the response with the message ID
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// reportMessage handles the report message action
// This action is used to report a message
func reportMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// TODO : Implement the reportMessage action
	f.DebugPrintf("Report message action not implemented yet\n")
	http.Error(w, "Report message action not implemented yet", http.StatusNotImplemented)
}

// upVoteMessage handles the upvote message action
// This action is used to upvote a message
func upVoteMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	id := _checkMessageApiCallValidity(w, r, thread)
	if id < 0 {
		return
	}
	// Check if the user is allowed to upvote the message (if the user is a member of the thread)
	if !f.IsUserInThread(thread, user) {
		f.DebugPrintf("User is not allowed to downvote this message\n")
		http.Error(w, "User is not allowed to downvote this message", http.StatusForbidden)
		return
	}

	// Check if the user has already up/downvoted the message
	vote := f.HasUserAlreadyVoted(user, id)

	if vote == 0 { // User has not voted yet
		// Add the upvote
		err := f.ThreadMessageUpVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while upvoting the message: %v\n", err)
			http.Error(w, "Error while upvoting the message", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s upvoted message %d\n", user.Username, id)
		return
	}

	if vote == 1 { // User has already upvoted the message so we remove the vote
		err := f.ThreadMessageRemoveVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while removing the vote: %v\n", err)
			http.Error(w, "Error while removing the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s removed his upvote on message %d\n", user.Username, id)
		return
	}

	if vote == -1 { // User has already downvoted the message so we change the vote to upvote
		err := f.ThreadMessageUpdateVote(id, user.UserID, true)
		if err != nil {
			f.ErrorPrintf("Error while updating the vote: %v\n", err)
			http.Error(w, "Error while updating the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s changed his downvote to upvote on message %d\n", user.Username, id)
		return
	}
}

// downVoteMessage handles the downvote message action
// This action is used to downvote a message
func downVoteMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	id := _checkMessageApiCallValidity(w, r, thread)
	if id < 0 {
		return
	}
	// Check if the user is allowed to downvote the message (if the user is a member of the thread)
	if !f.IsUserInThread(thread, user) {
		f.DebugPrintf("User is not allowed to downvote this message\n")
		http.Error(w, "User is not allowed to downvote this message", http.StatusForbidden)
		return
	}

	// Check if the user has already up/downvoted the message
	vote := f.HasUserAlreadyVoted(user, id)

	if vote == 0 { // User has not voted yet
		// Add the downvote
		err := f.ThreadMessageDownVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while downvoting the message: %v\n", err)
			http.Error(w, "Error while downvoting the message", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s downvoted message %d\n", user.Username, id)
		return
	}

	if vote == -1 { // User has already downvoted the message so we remove the vote
		err := f.ThreadMessageRemoveVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while removing the vote: %v\n", err)
			http.Error(w, "Error while removing the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s removed his downvote on message %d\n", user.Username, id)
		return
	}

	if vote == 1 { // User has already upvoted the message so we change the vote to downvote
		err := f.ThreadMessageUpdateVote(id, user.UserID, false)
		if err != nil {
			f.ErrorPrintf("Error while updating the vote: %v\n", err)
			http.Error(w, "Error while updating the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s changed his upvote to downvote on message %d\n", user.Username, id)
		return
	}
}

// joinThread handles the join thread action
// This action is used to join a thread
func joinThread(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	if r.Method != "POST" {
		f.DebugPrintf("Method is not POST\n")
		http.Error(w, "Method is not POST", http.StatusMethodNotAllowed)
		return
	}
	if f.IsUserInThread(thread, user) {
		f.DebugPrintf("User is already in the thread\n")
		http.Error(w, "User is already in the thread", http.StatusBadRequest)
		return
	}

	err := f.JoinThread(thread, user)
	if err != nil {
		f.ErrorPrintf("Error while joining the thread: %v\n", err)
		http.Error(w, "Error while joining the thread", http.StatusInternalServerError)
		return
	}
	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// leaveThread handles the leave thread action
// This action is used to leave a thread
func leaveThread(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	if r.Method != "POST" {
		f.DebugPrintf("Method is not POST\n")
		http.Error(w, "Method is not POST", http.StatusMethodNotAllowed)
		return
	}
	if !f.IsUserInThread(thread, user) {
		f.DebugPrintf("User is not in the thread\n")
		http.Error(w, "User is not in the thread", http.StatusBadRequest)
		return
	}
	// If the user is the owner do not allow him to escape from his responsibilities (prevent him from leaving the thread)
	if f.IsThreadOwner(thread, user) {
		f.DebugPrintf("User is the owner of the thread\n")
		http.Error(w, "User is the owner of the thread", http.StatusBadRequest)
	}

	err := f.LeaveThread(thread, user)
	if err != nil {
		f.ErrorPrintf("Error while leaving the thread: %v\n", err)
		http.Error(w, "Error while leaving the thread", http.StatusInternalServerError)
		return
	}
	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// _checkMessageApiCallValidity checks if the message API call is valid
// It checks if the method is POST, if the form is valid and if the message ID is valid
// It returns the message ID if valid, -1 if not
func _checkMessageApiCallValidity(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum) int {
	if r.Method != "POST" {
		f.DebugPrintf("Method is not POST\n")
		http.Error(w, "Method is not POST", http.StatusMethodNotAllowed)
		return -1
	}

	// Parse the form
	err := r.ParseForm()
	if err != nil {
		f.ErrorPrintf("Error while parsing the form: %v\n", err)
		http.Error(w, "Error while parsing the form", http.StatusBadRequest)
		return -1
	}

	// Getting the form values
	var msg jsonMessageDesignator
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return -1
	}

	// Check if the message ID is empty
	if msg.ID == 0 {
		f.DebugPrintf("Message ID is empty\n")
		http.Error(w, "Message ID is empty", http.StatusBadRequest)
		return -1
	}

	// Check if the message ID is valid
	if !f.MessageExistsInThread(thread, msg.ID) {
		f.DebugPrintf("Message ID is not valid\n")
		http.Error(w, "Message ID is not valid", http.StatusBadRequest)
		return -1
	}
	return msg.ID
}
