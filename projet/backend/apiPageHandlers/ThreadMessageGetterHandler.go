package apiPageHandlers

import (
	f "GoForum/functions"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
)

func ThreadMessageGetter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	threadName := query.Get("thread")
	offset := query.Get("offset")
	order := query.Get("order")
	tags := query.Get("tags")

	// Check if the thread name is empty or does not exist
	if threadName == "" || !f.CheckIfThreadNameExists(threadName) {
		f.DebugPrintf("Thread \"%s\" does not exist\n", threadName)
		http.Error(w, "Thread does not exist or was not specified !", http.StatusNotFound)
		return
	}
	// Check if the offset is empty or not a number
	if offset == "" {
		f.DebugPrintf("Offset is empty\n")
		http.Error(w, "Offset is empty", http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(offset); err != nil {
		f.ErrorPrintf("Error parsing offset: %s\n", err)
		http.Error(w, "Offset is not a number", http.StatusBadRequest)
		return
	}

	// Check if the order is empty or not a valid order
	if order == "" || !slices.Contains(f.OrderingList, order) {
		f.DebugPrintf("Order is empty or not a valid order\n")
		http.Error(w, "Order is empty or not a valid order", http.StatusBadRequest)
		return
	}

	// Convert the offset to an int
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		f.ErrorPrintf("Error converting offset to int: %s\n", err)
		http.Error(w, "Offset is not a number", http.StatusBadRequest)
		return
	}

	thread := f.GetThreadFromName(threadName)
	threadConfig := f.GetThreadConfigFromThread(thread)
	user := f.GetUser(r)

	realTags := []f.ThreadTag{}
	// Check if the tags are empty or not a valid tag
	if tags != "" {
		var tagsList []string
		err := json.Unmarshal([]byte(tags), &tagsList)
		if err != nil {
			f.ErrorPrintf("Error parsing tags: %s\n", err)
			http.Error(w, "Tags are not a valid JSON array", http.StatusBadRequest)
			return
		}
		realTags, err = stringTagsToThreadTags(tagsList, thread)
	}

	// Check if the user is banned from the thread
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

	Messages, err := f.GetMessagesFromThreadWithPOV(
		f.GetThreadFromName(threadName),
		offsetInt,
		order,
		user,
		realTags,
	)
	if err != nil {
		f.ErrorPrintf("Error getting messages from thread: %s\n", err)
		http.Error(w, "Error getting messages from thread", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Messages)
	if err != nil {
		f.ErrorPrintf("Error encoding messages to JSON: %s\n", err)
		http.Error(w, "Error encoding messages to JSON", http.StatusInternalServerError)
		return
	}
}

func stringTagsToThreadTags(tags []string, thread f.ThreadGoForum) ([]f.ThreadTag, error) {
	threadTags, err := f.GetThreadTags(thread)
	if err != nil {
		f.ErrorPrintf("Error getting thread tags: %s\n", err)
		return nil, err
	}
	var limitedTags []f.ThreadTag
	for _, tag := range tags {
		for _, threadTag := range threadTags {
			if tag == threadTag.TagName {
				limitedTags = append(limitedTags, threadTag)
			}
		}
	}
	return limitedTags, nil
}
