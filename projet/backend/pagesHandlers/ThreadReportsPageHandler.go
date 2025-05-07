package pagesHandlers

import (
	f "GoForum/functions"
	"github.com/gorilla/mux"
	"net/http"
)

func ThreadReportsPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadName := vars["threadName"]
	PageInfo := f.NewContentInterface("thread_reports", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Thread option page accessed at %s by unverified : %s\n", f.GetIP(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirmMail", http.StatusFound)
			return
		}
		if !(f.GetUserRankInThread(f.GetThreadFromName(threadName), f.GetUser(r)) >= 1) {
			f.InfoPrintf("Thread option page accessed at %s by verified non moderation team member : %s\n", f.GetIP(r), f.GetUserEmail(r))
			ErrorPage403(w, r) // Forbidden access
			return
		}
		f.InfoPrintf("Thread option page accessed at %s by verified moderation team member : %s\n", f.GetIP(r), f.GetUserEmail(r))
	} else {
		// If not authenticated, redirect to the login page
		f.InfoPrintf("Thread option page accessed at %s\n", f.GetIP(r))
		RedirectToLogin(w, r)
		return
	}

	// Check if the thread name is empty or does not exist
	if threadName == "" || !f.CheckIfThreadNameExists(threadName) {
		f.DebugPrintf("Thread name is empty or does not exist : %s\n", threadName)
		ErrorPage404(w, r)
		return
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	thread := f.GetThreadFromName(threadName)

	reports, err := f.GetReportedContentInThread(thread)
	if err != nil {
		f.ErrorPrintf("Error while getting the reports for thread %s : %s\n", threadName, err)
		ErrorPage404(w, r)
		return
	}
	PageInfo["Reports"] = reports
	PageInfo["ThreadName"] = threadName

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "/css/threadReports.css")
	f.AddAdditionalScriptsToContentInterface(&PageInfo, "/js/threadScript.js", "/js/threadReports.js")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/threadReports.html")
}
