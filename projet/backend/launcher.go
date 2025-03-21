package backend

import (
	"GoForum/backend/pagesHandlers"
	f "GoForum/functions"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

// LaunchWebApp launches the web application and handles the routes.
func LaunchWebApp() {
	// Getting the environment variables
	f.DebugPrintf("Getting the environment variables\n")
	err := godotenv.Load()
	if err != nil {
		f.ErrorPrintln("Error loading .env file")
	} else {
		f.SuccessPrintln("Environment variables loaded")
	}

	/* // Initialize the database
	f.InitDatabaseConnection()
	err = f.AddUserToModeration("mcordonn9277@gmail.com", 2)
	if err != nil {
		f.ErrorPrintf("Error while adding the user to the moderation: %v\n", err)
		return
	} */

	// Gestion de l'arrêt de l'application web via le terminal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				f.ClearCmd()
				f.CloseDatabase()
				os.Exit(1)
			}
		}
	}()

	// Managing the program arguments
	f.AddValueArg(f.ArgIntValue, "port", "p") // Argument to change the port
	f.AddNoValueArg("debug", "d")             // Argument to enable the debug mode
	f.AddNoValueArg("log", "l")               // Argument to enable the log mode
	if isPresent, err := f.GetArgNoValue("debug", "d"); isPresent && err == nil {
		f.SetShouldLogDebug(true)
	}
	if isPresent, err := f.GetArgNoValue("log", "l"); isPresent && err == nil {
		f.InitLogger()
	}
	finalPort := fmt.Sprintf(":%s", strconv.Itoa(getPort()))

	// Create the router
	r := mux.NewRouter()

	// Handle the static files
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("./statics/css"))))
	r.PathPrefix("/img/").Handler(http.StripPrefix("/img", http.FileServer(http.Dir("./statics/img"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("./statics/js"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts", http.FileServer(http.Dir("./statics/fonts"))))

	// Set the base template
	f.AddBaseTemplate("templates/base.html")

	// Handle the routes
	r.HandleFunc("/", pagesHandlers.HomePage)
	r.HandleFunc("/login", pagesHandlers.LoginPage)
	r.HandleFunc("/register", pagesHandlers.RegisterPage)
	r.HandleFunc("/profile", pagesHandlers.UserProfilePage)
	r.HandleFunc("/settings", pagesHandlers.UserSettingsPage)
	r.HandleFunc("/reset-password", pagesHandlers.ResetPasswordPage)
	r.HandleFunc("/confirm-email-address", pagesHandlers.ConfirmMailPage)

	// Handle error 404 & 405
	r.NotFoundHandler = http.HandlerFunc(pagesHandlers.ErrorPage404)
	r.MethodNotAllowedHandler = http.HandlerFunc(pagesHandlers.ErrorPage405)

	// Creating the session store
	f.SetupCookieStore()

	// Initialize the certificate
	f.InitServerCertification()

	// Initialize the OAuth keys and routes
	f.InitOAuthKeys(finalPort, r)

	// Initialize the mail configuration
	f.InitMail("MailConfig.json")

	// Launch the server
	f.LaunchServer(r, finalPort)
}

// getPort returns the port number to use for the server.
func getPort() int {
	strPort, err := f.GetArgValue("port", "p")
	if err != nil {
		f.ErrorPrintf("Error while getting the port value: %v\n", err)
	}
	f.DebugPrintf("Getting the port value: %v\n", strPort)
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
