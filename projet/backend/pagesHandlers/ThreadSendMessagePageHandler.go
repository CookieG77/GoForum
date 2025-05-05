package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
)

// TODO : Remove this page as we move the send message in the thread page
func ThreadSendMessagePage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("sendMessage", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Thread send message page accessed at %s by unverified : %s\n", f.GetIP(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirm-email-address", http.StatusFound)
			return
		}
		f.InfoPrintf("Thread send message page accessed at %s by verified : %s\n", f.GetIP(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("Thread send message page accessed at %s\n", f.GetIP(r))
		RedirectToLogin(w, r)
		return
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	PageInfo["UserThreads"] = f.GetUserThreads(f.GetUser(r))

	// Add additional styles to the content interface
	f.AddAdditionalStylesToContentInterface(&PageInfo, "/css/threadSendMessage.css", "/css/generalElementStyling.css")
	f.AddAdditionalScriptsToContentInterface(&PageInfo, "/js/threadScript.js", "/js/threadSendMessage.js", "/js/imgUploaderScript.js")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/threadSendMessage.html")
}
