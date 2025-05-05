package pagesHandlers

import (
	f "GoForum/functions"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func ThreadEditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadName := vars["threadName"]

	PageInfo := f.NewContentInterface("thread_config", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Thread config page accessed at %s by unverified : %s\n", f.GetIP(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirmMail", http.StatusFound)
			return
		}
		f.InfoPrintf("Thread config page accessed at %s by verified : %s\n", f.GetIP(r), f.GetUserEmail(r))
	} else {
		// If not authenticated, redirect to the login page
		f.InfoPrintf("Thread config page accessed at %s\n", f.GetIP(r))
		http.Redirect(w, r, "/?openlogin=true", http.StatusFound)
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Check if the thread name is empty or does not exist
	if threadName == "" || !f.CheckIfThreadNameExists(threadName) {
		ErrorPage404(w, r)
		return
	}

	thread := f.GetThreadFromName(threadName)
	user := f.GetUser(r)

	// If not the thread owner, redirect to the thread page
	if thread.OwnerID != user.UserID {
		f.DebugPrintf("User is not the owner of the thread he's trying to edit\n")
		http.Redirect(w, r, fmt.Sprintf("/t/%s", threadName), http.StatusFound)
		return
	}

	PageInfo["ErrorEditingThread"] = false

	// Handle the thread edit form
	if r.Method == "POST" {
		// parse the form
		err := r.ParseForm()
		if err != nil {
			f.ErrorPrintf("Error parsing the form : %s\n", err)
			PageInfo["ErrorEditingThread"] = true
			return
		}
	}

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "/css/threadEdit.css")
	f.AddAdditionalScriptsToContentInterface(&PageInfo, "/js/threadEditScript.js", "/js/threadScript.js")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/threadEdit.html")
}
