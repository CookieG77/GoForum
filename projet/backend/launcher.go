package backend

import (
	"GoForum/backend/apiPageHandlers"
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

	// Initialize the default language
	f.InitDefaultLangConfig()

	// Initialize the default theme
	f.InitDefaultThemeConfig()

	// Initialize the Uploads directory
	f.InitUploadsDirectory()

	// Initialize the database
	f.InitDatabaseConnection()

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
	f.AddNoValueArg("debug", "d") // Argument to enable the debug mode
	f.AddNoValueArg("log", "l")   // Argument to enable the log mode
	if isPresent, err := f.GetArgNoValue("debug", "d"); isPresent && err == nil {
		f.SetShouldLogDebug(true)
	}
	if isPresent, err := f.GetArgNoValue("log", "l"); isPresent && err == nil {
		f.InitLogger()
	}
	finalPort := fmt.Sprintf(":%s", strconv.Itoa(getPort()))

	// Setting up the rate limiter
	// TODO : Add the rate limiter and his middleware

	// Create the router
	r := mux.NewRouter()

	// Handle the static files
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("./statics/css"))))
	r.PathPrefix("/img/").Handler(http.StripPrefix("/img", http.FileServer(http.Dir("./statics/img"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("./statics/js"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts", http.FileServer(http.Dir("./statics/fonts"))))
	// Upload folder
	r.PathPrefix("/upload/").Handler(http.StripPrefix("/upload", http.FileServer(http.Dir(fmt.Sprintf("./%s", f.GetImgUploadFolder())))))

	// Set the base template
	f.AddBaseTemplate("templates/base.html")

	// Handle the routes
	r.HandleFunc("/", pagesHandlers.HomePage).Methods("GET", "POST")
	r.HandleFunc("/register", pagesHandlers.RegisterPage).Methods("GET", "POST")
	r.HandleFunc("/auth/callback/{provider}", pagesHandlers.CallbackRedirection).Methods("GET", "POST")
	r.HandleFunc("/profile", pagesHandlers.UserSelfProfilePage).Methods("GET", "POST")
	r.HandleFunc("/profile/{user}", pagesHandlers.UserOtherProfilePage).Methods("GET", "POST")
	r.HandleFunc("/settings", pagesHandlers.UserSettingsPage).Methods("GET", "POST")
	r.HandleFunc("/reset-password", pagesHandlers.ResetPasswordPage).Methods("GET", "POST")
	r.HandleFunc("/confirm-email-address", pagesHandlers.ConfirmMailPage).Methods("GET", "POST")
	r.HandleFunc("/nt", pagesHandlers.ThreadCreationPage).Methods("GET", "POST")
	r.HandleFunc("/t/{threadName}", pagesHandlers.ThreadPage).Methods("GET", "POST")
	r.HandleFunc("/t/{threadName}/edit", pagesHandlers.ThreadEditPage).Methods("GET", "POST")
	r.HandleFunc("/t/{threadName}/p/{post}", pagesHandlers.ThreadPostPage).Methods("GET", "POST")
	r.HandleFunc("/tnm", pagesHandlers.ThreadSendMessagePage).Methods("GET", "POST")
	r.HandleFunc("/api/messages", apiPageHandlers.ThreadMessageGetter).Methods("GET")
	r.HandleFunc("/api/comments", apiPageHandlers.MessageCommentGetter).Methods("GET")
	r.HandleFunc("/api/threadTags", apiPageHandlers.ThreadTagsGetterHandler).Methods("GET")
	r.HandleFunc("/api/thread/{threadName}/{action}", apiPageHandlers.ThreadContentHandler).Methods("POST")
	r.HandleFunc("/api/upload/{type}", apiPageHandlers.ImgUploader).Methods("POST")

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
	f.InitMail()

	// Launch the server
	f.LaunchServer(r, finalPort)
}

// getPort returns the port number to use for the server.
// Get it from the environment variable
func getPort() int {
	strPort := os.Getenv("PORT")
	f.DebugPrintf("Getting the port value: %v\n", strPort)
	var port int
	if strPort == "" { // If the port is not provided
		f.ErrorPrintf("PORT environment variable not set, switching to default '8080'\n")
		port = 8080
	} else { // If the port is provided
		portInt, err := strconv.Atoi(strPort)
		if err != nil {
			f.ErrorPrintf("Error while converting the port value to int: %v\n", err)
			f.ErrorPrintf("Switching to default port '8080'\n")
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
