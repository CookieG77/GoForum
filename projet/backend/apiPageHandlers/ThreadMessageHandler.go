package apiPageHandlers

import (
	f "GoForum/functions"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type jsonMessage struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ThreadMessageHandler handles the thread message requests from ajax calls
// It's path is /api/thread/{thread}/{action}?id={id}
// The "thread" is the name of the thread
// The "action" can be "sendMessage", "deleteMessage", "editMessage", "reportMessage", "upvoteMessage", "downvoteMessage"
// The "id" is the id of the message to edit/delete/report
func ThreadMessageHandler(w http.ResponseWriter, r *http.Request) {
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
		action == "editMessage" ||
		action == "reportMessage" ||
		action == "upvoteMessage" ||
		action == "downvoteMessage") {

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

	// Execute the action
	switch action {
	case "sendMessage":
		sendMessage(w, r, thread, user)
		return
	case "deleteMessage":
		deleteMessage(w, r, thread, user)
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
	default:
		f.DebugPrintf("Action \"%s\" does not exist\n", action)
		http.Error(w, "Action does not exist !", http.StatusNotFound)
		return
	}
}

// sendMessage handles the send message action
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
	// Return the message ID
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status":"success"}`))
	if err != nil {
		f.ErrorPrintf("Error while writing the response: %v\n", err)
		http.Error(w, "Error while writing the response", http.StatusInternalServerError)
		return
	}

}

// deleteMessage handles the delete message action
func deleteMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// TODO : Implement the deleteMessage action
	f.DebugPrintf("Delete message action not implemented yet\n")
	http.Error(w, "Delete message action not implemented yet", http.StatusNotImplemented)
}

// editMessage handles the edit message action
func editMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// TODO : Implement the editMessage action
	f.DebugPrintf("Edit message action not implemented yet\n")
	http.Error(w, "Edit message action not implemented yet", http.StatusNotImplemented)
}

// reportMessage handles the report message action
func reportMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// TODO : Implement the reportMessage action
	f.DebugPrintf("Report message action not implemented yet\n")
	http.Error(w, "Report message action not implemented yet", http.StatusNotImplemented)
}

// upVoteMessage handles the upvote message action
func upVoteMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// TODO : Implement the upVoteMessage action
	f.DebugPrintf("Upvote message action not implemented yet\n")
	http.Error(w, "Upvote message action not implemented yet", http.StatusNotImplemented)
}

// downVoteMessage handles the downvote message action
func downVoteMessage(w http.ResponseWriter, r *http.Request, thread f.ThreadGoForum, user f.User) {
	// TODO : Implement the downVoteMessage action
	f.DebugPrintf("Downvote message action not implemented yet\n")
	http.Error(w, "Downvote message action not implemented yet", http.StatusNotImplemented)
}
