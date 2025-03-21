package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("home", w, r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Home page accessed at %s by unverified %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirmMail", http.StatusFound)
			return
		}
		f.InfoPrintf("Home page accessed at %s by verified %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("Home page accessed at %s\n", f.GetIP(r))
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/home.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/home.html")
}
