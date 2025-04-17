package pagesHandlers

import (
	m "GoForum/backend/emailsHandlers"
	f "GoForum/functions"
	"net/http"
)

func ConfirmMailPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("confirmMail", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Confirm Mail page accessed at %s by unverified %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
		} else {
			f.InfoPrintf("Confirm Mail page accessed at %s by verified %s : %s\n", f.GetIP(r), f.GetUserRankString(r), f.GetUserEmail(r))
		}
	} else {
		f.InfoPrintf("Confirm Mail page accessed at %s\n", f.GetIP(r))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Getting the user
	user := f.GetUser(r)

	// Setting the pages values to default
	PageInfo["Error"] = ""
	PageInfo["Success"] = false
	PageInfo["UserMail"] = user.Email
	PageInfo["UserUsername"] = user.Username

	// If we receive a POST request, we resend the email
	if r.Method == "POST" {
		// Check if the user is already confirmed
		if f.IsUserVerified(r) {
			PageInfo["Error"] = "alreadyVerified"
		} else {
			// Send the email
			PageInfo["Success"] = true
			m.SendConfirmEmail(user.Email)

		}
	} else if r.Method == "GET" {
		// If we receive a GET request, we check if the user is already confirmed
		if f.IsUserVerified(r) {
			PageInfo["Error"] = "alreadyVerified"
		} else {
			// If the user is not verified, we check if the token is valid
			token := r.URL.Query().Get("token")
			if token != "" {
				if !f.CheckEmailIdentification(token, f.VerifyEmailEmail) {
					// The token is invalid or expired
					PageInfo["Error"] = "invalidToken"
					f.DebugPrintln("Given token is invalid")
				} else {
					// The token is valid
					err := f.VerifyEmail(user.Email)
					if err != nil {
						f.ErrorPrintf("Error while verifying the email address: %s\n", err)
						ErrorPage(w, r, http.StatusInternalServerError)
						return
					}
					err = f.RemoveEmailIdentificationWithID(token)
					if err != nil {
						f.ErrorPrintf("Error while removing the email identification: %s\n", err)
						ErrorPage(w, r, http.StatusInternalServerError)
					}
					// If the email is verified, we redirect the user to the home page
					http.Redirect(w, r, "/", http.StatusFound)
				}
			}
		}
	}

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/confirmEmailAddress.css", "generalElementStyling.css")
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/confirmEmailAddress.html")
}
