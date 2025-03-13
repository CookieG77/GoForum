package backend

import (
	"GoForum/backend/pages"
	f "GoForum/functions"
	"fmt"
	"net/http"
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

	// Handle the static files
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./statics/css"))))
	http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./statics/img"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./statics/js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir("./statics/fonts"))))

	// Handle the routes
	http.HandleFunc("/", pages.HomePage)

	// Initialize the mail configuration
	f.InitMail("MailConfig.json")

	// Launch the server
	f.SuccessPrintf("Server started at -> http://localhost%s\n", finalPort)
	if err := http.ListenAndServe(finalPort, nil); err != nil {
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
	port, err2 := strPort.(int)
	if !err2 { // If the port is not an int
		f.ErrorPrintf("Error while converting the port value to int: %v\n", err2)
		port = 8080
	}
	if strPort == nil { // If the port is not provided
		port = 8080
	}
	if port < 1 || port > 65535 {
		f.ErrorPrintf("The port %d is not a valid port number, switching back to default port (8080). Please provide a valid port number. (1-65535)\n", port)
		port = 8080
	}
	return port
}
