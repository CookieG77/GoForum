package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
	"slices"
)

func UserSettingsPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("user_settings", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("User Settings page accessed at %s by unverified : %s\n", f.GetIP(r), f.GetUserEmail(r))
			// Redirect the user to the confirm email address page
			http.Redirect(w, r, "/confirm-email-address", http.StatusFound)
			return
		}
		f.InfoPrintf("User Settings page accessed at %s by verified : %s\n", f.GetIP(r), f.GetUserEmail(r))
	} else {
		f.InfoPrintf("User Settings page accessed at %s\n", f.GetIP(r))
		// If the user is not authenticated, redirect to the login page
		RedirectToLogin(w, r)
		return
	}
	user := f.GetUser(r)
	userConfig := f.GetUserConfig(user)

	// Check if the user is changing his settings
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			f.ErrorPrintf("Error while parsing the user settings form : %s\n", err)
			ErrorPage(w, r, http.StatusInternalServerError)
			return
		}
		lang := r.Form.Get("lang")
		theme := r.Form.Get("theme")
		if lang == "" || theme == "" {
			f.ErrorPrintf("User settings form is missing fields\n")
			ErrorPage(w, r, http.StatusBadRequest)
			return
		}
		if !slices.Contains(f.LangListToStrList(f.GetLangList()), lang) {
			f.ErrorPrintf("User settings form has an invalid lang field\n")
			ErrorPage(w, r, http.StatusBadRequest)
			return
		}
		if !slices.Contains(f.ThemeListToStrList(f.GetThemeList()), theme) {
			f.ErrorPrintf("User settings form has an invalid theme field\n")
			ErrorPage(w, r, http.StatusBadRequest)
			return
		}
		userConfig.Lang = lang
		userConfig.Theme = theme
		err = f.UpdateUserConfig(userConfig)
		if err != nil {
			f.ErrorPrintf("Error while saving the user settings : %s\n", err)
			ErrorPage(w, r, http.StatusInternalServerError)
			return
		}

		// We need to update the PageInfo with the new userConfig values
		PageInfo = f.NewContentInterface("home", r)
		// Check the user rights
		f.GiveUserHisRights(&PageInfo, r)
	}

	PageInfo["LangList"] = f.LangListToStrList(f.GetLangList())
	PageInfo["UserLang"] = userConfig.Lang
	PageInfo["UserTheme"] = userConfig.Theme

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Add additional styles to the content interface
	f.AddAdditionalStylesToContentInterface(&PageInfo, "/css/userSettings.css")
	f.AddAdditionalScriptsToContentInterface(&PageInfo, "/js/userSettingsScript.js", "/js/imgUploaderScript.js")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/userSettings.html")
}
