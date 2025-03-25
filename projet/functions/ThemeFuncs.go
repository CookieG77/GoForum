package functions

import (
	"net/http"
	"os"
)

// Theme is a type representing a theme
type Theme string

// Constants representing the different theme
const (
	Light Theme = "light"
	Dark  Theme = "dark"
)

// themeList is a list of all the theme
var themeList = []Theme{Light, Dark}

// DefaultTheme is the default theme
var DefaultTheme Theme

// InitDefaultThemeConfig set up the default theme
func InitDefaultThemeConfig() {
	defaultTheme := os.Getenv("DEFAULT_THEME")
	if defaultTheme == "" {
		DefaultTheme = Light
		WarningPrintf("No default theme was given switching to : %s\n", DefaultTheme)
		return
	}
	DefaultTheme = StrToTheme(defaultTheme)
	SuccessPrintf("Default theme set to : %s\n", DefaultTheme)
}

// StrToTheme convert a string to a Theme.
// If the string is not a valid theme, it returns the default theme (DefaultTheme).
func StrToTheme(s string) Theme {
	for _, theme := range themeList {
		if s == string(theme) {
			return theme
		}
	}
	return Light
}

// GetUserTheme return the theme of the user.
func GetUserTheme(r *http.Request) Theme {
	if !IsAuthenticated(r) {
		return DefaultTheme
	}
	c := GetUserConfig(r)
	return StrToTheme(c.Theme)
}
