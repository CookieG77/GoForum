package oauthCallbackPageHandlers

import (
	f "GoForum/functions"
	"github.com/markbates/goth/gothic"
	"net/http"
)

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	f.DebugPrintf("GooglebackHandler called")
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		f.ErrorPrintf("Error while completing the user auth: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	f.DebugPrintf("GooglebackHandler completed user: %+v", user)
}
