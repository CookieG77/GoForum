package pagesHandlers

import (
	f "GoForum/functions"
	"github.com/gorilla/mux"
	"net/http"
)

func ThreadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadName := vars["thread"]

	PageInfo := f.NewContentInterface("thread", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Thread \"%s\" page accessed at %s by unverified %s : %s\n", threadName, f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirm-email-address", http.StatusFound)
			return
		}
		f.InfoPrintf("Thread \"%s\" page accessed at %s by verified %s : %s\n", threadName, f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("Thread \"%s\" page accessed at %s\n", threadName, f.GetIP(r))
	}
	PageInfo["ThreadName"] = threadName

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Check if the thread name is empty or does not exist
	if threadName == "" || !f.CheckIfThreadNameExists(threadName) {
		ErrorPage404(w, r)
		return
	}

	thread := f.GetThreadFromName(threadName)
	threadConfig := f.GetThreadConfigFromThread(thread)
	threadIcon := f.GetMediaLinkFromID(threadConfig.ThreadIconID).MediaAddress
	if threadIcon == "" {
		f.DebugPrintf("Thread \"%s\" has no icon, using default icon\n", threadName)
		threadIcon = "default"
	}
	threadBanner := f.GetMediaLinkFromID(threadConfig.ThreadBannerID).MediaAddress
	if threadBanner == "" {
		f.DebugPrintf("Thread \"%s\" has no banner, using default banner\n", threadName)
		threadBanner = "default"
	}

	// Set the page variables
	PageInfo["MustJoinThread"] = false
	PageInfo["ThreadInfos"] = thread
	PageInfo["ThreadComplementaryInfos"] = threadConfig
	PageInfo["ShowContent"] = false
	PageInfo["ThreadIcon"] = threadIcon
	PageInfo["ThreadBanner"] = threadBanner

	//
	if (r.Method == "POST") && (PageInfo["IsAuthenticated"].(bool)) {
		err := r.ParseForm()
		if err != nil {
			f.ErrorPrintf("Error while parsing the form: %v\n", err)
			ErrorPage(w, r, http.StatusInternalServerError)
			return
		}
		if r.FormValue("action") == "join" {
			err := f.JoinThread(thread, r)
			if err != nil {
				f.ErrorPrintf("Error while joining thread %s : %v\n", thread.ThreadName, err)
				ErrorPage500(w, r)
				return
			}
		} else if r.FormValue("action") == "leave" {
			err := f.LeaveThread(thread, r)
			if err != nil {
				f.ErrorPrintf("Error while leaving thread %s : %v\n", thread.ThreadName, err)
				ErrorPage500(w, r)
				return
			}
		}
	}

	// If the user is not verified and this thread does not accept non-connected users, do not display the thread and open the login popup
	if !threadConfig.IsOpenToNonConnectedUsers && !PageInfo["IsAuthenticated"].(bool) {
		PageInfo["ShowLoginPage"] = true
		// If the user is not a member of the thread and this thread does not accept non-members, do not display the thread messages and show the 'must join' message
	} else if !threadConfig.IsOpenToNonMembers && !f.IsThreadMember(thread, r) {
		PageInfo["MustJoinThread"] = true
		// If the user is a member of the thread, display the thread normally
	} else {
		PageInfo["ShowContent"] = true
		PageInfo["IsAMember"] = f.IsThreadMember(thread, r)
	}

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "/css/thread.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/thread.html")
}
