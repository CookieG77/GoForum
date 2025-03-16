package pages

import (
	f "GoForum/functions"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("register", w, r)
	if f.IsAuthenticated(r) {
		PageInfo["IsAuthenticated"] = true
		f.InfoPrintf("Login page accessed at %s by %s\n", f.GetIP(r), f.GetUserEmail(r))
		// If the user is already authenticated, redirect him to the home page
		f.DebugPrintln("User is already authenticated")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		PageInfo["IsAuthenticated"] = false
		f.InfoPrintf("Login page accessed at %s\n", f.GetIP(r))
	}

	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/loginAndRegister.css")
	PageInfo["bareboneBase"] = true // This is to remove most of the base template leaving only the logo
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/login.html")
}
