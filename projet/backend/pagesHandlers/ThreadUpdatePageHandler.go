package pagesHandlers

import (
	f "GoForum/functions"
	"github.com/gorilla/mux"
	"net/http"
)

func ThreadOptionPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadName := vars["thread"]

	PageInfo := f.NewContentInterface("thread_update", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Thread option page accessed at %s by unverified %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirmMail", http.StatusFound)
			return
		}
		if !f.IsThreadOwner(f.GetThreadFromName(threadName), r) {
			f.InfoPrintf("Thread option page accessed at %s by verified non owner %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
			ErrorPage404(w, r)
			return
		}
		f.InfoPrintf("Thread option page accessed at %s by verified owner %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		// If not authenticated, redirect to the login page
		f.InfoPrintf("Thread option page accessed at %s\n", f.GetIP(r))
		http.Redirect(w, r, "/?openlogin=true", http.StatusFound)
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Check if the thread name is empty or does not exist
	if threadName == "" || !f.CheckIfThreadNameExists(threadName) {
		ErrorPage404(w, r)
		return
	}

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/threadOption.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/threadOption.html")
}
