package functions

import (
	"bytes"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
)

var isCertified bool = false
var isInitialized bool = false

// baseTemplates is a list of base templates to be used by the templates.
var baseTemplates []string

// AddBaseTemplate add a base template to the list of base templates.
func AddBaseTemplate(templatePath ...string) {
	baseTemplates = append(baseTemplates, templatePath...)
}

// MakeTemplate create a template from one or more template files given as parameter in the form of their path in string.
// It will use the base templates defined in the baseTemplates variable.
func MakeTemplate(w http.ResponseWriter, templatesDir ...string) *template.Template {
	templatesDir = append(templatesDir, baseTemplates...)
	tmpl, err := template.New("base.html").Funcs(template.FuncMap{
		"interfaceToString": interfaceToString,
	}).ParseFiles(templatesDir...)
	if err != nil {
		ErrorPrintf("An error occurred while trying to parse the template -> %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return tmpl
}

// ExecuteTemplate execute a template given as parameter.
func ExecuteTemplate(w http.ResponseWriter, tmpl *template.Template, content interface{}) {
	if tmpl == nil {
		ErrorPrintln("An error occurred while trying to execute a template -> Template is nil")
		http.Error(w, "Template is nil", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, content); err != nil {
		ErrorPrintf("An error occurred while trying to execute a template -> %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// MakeTemplateAndExecute is to be used for pagesHandlers that do not require any particular processing
func MakeTemplateAndExecute(w http.ResponseWriter, r *http.Request, content interface{}, templatesDir ...string) {
	tmpl := MakeTemplate(w, templatesDir...)
	ExecuteTemplate(w, tmpl, content)
}

// TemplateToText execute a template and return the result as a string.
func TemplateToText(tmpl *template.Template, content interface{}) string {
	// Check if the template is nil
	if tmpl == nil {
		ErrorPrintln("An error occurred while trying to execute a template -> Template is nil")
		return ""
	}
	// Create a buffer to store the content
	var contentBuffer bytes.Buffer

	// Execute the template
	if err := tmpl.Execute(&contentBuffer, content); err != nil {
		ErrorPrintf("An error occurred while trying to execute a template in TemplateToText -> %v\n", err)
		return ""
	}
	return contentBuffer.String()
}

// NewContentInterface return a map[string]interface{} with a title given as parameter
// It also set the language of the user and the list of available languages, as well as the page theme.
func NewContentInterface(pageTitleKey string, w http.ResponseWriter, r *http.Request) map[string]interface{} {
	ContentInterface := make(map[string]interface{})
	// Getting the user language
	currentLang := GetUserLang(r)
	langText, err := GetLangContent(currentLang)
	if err != nil {
		ErrorPrintf("An error occurred while trying to get the language content -> %v\n", err)
	} else {
		ContentInterface["Lang"] = langText
		ContentInterface["Title"] = langText["pageNames"].(map[string]interface{})[pageTitleKey]
	}
	// On va initialiser les listes de styles et de scripts supplémentaires.
	// Ces listes serviront à ajouter des styles et des scripts supplémentaires pour qu'ils soient chargés par le template.
	ContentInterface["AdditionalStyles"] = []string{}
	ContentInterface["AdditionalScripts"] = []string{}

	// Setting the language
	ContentInterface["LangList"] = LangListToStrList(langList)
	ContentInterface["CurrentLang"] = string(currentLang)

	// Setting the theme
	currentTheme := GetUserTheme(r)
	ContentInterface["CurrentTheme"] = string(currentTheme)

	return ContentInterface
}

// AddAdditionalScriptsToContentInterface add additional JS scripts to be loaded by the template.
func AddAdditionalScriptsToContentInterface(content *map[string]interface{}, scripts ...string) {
	for _, scriptName := range scripts {
		(*content)["AdditionalScripts"] = append((*content)["AdditionalScripts"].([]string), scriptName)
	}
}

// AddAdditionalStylesToContentInterface add additional CSS styles to be loaded by the template.
func AddAdditionalStylesToContentInterface(content *map[string]interface{}, styles ...string) {
	for _, styleName := range styles {
		(*content)["AdditionalStyles"] = append((*content)["AdditionalStyles"].([]string), styleName)
	}
}

// GetIP returns the IP address of the user.
func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

// InitServerCertification sets up the certification for the server.
// Required for OAuth and other services to work properly.
// If the CERT_FILE and CERT_KEY_FILE environment variables are not set, the server will run in HTTP mode.
// If they are set, the server will run in HTTPS mode.
func InitServerCertification() {
	if os.Getenv("CERT_FILE") == "" || os.Getenv("CERT_KEY_FILE") == "" {
		WarningPrintln("No certificate file or key file provided, the server will run in HTTP mode.")
		isCertified = false
	} else {
		SuccessPrintln("Certificate file and key file provided, the server will run in HTTPS mode.")
		isCertified = true
	}
	isInitialized = true
}

// LaunchServer launches the server with the given router and port.
// If the CERT_FILE and CERT_KEY_FILE environment variables are not set, the server will run in HTTP mode.
// If they are set, the server will run in HTTPS mode.
// r is the router to use for the server.
// port is the port to use for the server. It should be given as a string. (e.g. ":8080")
func LaunchServer(r *mux.Router, port string) {
	if !isInitialized {
		InitServerCertification()
		WarningPrintln("Server certification was initialized automatically, OAuth and other services may not work properly.")
	}
	// Launch the server
	if !isCertified {
		SuccessPrintf("Server started at -> http://localhost%s\n", port)
		if err := http.ListenAndServe(port, r); err != nil {
			panic(err)
		}
	} else {
		SuccessPrintf("Server started at -> https://localhost%s\n", port)
		if err := http.ListenAndServeTLS(port, os.Getenv("CERT_FILE"), os.Getenv("CERT_KEY_FILE"), r); err != nil {
			panic(err)
		}
	}
}

// IsCertified returns true if the server is running in HTTPS mode, false otherwise.
func IsCertified() bool {
	return isCertified
}

// ================================
// |Functions for inside the pages|
// ================================

// interfaceToString is a function to convert an interface to a string.
func interfaceToString(i interface{}) string {
	return i.(string)
}
