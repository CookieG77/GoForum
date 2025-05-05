package pagesHandlers

import (
	oauthCallbacks "GoForum/backend/oauthCallbackPageHandlers"
	f "GoForum/functions"
	"github.com/gorilla/mux"
	"net/http"
)

func CallbackRedirection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	f.DebugPrintf("CallbackRedirection called: provider: %s", provider)
	switch provider {
	case "google":
		f.DebugPrintf("GoogleCallbackHandler called")
		oauthCallbacks.GoogleCallback(w, r)
	case "github":
		f.DebugPrintf("GitHubCallbackHandler called")
		oauthCallbacks.GitHubCallback(w, r)
	case "discord":
		f.DebugPrintf("DiscordCallbackHandler called")
		oauthCallbacks.DiscordCallback(w, r)
	default:
		f.ErrorPrintf("CallbackRedirection: provider is empty")
		ErrorPage404(w, r)
	}
}
