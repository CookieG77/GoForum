package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("home", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Home page accessed at %s by unverified : %s\n", f.GetIP(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirmMail", http.StatusFound)
			return
		}
		f.InfoPrintf("Home page accessed at %s by verified : %s\n", f.GetIP(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("Home page accessed at %s\n", f.GetIP(r))
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// If we come from the register page with the value "openlogin" in the URL, we open the login popup
	if r.URL.Query().Get("openlogin") == "true" {
		PageInfo["ShowLoginPage"] = true
	}

	PageInfo["AllThreads"] = f.GetAllFormattedThreads()

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/home.css")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/home.html")
}

// RedirectToLogin redirects the user to the login page if they are not authenticated.
func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect to the login page if the user is not authenticated
	http.Redirect(w, r, "/?openlogin=true", http.StatusFound)
}
