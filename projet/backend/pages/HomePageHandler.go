package pages

import (
	f "GoForum/functions"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("home", w, r)
	if f.IsAuthenticated(r) {
		PageInfo["IsAuthenticated"] = true
		f.InfoPrintf("Home page accessed at %s by %s\n", f.GetIP(r), f.GetUserEmail(r))
	} else {
		PageInfo["IsAuthenticated"] = false
		f.InfoPrintf("Home page accessed at %s\n", f.GetIP(r))
	}

	// Handle the user logout
	ConnectFromHeader(w, r, &PageInfo)
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/home.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/home.html")
}
