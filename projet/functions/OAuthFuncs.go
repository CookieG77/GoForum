package functions

import (
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"os"
)

var pathAPIKey = "OAuthKeys.json"

type OAuthConfig struct {
	ID     string `json:"ID"`
	Secret string `json:"SECRET"`
}

// SetupOAuth sets up the OAuth configuration.
func SetupOAuth(port string) {
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
