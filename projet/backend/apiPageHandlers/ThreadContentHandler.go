package apiPageHandlers

import (
	f "GoForum/functions"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type jsonMessage struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Medias  IntSlice `json:"medias"`
	Tags    IntSlice `json:"tags"`
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
	MessageID int `json:"messageId,string"`
}

type jsonComment struct {
	MessageID int    `json:"messageId,string"`
	Content   string `json:"content"`
}

type jsonUpdateComment struct {
	CommentID int    `json:"commentId,string"`
	MessageID int    `json:"messageId,string"`
	Content   string `json:"content"`
}

type jsonCommentDesignator struct {
	CommentID int `json:"commentId,string"`
	MessageID int `json:"messageId,string"`
}

// IntSlice is a custom type for handling string-to-int conversion
type IntSlice []int

// jsonReport is a custom type used to handle ajax calls that create a report
type jsonReport struct {
	ID         int    `json:"contentToReportID,string"`
	ReportType string `json:"reportType"`
	Content    string `json:"content"`
}

// jsonReportDesignator is a custom type used to handle ajax calls that target a report
type jsonReportDesignator struct {
	ReportID int `json:"reportId,string"`
}

// jsonUserDesignator is a custom type used to handle ajax calls that target a user
type jsonUserDesignator struct {
	Username string `json:"username"`
}

// jsonThreadTagDesignator is a custom type used to handle ajax calls that target a tag
type jsonThreadTagDesignator struct {
	TagID int `json:"tagId,string"`
}

// jsonThreadTag is a custom type used to handle ajax calls that create a thread tag
type jsonThreadTag struct {
	TagName  string `json:"tagName"`
	TagColor string `json:"tagColor"`
}

type jsonThreadTagUpdate struct {
	TagID    int    `json:"tagId,string"`
	TagName  string `json:"tagName"`
	TagColor string `json:"tagColor"`
}

// ThreadContentHandler handles the thread message requests from ajax calls
// Its path is /api/thread/{thread}/{action}?id={id}
// The "thread" is the name of the thread
// The "action" can be "sendMessage", "deleteMessage", "editMessage", "reportMessage", "upvoteMessage", "downvoteMessage", and other actions...
// The "id" is the id of the message/comment to edit/delete/report
func ThreadContentHandler(w http.ResponseWriter, r *http.Request) {
	f.DebugPrintln("ThreadContentHandler called")

	vars := mux.Vars(r)
	threadName := vars["threadName"]
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
		action == "leaveThread" ||
		action == "sendComment" ||
		action == "deleteComment" ||
		action == "editComment" ||
		action == "reportComment" ||
		action == "upvoteComment" ||
		action == "downvoteComment" ||
		action == "banUser" ||
		action == "setReportToResolved" ||
		action == "createThreadTag" ||
		action == "deleteThreadTag" ||
		action == "getThreadTags") {

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
	case "sendComment":
		sendComment(w, r, thread, user)
		return
	case "deleteComment":
		deleteComment(w, r, thread, user)
		return
	case "editComment":
		editComment(w, r, thread, user)
		return
	case "reportComment":
		reportComment(w, r, thread, user)
		return
	case "upvoteComment":
		upvoteComment(w, r, thread, user)
		return
	case "downvoteComment":
		downvoteComment(w, r, thread, user)
		return
	case "banUser":
		banUser(w, r, thread, user)
		return
	case "setReportToResolved":
		setReportToResolved(w, r, thread, user)
		return
	case "createThreadTag":
		createThreadTag(w, r, thread, user)
		return
	case "deleteThreadTag":
		deleteThreadTag(w, r, thread, user)
		return
	case "editThreadTag":
		editThreadTag(w, r, thread, user)
		return
	case "getThreadTags":
		getThreadTags(w, r, thread, user)
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
	if !f.IsMessageContentOrCommentContentValid(msg.Content) {
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

	// Check if the given media IDs are valid
	if len(msg.Medias) > 0 {
		for _, mediaID := range msg.Medias {
			if !f.NewMediaWithIdIsValid(mediaID) {
				f.DebugPrintf("Media MessageID %d is not valid\n", mediaID)
				http.Error(w, "Media MessageID is not valid", http.StatusBadRequest)
				return
			}
		}
		f.DebugPrintf("All media IDs are valid\n")
	} else {
		f.DebugPrintf("No media IDs provided\n")
	}

	// Send the message
	messageID, err := f.AddMessageInThread(thread, msg.Title, msg.Content, user, msg.Medias...)
	if err != nil {
		f.ErrorPrintf("Error while sending the message: %v\n", err)
		http.Error(w, "Error while sending the message", http.StatusInternalServerError)
		return
	}
	f.DebugPrintf("Message sent with MessageID: %d\n", messageID)

	// Add the tags to the message
	if len(msg.Tags) > 0 {
		for _, tagID := range msg.Tags {
			isCorrect, err := f.IsTagIDAssociatedWithThread(thread, tagID)
			if !isCorrect {
				f.DebugPrintf("Tag MessageID %d is not associated with thread %s\n", tagID, thread.ThreadID)
				http.Error(w, "Tag MessageID is not associated with thread", http.StatusBadRequest)
			} else {
				err = f.AddTagToMessage(messageID, tagID)
				if err != nil {
					f.ErrorPrintf("Error while adding the tag to the message: %v\n", err)
					http.Error(w, "Error while adding the tag to the message", http.StatusInternalServerError)
					return
				}
			}
		}
		f.DebugPrintf("Tags added to message with MessageID: %d\n", messageID)
	} else {
		f.DebugPrintf("No tags to add to the message with MessageID: %d\n", messageID)
	}

	// Return the response with the message MessageID
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

	f.DebugPrintf("Message with MessageID \"%d\" deleted by %s\n", msgID, user.Username)

	// Return the response
	// Return the message MessageID
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

	// Check if the message MessageID is empty
	if msg.ID < 1 {
		f.DebugPrintf("Message MessageID is empty\n")
		http.Error(w, "Message MessageID is empty", http.StatusBadRequest)
		return
	}

	// Check if the message MessageID is valid
	if !f.MessageExistsInThread(thread, msg.ID) {
		f.DebugPrintf("Message MessageID is not valid\n")
		http.Error(w, "Message MessageID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the media MessageID is empty
	if msg.MediaID < 1 {
		f.DebugPrintf("Media MessageID is empty\n")
		http.Error(w, "Media MessageID is empty", http.StatusBadRequest)
		return
	}

	// Check if the media MessageID is valid
	if !f.MediaExistsInMessage(msg.ID, msg.MediaID) {
		f.DebugPrintf("Media MessageID is not valid\n")
		http.Error(w, "Media MessageID is not valid", http.StatusBadRequest)
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
	f.DebugPrintf("Media with MessageID \"%d\" removed from message with MessageID \"%d\" by %s\n", msg.MediaID, msg.ID, user.Username)
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
	if !f.IsMessageContentOrCommentContentValid(msg.Content) {
		f.DebugPrintf("Message content is not valid\n")
		http.Error(w, "Message content is not valid", http.StatusBadRequest)
		return
	}

	// Check if the message MessageID is empty
	if msg.ID < 1 {
		f.DebugPrintf("Message MessageID is empty\n")
		http.Error(w, "Message MessageID is empty", http.StatusBadRequest)
		return
	}

	// Check if the message MessageID is valid
	if !f.MessageExistsInThread(thread, msg.ID) {
		f.DebugPrintf("Message MessageID is not valid\n")
		http.Error(w, "Message MessageID is not valid", http.StatusBadRequest)
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
	f.DebugPrintf("Message with MessageID \"%d\" was edited by %s\n", msg.ID, user.Username)
	// Return the response with the message MessageID
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
	// Check if the user is allowed to report the message
	if !f.IsUserAllowedToSendMessageInThread(thread, user) { // Same check as sendMessage
		f.DebugPrintf("User is not allowed to report this message\n")
		http.Error(w, "User is not allowed to report this message", http.StatusForbidden)
		return
	}

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
	var message jsonReport
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the message MessageID is empty
	if message.ID < 1 {
		f.DebugPrintf("Message MessageID is empty\n")
		http.Error(w, "Message MessageID is empty", http.StatusBadRequest)
		return
	}

	// Check if the message MessageID is valid
	if !f.MessageExistsInThread(thread, message.ID) {
		f.DebugPrintf("Message MessageID is not valid\n")
		http.Error(w, "Message MessageID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the reason is empty
	if message.ReportType == "" {
		f.DebugPrintf("Report reason is empty\n")
		http.Error(w, "Report reason is empty", http.StatusBadRequest)
		return
	}

	// Check if the reason is valid
	if !f.IsAReportType(message.ReportType) {
		f.DebugPrintf("Report reason is not valid\n")
		http.Error(w, "Report reason is not valid", http.StatusBadRequest)
		return
	}

	// Check if the comment is empty
	if message.Content == "" {
		f.DebugPrintf("Report comment is empty\n")
		http.Error(w, "Report comment is empty", http.StatusBadRequest)
		return
	}

	// Check if the comment is valid
	if !f.IsMessageContentOrCommentContentValid(message.Content) { // Same check as IsMessageContentOrCommentContentValid
		f.DebugPrintf("Report comment is not valid\n")
		http.Error(w, "Report comment is not valid", http.StatusBadRequest)
		return
	}

	reportType, _ := f.GetReportTypeFromString(message.ReportType)

	// Send the report
	err = f.AddReportedMessage(user, message.ID, reportType, message.Content)
	if err != nil {
		f.ErrorPrintf("Error while sending the report: %v\n", err)
		http.Error(w, "Error while sending the report", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Message with MessageID \"%d\" was reported by %s\n", message.ID, user.Username)

	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// upVoteMessage handles the upvote message action
// This action is used to upvote a message
func upVoteMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	id := _checkMessageApiCallValidity(w, r, thread)
	if id < 0 {
		return
	}

	// Check if the user has already up/downvoted the message
	vote := f.HasUserAlreadyVotedOnMessage(user, id)

	if vote == 0 { // User has not voted yet
		// Add the upvote
		err := f.ThreadMessageUpVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while upvoting the message: %v\n", err)
			http.Error(w, "Error while upvoting the message", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s upvoted message %d\n", user.Username, id)
	} else if vote == 1 { // User has already upvoted the message so we remove the vote
		err := f.ThreadMessageRemoveVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while removing the vote: %v\n", err)
			http.Error(w, "Error while removing the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s removed his upvote on message %d\n", user.Username, id)
	} else if vote == -1 { // User has already downvoted the message so we change the vote to upvote
		err := f.ThreadMessageUpdateVote(id, user.UserID, true)
		if err != nil {
			f.ErrorPrintf("Error while updating the vote: %v\n", err)
			http.Error(w, "Error while updating the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s changed his downvote to upvote on message %d\n", user.Username, id)
	} else {
		f.DebugPrintf("User %s has an unknown vote value: %d\n", user.Username, vote)
		http.Error(w, "Unknown vote value", http.StatusInternalServerError)
		return
	}
	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
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

	// Check if the user has already up/downvoted the message
	vote := f.HasUserAlreadyVotedOnMessage(user, id)

	if vote == 0 { // User has not voted yet
		// Add the downvote
		err := f.ThreadMessageDownVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while downvoting the message: %v\n", err)
			http.Error(w, "Error while downvoting the message", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s downvoted message %d\n", user.Username, id)
	} else if vote == -1 { // User has already downvoted the message so we remove the vote
		err := f.ThreadMessageRemoveVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while removing the vote: %v\n", err)
			http.Error(w, "Error while removing the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s removed his downvote on message %d\n", user.Username, id)
	} else if vote == 1 { // User has already upvoted the message so we change the vote to downvote
		err := f.ThreadMessageUpdateVote(id, user.UserID, false)
		if err != nil {
			f.ErrorPrintf("Error while updating the vote: %v\n", err)
			http.Error(w, "Error while updating the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s changed his upvote to downvote on message %d\n", user.Username, id)
	} else {
		f.DebugPrintf("User %s has an unknown vote value: %d\n", user.Username, vote)
		http.Error(w, "Unknown vote value", http.StatusInternalServerError)
		return
	}
	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
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

// sendComment handles the edit comment action
// This action is used to send a comment
func sendComment(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// Check if the method is POST
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
	var comment jsonComment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&comment); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the comment MessageID is empty
	if comment.MessageID < 1 {
		f.DebugPrintf("Content MessageID is empty\n")
		http.Error(w, "Content MessageID is empty", http.StatusBadRequest)
		return
	}

	// Check if the comment MessageID is valid
	if !f.MessageExistsInThread(thread, comment.MessageID) {
		f.DebugPrintf("Content MessageID is not valid\n")
		http.Error(w, "Content MessageID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the comment content is empty
	if comment.Content == "" {
		f.DebugPrintf("Content content is empty\n")
		http.Error(w, "Content content is empty", http.StatusBadRequest)
		return
	}

	// Check if the comment content is valid
	if !f.IsMessageContentOrCommentContentValid(comment.Content) {
		f.DebugPrintf("Content content is not valid\n")
		http.Error(w, "Content content is not valid", http.StatusBadRequest)
		return
	}

	// Check if the user is allowed to send a comment
	if !f.IsUserAllowedToSendMessageInThread(thread, user) {
		f.DebugPrintf("User is not allowed to send a comment in this thread\n")
		http.Error(w, "User is not allowed to send a comment in this thread", http.StatusForbidden)
		return
	}

	// Send the comment
	commentID, err := f.AddCommentToPost(user, comment.MessageID, comment.Content)
	if err != nil {
		f.ErrorPrintf("Error while sending the comment: %v\n", err)
		http.Error(w, "Error while sending the comment", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Content sent with CommentID '%d' on message '%d'\n", commentID, comment.MessageID)
	// Return the response with the comment CommentID
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success", "commentId":` + strconv.Itoa(commentID) + `}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// editComment handles the edit comment action
// This action is used to edit a comment
func editComment(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// Check if the method is POST
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
	var comment jsonUpdateComment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&comment); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the comment MessageID is empty
	if comment.MessageID < 1 {
		f.DebugPrintf("Content MessageID is empty\n")
		http.Error(w, "Content MessageID is empty", http.StatusBadRequest)
		return
	}

	// Check if the comment MessageID is valid
	if !f.MessageExistsInThread(thread, comment.MessageID) {
		f.DebugPrintf("Content MessageID is not valid\n")
		http.Error(w, "Content MessageID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the comment CommentID is empty
	if comment.CommentID < 1 {
		f.DebugPrintf("Content CommentID is empty\n")
		http.Error(w, "Content CommentID is empty", http.StatusBadRequest)
		return
	}

	// Check if the comment CommentID is valid
	if !f.CommentExistsOnMessage(comment.MessageID, comment.CommentID) {
		f.DebugPrintf("Content CommentID is not valid\n")
		http.Error(w, "Content CommentID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the user is allowed to edit the comment
	if !f.IsUserAllowedToEditComment(thread, user, comment.CommentID) {
		f.DebugPrintf("User is not allowed to edit this comment\n")
		http.Error(w, "User is not allowed to edit this comment", http.StatusForbidden)
		return
	}

	// Update the comment
	err = f.EditCommentFromPost(comment.CommentID, comment.Content)
	if err != nil {
		f.ErrorPrintf("Error while updating the comment: %v\n", err)
		http.Error(w, "Error while updating the comment", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Content with CommentID \"%d\" was edited by %s\n", comment.CommentID, user.Username)

	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// deleteComment handles the delete comment action
// This action is used to delete a comment
func deleteComment(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	commentID := _checkCommentApiCallValidity(w, r, thread)
	if commentID <= 0 {
		return
	}
	// Check if the user is allowed to delete the comment
	if !f.IsUserAllowedToDeleteComment(thread, user, commentID) {
		f.DebugPrintf("User is not allowed to delete this comment\n")
		http.Error(w, "User is not allowed to delete this comment", http.StatusForbidden)
		return
	}

	// Delete the comment
	err := f.RemoveCommentFromPost(commentID)
	if err != nil {
		f.ErrorPrintf("Error while deleting the comment: %v\n", err)
		http.Error(w, "Error while deleting the comment", http.StatusInternalServerError)
		return
	}
}

// reportComment handles the report comment action
// This action is used to report a comment
func reportComment(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// Check if the user is allowed to report the comment
	if !f.IsUserAllowedToSendMessageInThread(thread, user) { // Same check as sendMessage
		f.DebugPrintf("User is not allowed to report this comment\n")
		http.Error(w, "User is not allowed to report this comment", http.StatusForbidden)
		return
	}

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
	var comment jsonReport
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&comment); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the message MessageID is empty
	if comment.ID < 1 {
		f.DebugPrintf("comment CommentID is empty\n")
		http.Error(w, "comment CommentID is empty", http.StatusBadRequest)
		return
	}

	// Check if the message MessageID is valid
	if !f.MessageExistsInThread(thread, comment.ID) {
		f.DebugPrintf("comment CommentID is not valid\n")
		http.Error(w, "comment CommentID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the reason is empty
	if comment.ReportType == "" {
		f.DebugPrintf("Report reason is empty\n")
		http.Error(w, "Report reason is empty", http.StatusBadRequest)
		return
	}

	// Check if the reason is valid
	if !f.IsAReportType(comment.ReportType) {
		f.DebugPrintf("Report reason is not valid\n")
		http.Error(w, "Report reason is not valid", http.StatusBadRequest)
		return
	}

	// Check if the comment is empty
	if comment.Content == "" {
		f.DebugPrintf("Report comment is empty\n")
		http.Error(w, "Report comment is empty", http.StatusBadRequest)
		return
	}

	// Check if the comment is valid
	if !f.IsMessageContentOrCommentContentValid(comment.Content) { // Same check as IsMessageContentOrCommentContentValid
		f.DebugPrintf("Report comment is not valid\n")
		http.Error(w, "Report comment is not valid", http.StatusBadRequest)
		return
	}

	reportType, _ := f.GetReportTypeFromString(comment.ReportType)

	// Send the report
	err = f.AddReportedComment(user, comment.ID, reportType, comment.Content)
	if err != nil {
		f.ErrorPrintf("Error while sending the report: %v\n", err)
		http.Error(w, "Error while sending the report", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Message with MessageID \"%d\" was reported by %s\n", comment.ID, user.Username)

	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// upvoteComment handles the upvote comment action
// This action is used to upvote a comment
func upvoteComment(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	id := _checkCommentApiCallValidity(w, r, thread)
	if id <= 0 {
		return
	}

	// Check if the user has already up/downvoted the comment
	vote := f.HasUserAlreadyVotedOnComment(user, id)

	if vote == 0 { // User has not voted yet
		// Add the upvote
		err := f.MessageCommentUpVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while upvoting the comment: %v\n", err)
			http.Error(w, "Error while upvoting the comment", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s upvoted comment %d\n", user.Username, id)
	} else if vote == 1 { // User has already upvoted the comment so we remove the vote
		err := f.MessageCommentRemoveVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while removing the vote: %v\n", err)
			http.Error(w, "Error while removing the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s removed his upvote on comment %d\n", user.Username, id)
	} else if vote == -1 { // User has already downvoted the comment so we change the vote to upvote
		err := f.MessageCommentUpdateVote(id, user.UserID, true)
		if err != nil {
			f.ErrorPrintf("Error while updating the vote: %v\n", err)
			http.Error(w, "Error while updating the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s changed his downvote to upvote on comment %d\n", user.Username, id)
	} else {
		f.DebugPrintf("User %s has an unknown vote value: %d\n", user.Username, vote)
		http.Error(w, "Unknown vote value", http.StatusInternalServerError)
		return
	}
	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// downvoteComment handles the downvote comment action
// This action is used to downvote a comment
func downvoteComment(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	id := _checkCommentApiCallValidity(w, r, thread)
	if id <= 0 {
		return
	}

	// Check if the user has already up/downvoted the comment
	vote := f.HasUserAlreadyVotedOnComment(user, id)

	if vote == 0 { // User has not voted yet
		// Add the downvote
		err := f.MessageCommentDownVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while downvoting the comment: %v\n", err)
			http.Error(w, "Error while downvoting the comment", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s downvoted comment %d\n", user.Username, id)
	} else if vote == -1 { // User has already downvoted the comment so we remove the vote
		err := f.MessageCommentRemoveVote(id, user.UserID)
		if err != nil {
			f.ErrorPrintf("Error while removing the vote: %v\n", err)
			http.Error(w, "Error while removing the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s removed his downvote on comment %d\n", user.Username, id)
	} else if vote == 1 { // User has already upvoted the comment so we change the vote to downvote
		err := f.MessageCommentUpdateVote(id, user.UserID, false)
		if err != nil {
			f.ErrorPrintf("Error while updating the vote: %v\n", err)
			http.Error(w, "Error while updating the vote", http.StatusInternalServerError)
			return
		}
		f.DebugPrintf("User %s changed his upvote to downvote on comment %d\n", user.Username, id)
	} else {
		f.DebugPrintf("User %s has an unknown vote value: %d\n", user.Username, vote)
		http.Error(w, "Unknown vote value", http.StatusInternalServerError)
		return
	}
	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

// banUser handles the ban user action
// This action is used to ban a user from the thread
func banUser(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// Check if the user is allowed to ban a user
	if !f.IsUserAllowedToBanUserInThread(thread, user) {
		f.DebugPrintf("User is not allowed to ban a user in this thread\n")
		http.Error(w, "User is not allowed to ban a user in this thread", http.StatusForbidden)
		return
	}
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
	var msg jsonUserDesignator
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the user username is empty
	if msg.Username < "" {
		f.DebugPrintf("User MessageID is empty\n")
		http.Error(w, "User MessageID is empty", http.StatusBadRequest)
		return
	}

	// Check if the user username is valid
	userToBan, err := f.GetUserFromUsername(msg.Username)
	if err != nil {
		f.DebugPrintf("User MessageID is not valid\n")
		http.Error(w, "User MessageID is not valid", http.StatusBadRequest)
		return
	}
	if (userToBan == f.User{}) {
		f.DebugPrintf("User MessageID is not valid\n")
		http.Error(w, "User MessageID is not valid", http.StatusBadRequest)
		return
	}

	// Check if the user is already banned from the thread
	if f.IsUserBannedFromThread(thread, userToBan) {
		f.DebugPrintf("User is already banned from the thread\n")
		http.Error(w, "User is already banned from the thread", http.StatusBadRequest)
		return
	}

	// Ban the user
	err = f.BanUserFromThread(thread, userToBan)
	if err != nil {
		f.ErrorPrintf("Error while banning the user: %v\n", err)
		http.Error(w, "Error while banning the user", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("User %s was banned from thread %s by %s\n", userToBan.Username, thread.ThreadID, user.Username)

	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

func setReportToResolved(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// Check if the user is allowed to set the report to resolved
	if !f.IsUserAllowedToBanUserInThread(thread, user) { // Same check as banUser
		f.DebugPrintf("User is not allowed to set the report to resolved in this thread\n")
		http.Error(w, "User is not allowed to set the report to resolved in this thread", http.StatusForbidden)
		return
	}
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
	var report jsonReportDesignator
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&report); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the report ReportID is empty
	if report.ReportID < 1 {
		f.DebugPrintf("Report ReportID is empty\n")
		http.Error(w, "Report ReportID is empty", http.StatusBadRequest)
		return
	}

	// Check if the report ReportID is valid
	if !f.ReportExistsInThread(thread, report.ReportID) {
		f.DebugPrintf("Report ReportID is not valid\n")
		http.Error(w, "Report ReportID is not valid", http.StatusBadRequest)
		return
	}

	// Set the report to resolved
	err = f.SetReportAsResolved(report.ReportID)
	if err != nil {
		f.ErrorPrintf("Error while setting the report as resolved: %v\n", err)
		http.Error(w, "Error while setting the report as resolved", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Report %d was set to resolved by %s\n", report.ReportID, user.Username)

	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

func createThreadTag(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	if !(f.GetUserRankInThread(thread, user) >= f.ThreadRankOwner) {
		f.DebugPrintf("User is not allowed to create a tag in this thread\n")
		http.Error(w, "User is not allowed to create a tag in this thread", http.StatusForbidden)
		return
	}

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
	var tag jsonThreadTag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tag); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the tag name is empty
	if tag.TagName == "" {
		f.DebugPrintf("Tag name is empty\n")
		http.Error(w, "Tag name is empty", http.StatusBadRequest)
		return
	}

	// Check if the tag name is valid
	if !f.IsTagNameValid(tag.TagName) {
		f.DebugPrintf("Tag name is not valid\n")
		http.Error(w, "Tag name is not valid", http.StatusBadRequest)
		return
	}

	// Check if the tag color is empty
	if tag.TagColor == "" {
		f.DebugPrintf("Tag color is empty\n")
		http.Error(w, "Tag color is empty", http.StatusBadRequest)
		return
	}

	// Check if the tag color is valid
	if !f.IsStringHexColor(tag.TagColor) {
		f.DebugPrintf("Tag color is not valid\n")
		http.Error(w, "Tag color is not valid", http.StatusBadRequest)
		return
	}

	// Check if the tag is already in the thread
	res, err := f.TagAlreadyExists(thread, tag.TagName)
	if err != nil {
		f.ErrorPrintf("Error while checking if the tag already exists: %v\n", err)
		http.Error(w, "Error while checking if the tag already exists", http.StatusInternalServerError)
		return
	}
	if res {
		f.DebugPrintf("Tag already exists in the thread\n")
		http.Error(w, "Tag already exists in the thread", http.StatusBadRequest)
		return
	}

	// Create the tag
	err = f.AddThreadTag(thread, tag.TagName, tag.TagColor)
	if err != nil {
		f.ErrorPrintf("Error while creating the tag: %v\n", err)
		http.Error(w, "Error while creating the tag", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Tag %s was created in thread %s by %s\n", tag.TagName, thread.ThreadID, user.Username)
	// Return the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

func deleteThreadTag(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	if !(f.GetUserRankInThread(thread, user) >= f.ThreadRankOwner) {
		f.DebugPrintf("User is not allowed to delete a tag in this thread\n")
		http.Error(w, "User is not allowed to delete a tag in this thread", http.StatusForbidden)
		return
	}

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
	var tag jsonThreadTagDesignator
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tag); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	// Check if the tag id is empty
	if tag.TagID == 0 {
		f.DebugPrintf("Tag id is empty\n")
		http.Error(w, "Tag id is empty", http.StatusBadRequest)
		return
	}

	// Get the full tag
	tagFull, err := f.GetTagByID(tag.TagID)

	err = f.RemoveThreadTag(thread, tagFull.TagName)
	if err != nil {
		f.ErrorPrintf("Error while deleting the tag: %v\n", err)
		http.Error(w, "Error while deleting the tag", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Tag %d was deleted in thread %s by %s\n", tagFull.TagName, thread.ThreadID, user.Username)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

func editThreadTag(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	if !(f.GetUserRankInThread(thread, user) >= f.ThreadRankOwner) {
		f.DebugPrintf("User is not allowed to edit a tag in this thread\n")
		http.Error(w, "User is not allowed to edit a tag in this thread", http.StatusForbidden)
		return
	}

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
	var tag jsonThreadTagUpdate
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tag); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return
	}

	err = f.UpdateThreadTag(thread, tag.TagID, tag.TagName, tag.TagColor)
	if err != nil {
		f.ErrorPrintf("Error while editing the tag: %v\n", err)
		http.Error(w, "Error while editing the tag", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Tag %d was edited in thread %s by %s\n", tag.TagName, thread.ThreadID, user.Username)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}
}

func getThreadTags(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
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

	tags, err := f.GetThreadTags(thread)
	if err != nil {
		f.ErrorPrintf("Error while getting the thread tags: %v\n", err)
		http.Error(w, "Error while getting the thread tags", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tags)
	if err != nil {
		f.ErrorPrintf("Error encoding comments to JSON: %s\n", err)
		http.Error(w, "Error encoding comments to JSON", http.StatusInternalServerError)
		return
	}
}

// _checkMessageApiCallValidity checks if the message API call is valid
// It checks if the method is POST, if the form is valid and if the message MessageID is valid
// It returns the message MessageID if valid, -1 if not
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

	// Check if the message MessageID is empty
	if msg.MessageID == 0 {
		f.DebugPrintf("Message MessageID is empty\n")
		http.Error(w, "Message MessageID is empty", http.StatusBadRequest)
		return -1
	}

	// Check if the message MessageID is valid
	if !f.MessageExistsInThread(thread, msg.MessageID) {
		f.DebugPrintf("Message MessageID is not valid\n")
		http.Error(w, "Message MessageID is not valid", http.StatusBadRequest)
		return -1
	}
	return msg.MessageID
}

// _checkCommentApiCallValidity checks if the comment API call is valid
// It checks if the method is POST, if the form is valid and if the comment CommentID and MessageID is valid
// It returns the comment CommentID if valid, -1 if not
func _checkCommentApiCallValidity(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum) int {
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
	var msg jsonCommentDesignator
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		f.ErrorPrintf("Error while decoding the JSON: %v\n", err)
		http.Error(w, "Error while decoding the JSON", http.StatusBadRequest)
		return -1
	}

	// Check if the message CommentID is empty
	if msg.CommentID == 0 {
		f.DebugPrintf("Message CommentID is empty\n")
		http.Error(w, "Message CommentID is empty", http.StatusBadRequest)
		return -1
	}

	// Check if the message MessageID is empty
	if msg.MessageID == 0 {
		f.DebugPrintf("Message MessageID is empty\n")
		http.Error(w, "Message MessageID is empty", http.StatusBadRequest)
		return -1
	}

	// Check if the message MessageID is valid
	if !f.CommentExistsOnMessage(msg.MessageID, msg.CommentID) {
		f.DebugPrintf("Content ID or Message ID is not valid\n")
		http.Error(w, "Content ID or Message ID is not valid", http.StatusBadRequest)
		return -1
	}
	return msg.CommentID
}

// UnmarshalJSON implement the json.Unmarshaler interface for IntSlice
func (s *IntSlice) UnmarshalJSON(data []byte) error {
	var raw []string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal as []string: %w", err)
	}

	var result []int
	for _, str := range raw {
		num, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("failed to convert string to int: %w", err)
		}
		result = append(result, num)
	}

	*s = result
	return nil
}
