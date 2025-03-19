package pagesHandlers

import (
	f "GoForum/functions"
	"net/http"
)

// ConnectFromHeader handles the forms in the header.
// If the form is a logout form, it logs the user out and returns false.
// If the form is a login form, it checks if the fields are empty, if an error occurs, it returns true.
// If the form is a login form and the fields are not empty, it connects the user and returns false.
// It is intended that the function is called at the beginning of the page handler so that the user is connected before the page is displayed.
func ConnectFromHeader(w http.ResponseWriter, r *http.Request, PageInfo *map[string]interface{}) bool {
	(*PageInfo)["LoginError"] = ""
	(*PageInfo)["LoginMissingField"] = map[string]bool{}
	(*PageInfo)["ShowLoginPage"] = false
	(*PageInfo)["Error"] = ""
	if r.Method == "POST" {
		// Parse the form
		err := r.ParseForm()
		if err != nil {
			f.ErrorPrintf("Error parsing the form: %v\n", err)
			(*PageInfo)["Error"] = "serverError"
			return true
		}
		switch r.Form.Get("headerForm") {
		case "logout":
			f.DebugPrintln("Logout form submitted")
			f.DebugPrintf("User %s logged out\n", f.GetUserEmail(r))
			err := f.EmptySessionCookie(w, r)
			if err != nil {
				f.ErrorPrintf("Error emptying the session cookie: %v\n", err)
			} else {
				(*PageInfo)["IsAuthenticated"] = false
			}
		case "login":
			f.DebugPrintln("Login form submitted")
			// Getting the form values
			emailOrUsername := r.Form.Get("email_or_username")
			password := r.Form.Get("password")
			// Check if a field is empty
			if emailOrUsername == "" || password == "" {
				// If a field is empty, add it to the PageInfo["MissingField"] map
				if emailOrUsername == "" {
					f.DebugPrintf("Email/Username is empty\n")
					(*PageInfo)["LoginMissingField"].(map[string]bool)["emailOrUsername"] = true
				}
				if password == "" {
					f.DebugPrintf("Password is empty\n")
					(*PageInfo)["LoginMissingField"].(map[string]bool)["password"] = true
				}
				(*PageInfo)["LoginError"] = "missingField"
				(*PageInfo)["ShowLoginPage"] = true
				return true
			}
			// Check if the user exists
			connectionMethod, provider := f.GetConnectionMethod(emailOrUsername)
			switch connectionMethod {
			case "": // If the credentials are invalid
				{
					f.DebugPrintf("User with mail/username '%s' does not exist\n", emailOrUsername)
					(*PageInfo)["LoginError"] = "invalidCredentials"
					(*PageInfo)["ShowLoginPage"] = true
					return true
				}
			case "oauth": // If the user entered an email that is registered with an OAuth
				{
					f.DebugPrintf("User with mail '%s' is registered with an OAuth from provider '%s'\n", emailOrUsername, provider)
					(*PageInfo)["LoginError"] = "userIsOAuth"
					(*PageInfo)["OAuthProvider"] = provider
					(*PageInfo)["ShowLoginPage"] = true
					return true
				}
			case "email": // If the user entered an email
				{
					f.DebugPrintf("Connection method: '%s'\n", connectionMethod)
					// Check if password is correct
					b, err := f.CheckUserConnectingWMail(emailOrUsername, password)
					if err != nil {
						f.ErrorPrintf("Error checking the user connecting: %v\n", err)
						(*PageInfo)["LoginError"] = "serverError"
						(*PageInfo)["ShowLoginPage"] = true
						return true
					}
					if !b {
						f.DebugPrintf("User with mail '%s' entered an incorrect password\n", emailOrUsername)
						(*PageInfo)["LoginError"] = "invalidCredentials"
						(*PageInfo)["ShowLoginPage"] = true
						return true
					}
					// Set the session cookie
					err = f.SetSessionCookie(w, r, emailOrUsername)
					if err != nil {
						f.ErrorPrintf("Error setting the session cookie: %v\n", err)
						(*PageInfo)["LoginError"] = "serverError"
						(*PageInfo)["ShowLoginPage"] = true
						return true
					}
					(*PageInfo)["IsAuthenticated"] = true
					f.InfoPrintf("User %s logged in\n", emailOrUsername)
					return false
				}
			case "username": // If the user entered a username
				{
					f.DebugPrintf("Connection method: '%s'\n", connectionMethod)
					// Check if password is correct
					b, err := f.CheckUserConnectingWUsername(emailOrUsername, password)
					if err != nil {
						f.ErrorPrintf("Error checking the user connecting: %v\n", err)
						(*PageInfo)["LoginError"] = "serverError"
						(*PageInfo)["ShowLoginPage"] = true
						return true
					}
					if !b {
						f.DebugPrintf("User with username '%s' entered an incorrect password\n", emailOrUsername)
						(*PageInfo)["LoginError"] = "invalidCredentials"
						(*PageInfo)["ShowLoginPage"] = true
						return true
					}
					email := f.GetEmailFromUsername(emailOrUsername)
					// Set the session cookie
					err = f.SetSessionCookie(w, r, email)
					if err != nil {
						f.ErrorPrintf("Error setting the session cookie: %v\n", err)
						(*PageInfo)["LoginError"] = "serverError"
						(*PageInfo)["ShowLoginPage"] = true
						return true
					}
					(*PageInfo)["IsAuthenticated"] = true
					f.InfoPrintf("User %s logged in\n", emailOrUsername)
					return false
				}
			default:
				{
					f.DebugPrintf("Unknown connection method: '%s\n'", connectionMethod)
				}
			}
		default:
			{
				f.DebugPrintf("Unknown connection method: '%s\n'", r.Form.Get("headerForm"))
			}
		}
	}
	return false
}
