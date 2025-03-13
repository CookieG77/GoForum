package pages

import (
	f "GoForum/functions"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	f.InfoPrintf("Home page accessed by %s", f.GetIP(r))
	PageInfo := f.NewContentInterface("home", w, r)
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/home.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/home.html")
}
