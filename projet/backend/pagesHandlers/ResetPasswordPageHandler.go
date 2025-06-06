package pagesHandlers

import (
	m "GoForum/backend/emailsHandlers"
	f "GoForum/functions"
	"net/http"
)

func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("resetPassword", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		f.InfoPrintf("Reset Password page accessed at %s by : %s\n", f.GetIP(r), f.GetUserEmail(r))
		// Redirect the user to the home page if he is already authenticated
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		f.InfoPrintf("Reset Password page accessed at %s\n", f.GetIP(r))
	}

	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/resetPassword.css", "generalElementStyling.css")

	PageInfo["bareboneBase"] = true    // This is a barebone page (no header or useless stuff)
	PageInfo["Error"] = ""             // No error message by default
	PageInfo["Provider"] = ""          // No provider message by default
	PageInfo["Success"] = false        // No success message by default
	PageInfo["ComingFromMail"] = false // Not coming from a mail by default
	PageInfo["MailToken"] = ""         // No mail id by default
	// If the request is a POST request, we try to reset the password
	if r.Method == "POST" {
		f.DebugPrintf("Accessing the reset password page with a POST request\n")
		err := r.ParseForm()
		if err != nil {
			f.ErrorPrintf("Error while parsing the form: %s\n", err)
			ErrorPage(w, r, 500)
			return
		}
		formType := r.FormValue("formType")
		if formType == "submitMail" {

			f.DebugPrintf("The formType is submitMail\n")
			// If the formType is submitMail, we send a reset password mail
			email := r.FormValue("email")
			if email == "" {
				PageInfo["Error"] = "noMailProvided"
			}
			if PageInfo["Error"] == "" {

				// Check if the email address exists in the database
				if f.CheckIfEmailExists(email) {
					// We only send a mail if the email address exists in the database and is not associated with an OAuth provider
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
		} else if formType == "submitPassword" {

			f.DebugPrintf("The formType is submitPassword\n")
			// If the formType is submitPassword, we try to change the password
			token := r.FormValue("token")
			password := r.FormValue("password")
			passwordConfirm := r.FormValue("passwordConfirm")

			PageInfo["ComingFromMail"] = true
			PageInfo["MailToken"] = token

			if token == "" || password == "" || passwordConfirm == "" { // Check if the fields are empty
				PageInfo["Error"] = "emptyFields"
			} else if password != passwordConfirm { // Check if the passwords match
				PageInfo["Error"] = "passwordsMismatch"
			} else if !f.CheckEmailIdentification(token, f.ResetPasswordEmail) {
				PageInfo["Error"] = "invalidToken"
			} else if !f.CheckPasswordStrength(password) { // Check if the password is valid
				PageInfo["Error"] = "passwordIncorrect"
			} else {
				userMail := f.GetEmailFromEmailIdentification(token)
				if userMail == "" {
					ErrorPage(w, r, 500)
					return
				}
				err := f.ChangeUserPassword(userMail, password)
				if err != nil {
					PageInfo["Error"] = "errorChangingPassword"
				} else {
					PageInfo["Success"] = true
					err := f.RemoveEmailIdentificationWithID(token)
					if err != nil {
						f.ErrorPrintf("Error while removing the email identification: %s\n", err)
					} else {
						f.InfoPrintf("User %s changed his password\n", userMail)
					}
				}
			}
		} else { // If the formType is not valid (Meaning the user tried to change the formType in the HTML code)
			ErrorPage(w, r, 400)
			return
		}
	} else if r.Method == "GET" {
		// If it's a GET request, we should have a token parameter in the URL
		// If it's not a GET request, we should not have any parameter in the URL
		f.DebugPrintf("Accessing the reset password page with a GET request\n")
		token := r.URL.Query().Get("token")
		if token != "" {
			f.DebugPrintf("A token was found in the URL: %s\n", token)
			if !f.CheckEmailIdentification(token, f.ResetPasswordEmail) {
				// The token from the URL is not valid
				PageInfo["Error"] = "invalidToken"
				f.DebugPrintln("Given token is invalid")
			} else {
				PageInfo["ComingFromMail"] = true
				PageInfo["MailToken"] = token
			}
		} else {
			f.DebugPrintln("No token was found in the URL")
		}
	}

	f.MakeTemplateAndExecute(w, PageInfo, "templates/resetPassword.html")
}
