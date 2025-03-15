package pages

import (
	f "GoForum/functions"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	f.InfoPrintf("Login page accessed by %s", f.GetIP(r))
	PageInfo := f.NewContentInterface("login", w, r)
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/loginAndRegister.css")
	PageInfo["bareboneBase"] = true // This is to remove most of the base template leaving only the logo
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/login.html")
}
