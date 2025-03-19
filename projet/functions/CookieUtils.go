package functions

import (
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

var cookieStore *sessions.CookieStore
var isInitialised = false

// GetCookie returns the cookie with the given name.
func GetCookie(r *http.Request, name string) *http.Cookie {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil
	}
	return cookie
}

// SetCookie set the cookie with the given name to the given value.
// This cookie is not meant to be used for marketing or data analysing of the user.
// This implementation only serve to store a value in the user browser.
func SetCookie(w http.ResponseWriter, name string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// SetupCookieStore sets up the cookie store.
func SetupCookieStore() {
	cookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	isInitialised = true
}

// GetCookieStore returns the cookie store.
func GetCookieStore() (*sessions.CookieStore, error) {
	if !isInitialised {
		// If the cookie store is not initialised, return an error
		return nil, errors.New("session store not initialised")
	}
	return cookieStore, nil
}

// GetSession returns the session.
func GetSession(r *http.Request) (*sessions.Session, error) {
	if !isInitialised {
		// If the cookie store is not initialised, return an error
		return nil, errors.New("session store not initialised")
	}
	return cookieStore.Get(r, "session")
}

// SetSessionCookie sets the session cookie for the user.
// The session cookie is used to store the id of the user.
// The cookie is stored in the session store.
// Returns an error if there is one.
func SetSessionCookie(w http.ResponseWriter, r *http.Request, email string) error {
	session, err := GetSession(r)
	if err != nil {
		ErrorPrintf("Error getting the session: %v\n", err)
		return err
	}
	session.Values["email"] = email
	err = session.Save(r, w)
	if err != nil {
		ErrorPrintf("Error saving the session: %v\n", err)
		return err
	}
	DebugPrintf("Setting the session cookie for email: %v\n", email)
	return nil
}

// GetSessionCookie returns the session cookie for the user.
// Returns the session cookie and an error if there is one.
func GetSessionCookie(r *http.Request) (string, error) {
	session, err := GetSession(r)
	if err != nil {
		ErrorPrintf("Error getting the session: %v\n", err)
		return "", err
	}
	email := session.Values["email"].(string)
	return email, nil
}

// EmptySessionCookie empties the session cookie for the user.
// Returns an error if there is one.
func EmptySessionCookie(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSession(r)
	if err != nil {
		ErrorPrintf("Error getting the session: %v\n", err)
		return err
	}
	session.Values["email"] = ""
	err = session.Save(r, w)
	if err != nil {
		ErrorPrintf("Error saving the session: %v\n", err)
		return err
	}
	return nil
}
