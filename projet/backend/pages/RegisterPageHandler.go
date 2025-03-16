package pages

import (
	f "GoForum/functions"
	"net/http"
)

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("register", w, r)
	if f.IsAuthenticated(r) {
		PageInfo["IsAuthenticated"] = true
		f.InfoPrintf("Register page accessed at %s by %s\n", f.GetIP(r), f.GetUserEmail(r))
		// If the user is already authenticated, redirect him to the home page
		f.DebugPrintln("User is already authenticated")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		PageInfo["IsAuthenticated"] = false
		f.InfoPrintf("Register page accessed at %s\n", f.GetIP(r))
	}

	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/loginAndRegister.css")
	PageInfo["bareboneBase"] = true // This is to remove most of the base template leaving only the logo

	// Those are the variables that will be used in the template to display errors
	PageInfo["Error"] = false
	PageInfo["MissingField"] = make(map[string]bool)
	PageInfo["PasswordError"] = ""
	PageInfo["EmailError"] = ""
	PageInfo["UsernameError"] = ""

	// Check if the form is submitted
	if r.Method == "POST" {
		f.DebugPrintf("Registration form submitted\n")
		// Parse the form
		err := r.ParseForm()
		if err != nil {
			f.ErrorPrintf("Error parsing the form: %v\n", err)
			PageInfo["Error"] = true
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		// Getting the form values
		firstName := r.Form.Get("first_name")
		lastName := r.Form.Get("last_name")
		email := r.Form.Get("email")
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		confirmPassword := r.Form.Get("confirm_password")
		acceptTerms := r.Form.Get("terms")
		// Check if a field is empty
		if firstName == "" || lastName == "" || email == "" || username == "" || password == "" || confirmPassword == "" {
			// If a field is empty, add it to the PageInfo["MissingField"] map
			if firstName == "" {
				f.DebugPrintf("First name is empty\n")
				PageInfo["MissingField"].(map[string]bool)["firstName"] = true
			}
			if lastName == "" {
				f.DebugPrintf("Last name is empty\n")
				PageInfo["MissingField"].(map[string]bool)["lastName"] = true
			}
			if email == "" {
				f.DebugPrintf("Email is empty\n")
				PageInfo["MissingField"].(map[string]bool)["email"] = true
			}
			if username == "" {
				f.DebugPrintf("Username is empty\n")
				PageInfo["MissingField"].(map[string]bool)["username"] = true
			}
			if password == "" {
				f.DebugPrintf("Password is empty\n")
				PageInfo["MissingField"].(map[string]bool)["password"] = true
			}
			if confirmPassword == "" {
				f.DebugPrintf("Confirm password is empty\n")
				PageInfo["MissingField"].(map[string]bool)["confirmPassword"] = true
			}
			if acceptTerms == "" {
				f.DebugPrintf("Accept terms is empty\n")
				PageInfo["MissingField"].(map[string]bool)["acceptTerms"] = true
			}
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		// Reading the form values
		PageInfo["ValueFirstName"] = firstName
		PageInfo["ValueLastName"] = lastName
		PageInfo["ValueEmail"] = email
		PageInfo["ValueUsername"] = username
		// We don't want to keep the password in the form values
		// The user will have to retype it if he makes a mistake

		// Check if the mail is valid
		if !f.IsEmailValid(email) {
			f.DebugPrintf("Email is not valid")
			PageInfo["EmailError"] = "invalid"
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		// Check if the email is already linked to an OAuth account
		linked, provider := f.CheckIfEmailLinkedToOAuth(email)
		if linked {
			f.DebugPrintf("Email is already linked to an OAuth account : %s", provider)
			PageInfo["EmailError"] = "alreadyInUseOAuth"
			PageInfo["Provider"] = provider
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		// Check if the email is already in the database
		if f.CheckIfEmailExists(email) {
			f.DebugPrintf("Email is already in use")
			PageInfo["EmailError"] = "alreadyInUse"
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		// Check if the username is already in the database
		if f.CheckIfUsernameExists(username) {
			f.DebugPrintf("Username is already in use")
			PageInfo["UsernameError"] = "alreadyInUse"
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		// Check if the username is valid
		if !f.IsUsernameValid(username) {
			f.DebugPrintf("Username is not valid")
			PageInfo["UsernameError"] = "invalid"
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		// Check if the password is strong enough
		if !f.CheckPasswordStrength(password) {
			f.DebugPrintf("Password is not strong enough")
			PageInfo["PasswordError"] = "invalid"
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		// Check if the password and the confirmation password are the same
		if password != confirmPassword {
			f.DebugPrintf("Passwords do not match")
			PageInfo["PasswordError"] = "different"
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		// Insert the user in the database
		err = f.AddUser(email, username, firstName, lastName, password)
		if err != nil {
			f.ErrorPrintf("Error adding the user: %v\n", err)
			PageInfo["Error"] = true
			f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
			return
		}
		f.DebugPrintf(
			"New user added to the database:\n\t- email : %s\n\t- username : %s\n\t- firstName : %s\n\t- lastName : %s\n",
			email, username, firstName, lastName,
		)
		// Set the session cookie
		err = f.SetSessionCookie(w, r, email)
		if err != nil {
			f.ErrorPrintf("Error setting the session cookie: %v\n", err)
			// Since the user is already in the database, we can still redirect him to the login page if the session cookie is not set
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
	f.MakeTemplateAndExecute(w, r, PageInfo, "templates/register.html")
}
