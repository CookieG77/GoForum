package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
)

func UserSettingsPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("home", w, r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		f.InfoPrintf("User Settings page accessed at %s by %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("User Settings page accessed at %s\n", f.GetIP(r))
	}
	// Redirect to error 403 if the user is not authenticated
	if !PageInfo["IsAuthenticated"].(bool) {
		ErrorPage(w, r, http.StatusForbidden)
		return
	}

	// TODO : Display the user settings page of the user logged in

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Add additional styles to the content interface
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/userSettings.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/userSettings.html")
}
