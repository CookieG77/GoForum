package apiPageHandlers

import (
	f "GoForum/functions"
	"encoding/json"
	"net/http"
)

func ThreadTagsGetterHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	threadName := query.Get("thread")

	// Check if the thread name is empty
	if threadName == "" {
		f.DebugPrintf("No thread name specified\n")
		http.Error(w, "No thread name specified", http.StatusBadRequest)
		return
	}

	// Check if the thread name exists
	if !f.CheckIfThreadNameExists(threadName) {
		f.DebugPrintf("Thread \"%s\" does not exist\n", threadName)
		http.Error(w, "Thread does not exist or was not specified !", http.StatusNotFound)
		return
	}

	thread := f.GetThreadFromName(threadName)
	threadConfig := f.GetThreadConfigFromThread(thread)
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

	// Get the tags of the thread
	tags, err := f.GetThreadTags(thread)
	if err != nil {
		http.Error(w, "Error retrieving tags", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("Returning thread %s tags!\n", threadName)

	// Return the tags as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tags)
	if err != nil {
		f.ErrorPrintf("Error encoding tags to JSON: %s\n", err)
		http.Error(w, "Error encoding tags to JSON", http.StatusInternalServerError)
		return
	}
}
