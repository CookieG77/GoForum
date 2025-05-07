package pagesHandlers

import (
	f "GoForum/functions"
	"fmt"
	"net/http"
)

func ThreadCreationPage(w http.ResponseWriter, r *http.Request) {
	PageInfo := f.NewContentInterface("thread_creation", r)
	// Check the user rights
	f.GiveUserHisRights(&PageInfo, r)
	if PageInfo["IsAuthenticated"].(bool) {
		// If the user is not verified, redirect him to the verify page
		if !PageInfo["IsAddressVerified"].(bool) {
			f.InfoPrintf("Thread creation page accessed at %s by unverified : %s\n", f.GetIP(r), f.GetUserEmail(r))
			http.Redirect(w, r, "/confirmMail", http.StatusFound)
			return
		}
		f.InfoPrintf("Thread creation page accessed at %s by verified : %s\n", f.GetIP(r), f.GetUserEmail(r))
	} else {
		// If not authenticated, redirect to the login page
		f.InfoPrintf("Thread creation page accessed at %s\n", f.GetIP(r))
		http.Redirect(w, r, "/?openlogin=true", http.StatusFound)
	}

	// Handle the user logout/login
	ConnectFromHeader(w, r, &PageInfo)

	PageInfo["NameNotValid"] = false
	PageInfo["DescriptionNotValid"] = false
	PageInfo["ErrorCreationThread"] = false

	// Handle the thread creation form
	if r.Method == "POST" {
		// parse the form
		err := r.ParseForm()
		if err == nil {
			// Get the form values
			threadName := r.FormValue("thread_name")
			threadDescription := r.FormValue("thread_description")

			// Check if the thread name is valid
			if f.IsThreadNameValid(threadName) {
				// Check if the thread description is valid
				if f.IsThreadDescriptionValid(threadDescription) {
					// Create the thread
					err := f.AddThread(f.GetUser(r), threadName, threadDescription)
					if err != nil {
						f.ErrorPrintf("Error creating the thread : %s\n", err)
						PageInfo["ErrorCreationThread"] = true
					} else {
						f.InfoPrintf("Thread created with name : %s\n", threadName)
						// Redirect to the thread page
						http.Redirect(w, r, fmt.Sprintf("/t/%s", threadName), http.StatusFound)
					}
				} else {
					f.DebugPrintf("Thread description not valid : %s\n", threadDescription)
					PageInfo["DescriptionNotValid"] = true
				}
			} else {
				f.DebugPrintf("Thread name not valid : %s\n", threadName)
				PageInfo["NameNotValid"] = true
			}
		} else {
			f.ErrorPrintf("Error parsing the form : %s\n", err)
			PageInfo["ErrorCreationThread"] = true
		}

	}

	// Add additional styles to the content interface and make the template
	f.AddAdditionalStylesToContentInterface(&PageInfo, "css/threadCreation.css")
	f.MakeTemplateAndExecute(w, PageInfo, "templates/threadCreation.html")
}
