package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Lang is a type representing a language
type Lang string

// Constants representing the different languages
const (
	En Lang = "en"
	Fr Lang = "fr"
)

// langList is a list of all the languages
var langList = []Lang{En, Fr}

// DefaultLang is the default language
var DefaultLang Lang

// InitDefaultLangConfig set up the default language
func InitDefaultLangConfig() {
	defaultLang := os.Getenv("DEFAULT_LANG")
	if defaultLang == "" {
		DefaultLang = En
		WarningPrintf("No default language was given switching to : %s\n", DefaultLang)
		return
	}
	DefaultLang = StrToLang(defaultLang)
	SuccessPrintf("Default Language set to : %s\n", DefaultLang)
}

// GetLangContent returns the map[string]string containing each field text adapted to the given language
// In case some fields are missing we first load the default version and overwrite it with the given language
// Also returns an error if the file can't be read or if the json can't be unmarshalled
func GetLangContent(language Lang) (map[string]interface{}, error) {

	// First we load the default language
	defaultLangFilePath := fmt.Sprintf("statics/lang/%s.json", DefaultLang)

	file, err := os.ReadFile(defaultLangFilePath)
	if err != nil {
		return nil, err
	}

	var defaultContent map[string]interface{}
	err = json.Unmarshal(file, &defaultContent)
	if err != nil {
		return nil, err
	}
	if language == DefaultLang {
		return defaultContent, nil
	}

	// Now we load the language asked
	langFilePath := fmt.Sprintf("statics/lang/%s.json", language)

	file, err = os.ReadFile(langFilePath)
	if err != nil {
		return nil, err
	}

	var content map[string]interface{}
	err = json.Unmarshal(file, &content)
	if err != nil {
		return nil, err
	}

	// We overwrite the default content with the language content
	// Since both defaultContent and content are map[string]interface{} that may also contain map[string]interface{} that may contain map[string]interface{} and so on...
	// We need to use a recursive function to merge the two maps
	mergedContent := mergeMap(defaultContent, content)
	return mergedContent, nil
}

// StrToLang converts a string to a Lang
// Returns En if the string doesn't match any Lang
func StrToLang(str string) Lang {
	for _, l := range langList {
		if string(l) == str {
			return l
		}
	}
	return En
}

// LangListToStrList return the given list of Lang and return it as list of string
func LangListToStrList(langList []Lang) []string {
	strList := make([]string, len(langList))
	for i, l := range langList {
		strList[i] = string(l)
	}
	return strList
}

// GetUserLang return the language of the user.
func GetUserLang(r *http.Request) Lang {
	if !IsAuthenticated(r) {
		return DefaultLang
	}
	u := GetUser(r)
	c := GetUserConfig(u)
	return StrToLang(c.Lang)
}

// GetLangList return the list of all the languages
func GetLangList() []Lang {
	return langList
}

// mergeMap is a recursive function that merge two maps
// It will overwrite the values of the first map with the values of the second map
// If a value is a map, it will call itself with the two maps
func mergeMap(m1, m2 map[string]interface{}) map[string]interface{} {
	for k, v := range m2 {
		if _, ok := m1[k]; ok {
			if m1Map, ok := m1[k].(map[string]interface{}); ok {
				if m2Map, ok := v.(map[string]interface{}); ok {
					m1[k] = mergeMap(m1Map, m2Map)
				} else {
					m1[k] = v
				}
			} else {
				m1[k] = v
			}
		} else {
			m1[k] = v
		}
	}
	return m1
}
