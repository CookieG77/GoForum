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
			f.InfoPrintf("Confirm Mail page accessed at %s by unverified : %s\n", f.GetIP(r), f.GetUserEmail(r))
		} else {
			f.InfoPrintf("Confirm Mail page accessed at %s by verified : %s\n", f.GetIP(r), f.GetUserEmail(r))
		}
	} else {
		f.InfoPrintf("Confirm Mail page accessed at %s\n", f.GetIP(r))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	// Check if the user is regular or oauth
	connectMethod, _ := f.GetConnectionMethod(f.GetUserEmail(r))
	// Getting the user
	user := f.GetUser(r)
	if connectMethod == "oauth" {
		PageInfo["Error"] = false
		PageInfo["MissingField"] = make(map[string]bool)
		PageInfo["UsernameError"] = ""

		// Add additional styles to the content interface and make the template
		PageInfo["bareboneBase"] = true // This is to remove most of the base template leaving only the logo
		f.AddAdditionalStylesToContentInterface(&PageInfo, "css/loginAndRegister.css")

		if r.Method == "POST" {
			f.DebugPrintf("Completeing Registration for OAuth form submitted\n") // Parse the form
			err := r.ParseForm()
			if err != nil {
				f.ErrorPrintf("Error parsing the form: %v\n", err)
				PageInfo["Error"] = true
				f.MakeTemplateAndExecute(w, PageInfo, "templates/register.html")
				return
			}
			// Check if the use is trying to log out
			if r.Form.Get("logout") == "true" {
				f.DebugPrintf("User is trying to log out\n")
				err := f.EmptySessionCookie(w, r)
				if err != nil {
					f.ErrorPrintf("Error emptying the session cookie: %v\n", err)
				}
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			// Getting the form values
			firstName := r.Form.Get("first_name")
			lastName := r.Form.Get("last_name")
			username := r.Form.Get("username")
			acceptTerms := r.Form.Get("terms")
			// Check if a field is empty
			if firstName == "" || lastName == "" || username == "" || acceptTerms == "" {
				if firstName == "" {
					f.DebugPrintf("First name is empty\n")
					PageInfo["MissingField"].(map[string]bool)["firstName"] = true
				}
				if lastName == "" {
					f.DebugPrintf("Last name is empty\n")
					PageInfo["MissingField"].(map[string]bool)["lastName"] = true
				}
				if username == "" {
					f.DebugPrintf("Username is empty\n")
					PageInfo["MissingField"].(map[string]bool)["username"] = true
				}
				if acceptTerms == "" {
					f.DebugPrintf("Accept terms is empty\n")
					PageInfo["MissingField"].(map[string]bool)["acceptTerms"] = true
				}
				f.MakeTemplateAndExecute(w, PageInfo, "templates/completeRegistration.html")
				return
			}
			// Check if the username is valid
			if !f.IsUsernameValid(username) {
				f.DebugPrintf("Username is not valid")
				PageInfo["UsernameError"] = "invalid"
				f.MakeTemplateAndExecute(w, PageInfo, "templates/completeRegistration.html")
				return
			}

			// Check if the username is already in the database
			if f.CheckIfUsernameExists(username) {
				f.DebugPrintf("Username is already in use")
				PageInfo["UsernameError"] = "alreadyInUse"
				f.MakeTemplateAndExecute(w, PageInfo, "templates/completeRegistration.html")
				return
			}

			// Check if the accept terms is checked
			if acceptTerms != "on" {
				f.DebugPrintf("Terms are not accepted")
				PageInfo["MissingField"].(map[string]bool)["acceptTerms"] = true
				f.MakeTemplateAndExecute(w, PageInfo, "templates/register.html")
				return
			}
			f.InfoPrintf("Completing registration for OAuth user %s", user.Email)
			// Complete the registration
			err = f.CompleteOAuthRegistration(user, username, firstName, lastName)
			if err != nil {
				f.ErrorPrintf("Error while completing the registration: %s\n", err)
				ErrorPage(w, r, http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound) // After the user is logged in, redirect to the home page
		}
		f.MakeTemplateAndExecute(w, PageInfo, "templates/completeRegistration.html")
	} else {

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
		f.MakeTemplateAndExecute(w, PageInfo, "templates/confirmEmailAddress.html")
	}
}
