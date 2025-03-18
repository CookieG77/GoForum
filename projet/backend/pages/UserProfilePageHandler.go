package pages

import (
	f "GoForum/functions"
	"net/http"
)

func UserProfilePage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("home", w, r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		f.InfoPrintf("User Profile page accessed at %s by %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("User Profile page accessed at %s\n", f.GetIP(r))
	}
	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	user := r.URL.Query().Get("user")
	if user != "" {
		// Check if the user exists
		if !f.CheckIfUsernameExists(user) { // If the user does not exist
			ErrorPage(w, r, http.StatusNotFound)
			return
		}
		f.DebugPrintf("Accessing the user profile of %s\n", user)
		// TODO : Display the user profile of the 'user'
	}
	if PageInfo["IsAuthenticated"].(bool) {
		f.DebugPrintf("Accessing the user profile of %s\n", f.GetUserRankString(r))
		// TODO : Display the user settings page of the user logged in
	}
	if user == "" {
		ErrorPage(w, r, http.StatusBadRequest)
		return
	}

	// Add additional styles to the content interface
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/userProfile.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/userProfile.html")
}
