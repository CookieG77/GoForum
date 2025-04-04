package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
)

func ErrorPage(w http.ResponseWriter, r *http.Request, status int) {
	PageInfo := f.NewContentInterface("error", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Error %d page accessed at %s by unverified %s : %s\n", status, f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirm-email-address", http.StatusFound)
			return
		}
		f.InfoPrintf("Error %d page accessed at %s by verified %s : %s\n", status, f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("Error %d page accessed at %s\n", status, f.GetIP(r))
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

func ErrorPage403(w http.ResponseWriter, r *http.Request) {
	ErrorPage(w, r, http.StatusForbidden)
}

func ErrorPage404(w http.ResponseWriter, r *http.Request) {
	ErrorPage(w, r, http.StatusNotFound)
}

func ErrorPage405(w http.ResponseWriter, r *http.Request) {
	ErrorPage(w, r, http.StatusMethodNotAllowed)
}

func ErrorPage500(w http.ResponseWriter, r *http.Request) {
	ErrorPage(w, r, http.StatusInternalServerError)
}
