package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
)

// UserSelfProfilePage handles the user profile page for the authenticated user
func UserSelfProfilePage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("profile", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("User Self Profile page accessed at %s by unverified %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirm-email-address", http.StatusFound)
			return
		}
		f.InfoPrintf("User Self Profile page accessed at %s by verified %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("User Self Profile page accessed at %s\n", f.GetIP(r))
		// If the user is not authenticated, show him a forbidden page
		ErrorPage403(w, r)
		return
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// TODO : Display the user profile of the 'user'

	// Add additional styles to the content interface
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/userSelfProfile.css", "generalElementStyling.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/userSelfProfile.html")
}
