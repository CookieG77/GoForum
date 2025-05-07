package pagesHandlers

import (
	f "GoForum/functions"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
	"net/http"
)

func CallbackRedirection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	f.DebugPrintf("CallbackRedirection called: provider: %s", provider)
	switch provider {
	case "google":
		f.DebugPrintf("GoogleCallbackHandler called")
		googleCallback(w, r)
	case "discord":
		f.DebugPrintf("DiscordCallbackHandler called")
		discordCallback(w, r)
	default:
		f.ErrorPrintf("CallbackRedirection: provider is empty")
		ErrorPage404(w, r)
	}
}

func googleCallback(w http.ResponseWriter, r *http.Request) {
	f.DebugPrintf("GoogleCallbackHandler called")
	oauthUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		f.ErrorPrintf("Error while completing the oauthUser auth: %v\n", err)
		ErrorPage500(w, r)
		return
	}
	f.DebugPrintf("GoogleCallbackHandler completed oauthUser: %+v", oauthUser)

	if f.UserWithProviderAndIDExist(f.GoogleOAuthProvider, oauthUser.UserID) { // Check if the user exists in the database
		user, err := f.GetUserFromOAuthProviderAndID(f.GoogleOAuthProvider, oauthUser.UserID)
		if err != nil {
			f.ErrorPrintf("Error getting the oauthUser from the database: %v\n", err)
			http.Error(w, "oauthUser not found in the database", http.StatusInternalServerError)
			return
		}
		if (user == f.User{}) {
			f.ErrorPrintf("Error: oauthUser not found in the database")
			http.Error(w, "oauthUser not found in the database", http.StatusInternalServerError)
			return
		}
		cookieMaxAge := 86400 // 1 day for all oauth users
		// Set the session cookie
		err = f.SetSessionCookie(w, r, user.Email, cookieMaxAge)
		if err != nil {
			f.ErrorPrintf("Error setting the session cookie: %v\n", err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther) // After the user is logged in, redirect to the home page
		return
	}
	if f.CheckIfEmailExists(oauthUser.Email) {
		f.DebugPrintf("GoogleCallbackHandler: email already exists in the database, redirecting login page")
		// If the email already exists in the database, redirect to the login page
		RedirectToLoginWithMessage(w, r, "An account with this email already exists. Please login.")
		return
	}

	// Create the user in the database, but we will use the verified status to check if the user finished the registration
	// Create temporary username
	temporaryUsername := fmt.Sprintf("user_%s_%s", oauthUser.UserID, uuid.New().String())
	err = f.AddUserWithOAuth(oauthUser.Email, temporaryUsername, f.GoogleOAuthProvider, oauthUser.UserID)
	if err != nil {
		f.ErrorPrintf("Error while creating the user in the database: %v\n", err)
		ErrorPage500(w, r)
		return
	}

	f.DebugPrintf("GoogleCallbackHandler: successfully created the user in the database")
	// log the user in
	err = f.SetSessionCookie(w, r, oauthUser.Email, 86400) // 1 day
	if err != nil {
		f.ErrorPrintf("Error setting the session cookie: %v\n", err)
		http.Error(w, "Error setting the session cookie", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/confirm-email-address", http.StatusSeeOther)
}

func discordCallback(w http.ResponseWriter, r *http.Request) {
	f.DebugPrintf("GitHubCallbackHandler called")
	oauthUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		f.ErrorPrintf("Error while completing the oauthUser auth: %v\n", err)
		ErrorPage500(w, r)
		return
	}
	f.DebugPrintf("DiscordCallbackHandler completed oauthUser: %+v", oauthUser)
	if f.UserWithProviderAndIDExist(f.DiscordOAuthProvider, oauthUser.UserID) { // Check if the user exists in the database
		user, err := f.GetUserFromOAuthProviderAndID(f.DiscordOAuthProvider, oauthUser.UserID)
		if err != nil {
			f.ErrorPrintf("Error getting the oauthUser from the database: %v\n", err)
			http.Error(w, "oauthUser not found in the database", http.StatusInternalServerError)
			return
		}
		if (user == f.User{}) {
			f.ErrorPrintf("Error: oauthUser not found in the database")
			http.Error(w, "oauthUser not found in the database", http.StatusInternalServerError)
			return
		}
		cookieMaxAge := 86400 // 1 day for all oauth users
		// Set the session cookie
		err = f.SetSessionCookie(w, r, user.Email, cookieMaxAge)
		if err != nil {
			f.ErrorPrintf("Error setting the session cookie: %v\n", err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther) // After the user is logged in, redirect to the home page
		return
	}
	if f.CheckIfEmailExists(oauthUser.Email) {
		f.DebugPrintf("DiscordCallbackHandler: email already exists in the database, redirecting login page")
		// If the email already exists in the database, redirect to the login page
		RedirectToLoginWithMessage(w, r, "An account with this email already exists. Please login.")
		return
	}
	// Create the user in the database, but we will use the verified status to check if the user finished the registration
	// Create temporary username
	temporaryUsername := fmt.Sprintf("user_%s_%s", oauthUser.UserID, uuid.New().String())
	err = f.AddUserWithOAuth(oauthUser.Email, temporaryUsername, f.DiscordOAuthProvider, oauthUser.UserID)
	if err != nil {
		f.ErrorPrintf("Error while creating the user in the database: %v\n", err)
		ErrorPage500(w, r)
		return
	}
	f.DebugPrintf("DiscordCallbackHandler: successfully created the user in the database")
	// log the user in
	err = f.SetSessionCookie(w, r, oauthUser.Email, 86400) // 1 day
	if err != nil {
		f.ErrorPrintf("Error setting the session cookie: %v\n", err)
		http.Error(w, "Error setting the session cookie", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/confirm-email-address", http.StatusSeeOther)
}
