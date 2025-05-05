package pagesHandlers

import (
	f "GoForum/functions"
	"github.com/gorilla/mux"
	"net/http"
)

// UserOtherProfilePage handles the user profile page for other users
func UserOtherProfilePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	// Check if the user exists
	if !f.CheckIfUsernameExists(user) { // If the user does not exist
		ErrorPage404(w, r)
		return
	}

	PageInfo := f.NewContentInterface("profile", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("User Profile page of '%s' accessed at %s by unverified : %s\n", user, f.GetIP(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirm-email-address", http.StatusFound)
			return
		}
		f.InfoPrintf("User Profile page of '%s' accessed at %s by verified : %s\n", user, f.GetIP(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("User Profile page of '%s' accessed at %s\n", user, f.GetIP(r))
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// TODO : Display the user profile of the 'user'

	// Add additional styles to the content interface
	f.AddAdditionalStylesToContentInterface(&PageInfo, "/css/userProfile.css", "/css/generalElementStyling.css")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/userProfile.html")
}
