package functions

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"net/http"
	"os"
)

// ConnectOAuth sets up the OAuth configuration.
func ConnectOAuth(port string) {
	// Adapt the scheme to the environment
	var scheme string
	if IsCertified() {
		scheme = "https"
	} else {
		scheme = "http"
	}

	// Add the OAuth providers here
	// TODO : add the other providers
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"),
			fmt.Sprintf("%s://localhost%s/auth/callback/google", scheme, port),
			"email",
		),
		github.New(
			os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"),
			fmt.Sprintf("%s://localhost%s/auth/callback/github", scheme, port),
			"user:email",
		),
		discord.New(
			os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"),
			fmt.Sprintf("%s://localhost%s/auth/callback/discord", scheme, port),
			"email", "identify",
		),
	)
}

// InitOAuthKeys initializes the OAuth keys and routes.
func InitOAuthKeys(finalPort string, r *mux.Router) {

	// Handle the OAuth routes
	ConnectOAuth(finalPort)
	r.HandleFunc("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})

	// link the store to the gothic package
	store, err := GetCookieStore()
	if err != nil {
		ErrorPrintf("Error getting the cookie store: %v\n", err)
		return
	}
	gothic.Store = store

	// Handle the OAuth callback routes
	r.HandleFunc("/auth/callback/{provider}", func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			ErrorPrintf("Error while completing the user auth: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// TODO : add the user to the database if we got enough information from the provider or redirect to a page to ask for more information.
		SuccessPrintf("User connected !\n\t%v\n", user) // TODO : remove this line in production.
		http.Redirect(w, r, "/", http.StatusSeeOther)   // TODO : redirect to the previous page instead of the home page.
	})
}
