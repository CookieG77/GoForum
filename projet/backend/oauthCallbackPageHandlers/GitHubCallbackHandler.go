package oauthCallbackPageHandlers

import (
	f "GoForum/functions"
	"net/http"
)

func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	f.DebugPrintf("GitHubCallbackHandler called")
}
