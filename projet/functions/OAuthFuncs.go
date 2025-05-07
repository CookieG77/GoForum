package functions

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
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
}
