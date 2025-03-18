package pages

import (
	f "GoForum/functions"
	"net/http"
)

func ErrorPage(w http.ResponseWriter, r *http.Request, status int) {
	PageInfo := f.NewContentInterface("home", w, r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		f.InfoPrintf("Error page accessed at %s by %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("Error page accessed at %s\n", f.GetIP(r))
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Set the error status
	PageInfo["ErrorStatus"] = status
	// Set the error message
	switch status {
	case http.StatusNotFound:
		PageInfo["ErrorMessage"] = "PageNotFound"
	case http.StatusForbidden:
		PageInfo["ErrorMessage"] = "Forbidden"
	case http.StatusInternalServerError:
		PageInfo["ErrorMessage"] = "InternalServerError"
	default:
		PageInfo["ErrorMessage"] = "Error"
	}

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/error.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/error.html")
}
