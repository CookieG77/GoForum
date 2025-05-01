package pagesHandlers

import (
	f "GoForum/functions"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

func ThreadPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadName := vars["thread"]

	PageInfo := f.NewContentInterface("thread", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Thread \"%s\" page accessed at %s by unverified : %s\n", threadName, f.GetIP(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirm-email-address", http.StatusFound)
			return
		}
		f.InfoPrintf("Thread \"%s\" page accessed at %s by verified : %s\n", threadName, f.GetIP(r), f.GetUserEmail(r))
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

	// Get User and user rights
	PageInfo["UserRank"] = 0
	user := f.GetUser(r)
	if (user != f.User{}) {
		userRank := f.GetUserRankInThread(thread, user)
		if userRank < 0 { // If the user is banned from the thread we show him the YOU ARE BANNED page
			f.MakeTemplateAndExecute(w, PageInfo, "templates/youAreBanned.html")
			return
		}
		PageInfo["UserRank"] = userRank
	}

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

	// If the user is not verified and this thread does not accept non-connected users, do not display the thread and open the login popup
	if !threadConfig.IsOpenToNonConnectedUsers && !PageInfo["IsAuthenticated"].(bool) {
		PageInfo["ShowLoginPage"] = true
		// If the user is not a member of the thread and this thread does not accept non-members, do not display the thread messages and show the 'must join' message
	} else if !threadConfig.IsOpenToNonMembers && !f.IsUserInThread(thread, user) {
		PageInfo["MustJoinThread"] = true
		// If the user is a member of the thread, display the thread normally
	} else {
		PageInfo["ShowContent"] = true
		PageInfo["IsAMember"] = f.IsThreadMember(thread, r)
		maxMessagesPerPageLoad := 10
		if os.Getenv("MAX_MESSAGES_PER_PAGE_LOAD") != "" {
			var err error
			maxMessagesPerPageLoad, err = strconv.Atoi(os.Getenv("MAX_MESSAGES_PER_PAGE_LOAD"))
			if err != nil {
				maxMessagesPerPageLoad = 10
			}
		}
		PageInfo["MaxMessagesPerPageLoad"] = maxMessagesPerPageLoad
	}
	// Add the thread moderation team to the page
	PageInfo["ThreadModerationTeam"] = f.GetThreadModerationTeam(thread)

	PageInfo["MessageOrdering"] = f.OrderingList

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "/css/thread.css", "/css/postStyle.css")
	f.AddAdditionalScriptsToContentInterface(&PageInfo, "/js/threadScript.js", "/js/threadPageScript.js")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/thread.html")
}
