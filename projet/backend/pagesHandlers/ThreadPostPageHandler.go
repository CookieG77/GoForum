package pagesHandlers

import (
	f "GoForum/functions"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ThreadPostPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadName := vars["threadName"]
	postID := vars["post"]

	PageInfo := f.NewContentInterface("threadPost", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Thread post \"%s\" page accessed at %s by unverified : %s\n", postID, f.GetIP(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirm-email-address", http.StatusFound)
			return
		}
		f.InfoPrintf("Thread post \"%s\" page accessed at %s by verified : %s\n", postID, f.GetIP(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("Thread post \"%s\" page accessed at %s\n", postID, f.GetIP(r))
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Check if the thread name is empty or does not exist
	if threadName == "" || !f.CheckIfThreadNameExists(threadName) {
		ErrorPage404(w, r)
		return
	}

	thread := f.GetThreadFromName(threadName)

	// Get User and user rights
	PageInfo["Username"] = ""
	PageInfo["UserRank"] = 0
	user := f.GetUser(r)
	if (user != f.User{}) {
		userRank := f.GetUserRankInThread(thread, user)
		if userRank < 0 { // If the user is banned from the thread we show him the YOU ARE BANNED page
			f.MakeTemplateAndExecute(w, PageInfo, "templates/youAreBanned.html")
			return
		}
		PageInfo["Username"] = user.Username
		PageInfo["UserRank"] = userRank
	}

	// Check if the post MessageID is empty or does not exist
	if postID == "" {
		ErrorPage404(w, r)
		return
	}
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		f.ErrorPrintf("Error converting post MessageID \"%s\" to int: %s\n", postID, err)
		ErrorPage404(w, r)
		return
	}

	// Check if the post exists
	if !f.MessageExistsInThread(thread, postIDInt) {
		f.ErrorPrintf("Post \"%s\" does not exist in thread \"%s\"\n", postID, threadName)
		ErrorPage404(w, r)
		return
	}

	var post f.FormattedThreadMessage
	if PageInfo["IsAddressVerified"].(bool) {
		post, err = f.GetMessageByIDWithPOV(postIDInt, user)
		PageInfo["IsAMember"] = f.IsUserInThread(thread, user)
	} else {
		post, err = f.GetMessageByID(postIDInt)
		PageInfo["IsAMember"] = false
	}
	if err != nil {
		f.ErrorPrintf("Error getting post \"%s\" : %s\n", postID, err)
		ErrorPage404(w, r)
		return
	}
	PageInfo["Post"] = post
	PageInfo["ThreadName"] = threadName
	PageInfo["ReportReasons"] = f.GetReportTypesAsStrings()

	f.AddAdditionalStylesToContentInterface(&PageInfo, "/css/threadPost.css", "/css/postStyle.css", "/css/thread.css")
	f.AddAdditionalScriptsToContentInterface(&PageInfo, "/js/threadPostScript.js", "/js/threadScript.js")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/threadPost.html")
}
