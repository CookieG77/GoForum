package pagesHandlers

import (
	m "GoForum/backend/emailsHandlers"
	f "GoForum/functions"
	"net/http"
)

func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("home", w, r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		f.InfoPrintf("Home page accessed at %s by %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
		// Redirect the user to the home page if he is already authenticated
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		f.InfoPrintf("Home page accessed at %s\n", f.GetIP(r))
	}

	PageInfo["bareboneBase"] = true // This is a barebone page (no header or useless stuff)
	PageInfo["Error"] = ""          // No error message by default
	PageInfo["Provider"] = ""       // No provider message by default
	PageInfo["Success"] = false     // No success message by default
	// If the request is a POST request, we try to reset the password
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			f.ErrorPrintf("Error while parsing the form: %s\n", err)
			ErrorPage(w, r, 500)
			return
		}
		email := r.FormValue("email")
		if email == "" {
			PageInfo["Error"] = "noMailProvided"
		}
		if PageInfo["Error"] == "" {
			// Check if the email address is valid
			if f.IsEmailValid(email) {
				// We only send a mail if the email address is valid and not associated with an OAuth provider
				if b, provider := f.CheckIfEmailLinkedToOAuth(email); b {
					PageInfo["Error"] = "linkedToOAuth"
					PageInfo["Provider"] = provider
				} else {
					f.DebugPrintf("Sending a reset password mail to %s\n", email)
					m.SendResetPasswordMail(email)
				}
			}
			// We don't tell the user if the email address is invalid
			// We don't want to give any information about the users
			PageInfo["Success"] = true
		}
	}

	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/resetPassword.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/resetPassword.html")
}
