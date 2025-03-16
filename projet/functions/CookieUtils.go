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
func GetCookie(w http.ResponseWriter, r *http.Request, name string) *http.Cookie {
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

// ClearSessionCookie empties the session cookie for the user.
func ClearSessionCookie(w http.ResponseWriter, r *http.Request) error {
	if !isInitialised {
		// If the cookie store is not initialised, return an error
		return errors.New("session store not initialised")
	}
	session, err := GetSession(r)
	if err != nil {
		ErrorPrintf("Error getting the session: %v\n", err)
		return err
	}
	delete(session.Values, "session")
	err = session.Save(r, w)
	if err != nil {
		ErrorPrintf("Error saving the session: %v\n", err)
		return err
	}
	return nil
}
