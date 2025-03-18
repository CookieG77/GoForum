package pages

import (
	f "GoForum/functions"
	"net/http"
)

func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("home", w, r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		f.InfoPrintf("Home page accessed at %s by %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("Home page accessed at %s\n", f.GetIP(r))
	}

	PageInfo["bareboneBase"] = true // This is a barebone page (no header or useless stuff)

}
