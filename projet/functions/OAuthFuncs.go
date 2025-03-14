package functions

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"net/http"
	"os"
)

// ConnectOAuth sets up the OAuth configuration.
func ConnectOAuth(port string) {
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"),
			fmt.Sprintf("http://localhost%s/auth/callback/google", port),
			"email", "profile",
		),
		github.New(
			os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"),
			fmt.Sprintf("http://localhost%s/auth/callback/github", port),
			"user:email",
		),
		/*twitter.New(
			os.Getenv("TWITTER_CLIENT_ID"), os.Getenv("TWITTER_CLIENT_SECRET"),
			fmt.Sprintf("http://localhost%s/auth/callback/twitter", port),
			),
		twitch.New(
			os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"),
			fmt.Sprintf("http://localhost%s/auth/callback/twitch", port),
			),
		*/
	)
}

// SetupCookieStore sets up the cookie store.
func SetupCookieStore() *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
}

// InitOAuthKeys initializes the OAuth keys and routes.
func InitOAuthKeys(finalPort string, r *mux.Router, store *sessions.CookieStore) {

	// Handle the OAuth routes
	ConnectOAuth(finalPort)
	r.HandleFunc("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})

	// link the store to the gothic package
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
		SuccessPrintf("User connected !\n\t- Name : %v\n\t- Email : %v\n", user.Name, user.Email) // TODO : remove this line in production.
		http.Redirect(w, r, "/", http.StatusSeeOther)                                             // TODO : redirect to the previous page instead of the home page.
	})
}
