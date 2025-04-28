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
	viewfrom := query.Get("viewfrom")

	isViewFromUserPOV := true
	User := f.User{}

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

	// Check if the viewfrom is empty or not a valid viewfrom
	if viewfrom == "" {
		isViewFromUserPOV = false
	} else {
		User, err = f.GetUserFromUsername(viewfrom)
		if err != nil {
			f.ErrorPrintf("Error getting user from username: %s\n", err)
			http.Error(w, "Error getting user from username", http.StatusBadRequest)
			return
		}
	}
	var Messages []f.FormattedThreadMessage
	if isViewFromUserPOV {
		Messages, err = f.GetMessagesFromThreadWithPOV(
			f.GetThreadFromName(threadName),
			offsetInt,
			order,
			User)
	} else {
		Messages, err = f.GetMessagesFromThread(
			f.GetThreadFromName(threadName),
			offsetInt,
			order)
		if err != nil {
			f.ErrorPrintf("Error getting messages from thread: %s\n", err)
			http.Error(w, "Error getting messages from thread", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Messages)
	if err != nil {
		f.ErrorPrintf("Error encoding messages to JSON: %s\n", err)
		http.Error(w, "Error encoding messages to JSON", http.StatusInternalServerError)
		return
	}
}
