package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
)

func UserSettingsPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("home", w, r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("User Settings page accessed at %s by unverified %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
			// Redirect the user to the confirm email address page
			http.Redirect(w, r, "/confirm-email-address", http.StatusFound)
			return
		}
		f.InfoPrintf("User Settings page accessed at %s by verified %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("User Settings page accessed at %s\n", f.GetIP(r))
		// If the user is not authenticated, show him a forbidden page
		http.Redirect(w, r, "/?openlogin=true", http.StatusSeeOther)
		return
	}

	userConfig := f.GetUserConfig(r)

	PageInfo["LangList"] = f.LangListToStrList(f.GetLangList())
	PageInfo["UserLang"] = userConfig.Lang
	PageInfo["UserTheme"] = userConfig.Theme

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Add additional styles to the content interface
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/userSettings.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/userSettings.html")
}
