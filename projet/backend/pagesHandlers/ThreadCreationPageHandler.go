package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
)

func ThreadCreationPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("thread_creation", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Thread creation page accessed at %s by unverified : %s\n", f.GetIP(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirmMail", http.StatusFound)
			return
		}
		f.InfoPrintf("Thread creation page accessed at %s by verified : %s\n", f.GetIP(r), f.GetUserEmail(r))
	} else {
		// If not authenticated, redirect to the login page
		f.InfoPrintf("Thread creation page accessed at %s\n", f.GetIP(r))
		http.Redirect(w, r, "/?openlogin=true", http.StatusFound)
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/threadCreation.css")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/threadCreation.html")
}
