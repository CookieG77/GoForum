package backend

import (
	"GoForum/backend/pages"
	f "GoForum/functions"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
	"net/http"
	"os"
	"strconv"
)

// LaunchWebApp launches the web application and handles the routes.
func LaunchWebApp() {
	// Managing the program arguments
	f.AddValueArg(f.ArgIntValue, "port", "p")
	f.AddNoValueArg("debug", "d")
	if f.IsNoValueArgPresent("debug") || f.IsNoValueArgPresent("d") {
		f.SetShouldLogInfo(true)
	}
	finalPort := fmt.Sprintf(":%s", strconv.Itoa(getPort()))

	// Create the router
	r := mux.NewRouter()

	// Handle the static files
	r.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./statics/css"))))
	r.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./statics/img"))))
	r.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./statics/js"))))
	r.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir("./statics/fonts"))))

	// Handle the routes
	r.HandleFunc("/", pages.HomePage)

	// Getting the environment variables
	f.InfoPrintf("Getting the environment variables\n")
	err := godotenv.Load()
	if err != nil {
		f.ErrorPrintln("Error loading .env file")
	}

	// Creating the session store
	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	// Handle the OAuth routes
	f.SetupOAuth(finalPort)
	r.HandleFunc("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})

	// link the store to the gothic package
	gothic.Store = store

	// Handle the OAuth callback routes
	r.HandleFunc("/auth/callback/{provider}", func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			f.ErrorPrintf("Error while completing the user auth: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		f.SuccessPrintf("User connected !\n\t- Name : %v\n\t- Email : %v\n", user.Name, user.Email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Initialize the mail configuration
	f.InitMail("MailConfig.json")

	// Launch the server
	f.SuccessPrintf("Server started at -> http://localhost%s\n", finalPort)
	if err := http.ListenAndServe(finalPort, r); err != nil {
		panic(err)
	}
}

// getPort returns the port number to use for the server.
func getPort() int {
	strPort, err := f.GetArgValue("port", "p")
	if err != nil {
		f.ErrorPrintf("Error while getting the port value: %v\n", err)
	}
	f.InfoPrintf("Getting the port value: %v\n", strPort)
	var port int
	if strPort == nil { // If the port is not provided
		port = 8080
	} else {
		portInt, isAnInt := strPort.(int)
		if !isAnInt { // If the port is not an int
			f.ErrorPrintf("Error while converting the port value to int: %v\n", err)
			port = 8080
		}
		port = portInt
	}

	if port < 1 || port > 65535 {
		f.ErrorPrintf("The port %d is not a valid port number, switching back to default port (8080). Please provide a valid port number. (1-65535)\n", port)
		port = 8080
	}
	return port
}
