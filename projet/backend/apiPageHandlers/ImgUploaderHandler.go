package apiPageHandlers

import (
	f "GoForum/functions"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

// ImgUploader handles the image upload for the web application.
// Its path is "/api/upload/{type}" where {type} is the type of the image.
// It checks if the user is authenticated and verified before allowing the upload.
// It also checks if the file is a valid image and if it is not too large.
// The image is saved in the upload folder and a media link is created for it.
func ImgUploader(w http.ResponseWriter, r *http.Request) {
	f.DebugPrintln("ImgUploaderHandler called")

	// Get the query parameters
	vars := mux.Vars(r)
	rawImgType := vars["type"]

	// Check if the image type is valid
	if !f.IsAMediaType(rawImgType) {
		f.DebugPrintf("Image type \"%s\" is not valid\n", rawImgType)
		http.Error(w, "Image type is not valid", http.StatusBadRequest)
		return
	}

	imgType, err := f.GetMediaTypeFromString(rawImgType)
	if err != nil {
		f.DebugPrintf("Error getting media type from string: %s\n", err)
		http.Error(w, "Image type is not valid", http.StatusBadRequest)
		return
	}

	var user f.User
	var userConfigs f.UserConfigs
	var thread f.ThreadGoForum
	var threadConfigs f.ThreadGoForumConfigs
	if imgType == f.ThreadIcon || imgType == f.ThreadBanner {
		// Check if the thread is given
		threadName := r.FormValue("thread")
		if threadName == "" {
			f.DebugPrintf("Thread threadName is not given\n")
			http.Error(w, "Thread threadName is not given", http.StatusBadRequest)
			return
		}
		// Check if the thread exists
		thread = f.GetThreadFromName(threadName)
		if (thread == f.ThreadGoForum{}) {
			f.DebugPrintf("Thread \"%s\" does not exist\n", threadName)
			http.Error(w, "Thread does not exist", http.StatusBadRequest)
			return
		}
		user = f.GetUser(r)
		// Check if the user is the owner of the thread
		if thread.OwnerID != user.UserID {
			f.DebugPrintf("User is not the owner of the thread\n")
			http.Error(w, "User is not the owner of the thread", http.StatusForbidden)
			return
		}
		threadConfigs = f.GetThreadConfigFromThread(thread)
	} else if imgType == f.UserProfilePicture {
		// Check if the user is authenticated
		user = f.GetUser(r)
		if (user == f.User{}) {
			f.DebugPrintf("User is not authenticated\n")
			http.Error(w, "User is not authenticated", http.StatusUnauthorized)
			return
		}
		userConfigs = f.GetUserConfig(r)
	}

	// Check if the user is authenticated
	if !f.IsAuthenticated(r) {
		f.DebugPrintf("User is not authenticated\n")
		http.Error(w, "User is not authenticated", http.StatusUnauthorized)
		return
	}

	// Check if the user is verified
	if !f.IsUserVerified(r) {
		f.DebugPrintf("User is not verified\n")
		http.Error(w, "User is not verified", http.StatusUnauthorized)
		return
	}

	// Check if the request method is POST
	if r.Method != http.MethodPost {
		f.DebugPrintf("Request method is not POST\n")
		http.Error(w, "Request method is not POST", http.StatusMethodNotAllowed)
		return
	}

	// Check if the request is too large
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // 20 Mo max

	// Check if the request is multipart/form-data
	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		f.DebugPrintf("Error parsing multipart form: %s\n", err)
		http.Error(w, "Fichier trop volumineux", http.StatusRequestEntityTooLarge)
		return
	}

	// Check if the file is present in the form
	file, handler, err := r.FormFile("image")
	if err != nil {
		f.DebugPrintf("Error getting file from form: %s\n", err)
		http.Error(w, "Erreur de fichier", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			f.DebugPrintf("Error closing file: %s\n", err)
			http.Error(w, "Erreur de fichier", http.StatusInternalServerError)
			return
		}
	}(file)

	// Check if the file is empty
	if handler.Size == 0 {
		f.DebugPrintf("File is empty\n")
		http.Error(w, "Fichier vide", http.StatusBadRequest)
		return
	}

	// Check if the file is too large

	if handler.Size > 20<<20 {
		f.DebugPrintf("File is too large\n")
		http.Error(w, "Fichier trop volumineux", http.StatusRequestEntityTooLarge)
		return
	}

	// Check if the file is a valid image
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		f.DebugPrintf("Error reading file: %s\n", err)
		http.Error(w, "Erreur lecture", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(buf)
	allowedTypes := map[string]string{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
		"image/gif":  ".gif",
	}
	ext, ok := allowedTypes[contentType]
	if !ok {
		http.Error(w, "Format non autorisé", http.StatusUnsupportedMediaType)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		f.DebugPrintf("Error seeking file: %s\n", err)
		http.Error(w, "Erreur de fichier", http.StatusInternalServerError)
		return
	}

	uniqueName := uuid.New().String() + ext
	fullPath := path.Join(f.GetImgUploadFolder(), uniqueName)

	// Create the file
	dst, err := os.Create(fullPath)
	if err != nil {
		f.DebugPrintf("Error creating file: %s\n", err)
		http.Error(w, "Erreur de fichier", http.StatusInternalServerError)
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			f.DebugPrintf("Error closing file: %s\n", err)
			http.Error(w, "Erreur de fichier", http.StatusInternalServerError)
			return
		}
	}(dst)

	// Copy the file
	_, err = io.Copy(dst, file)
	if err != nil {
		f.DebugPrintf("Error copying file: %s\n", err)
		http.Error(w, "Erreur de fichier", http.StatusInternalServerError)
		return
	}

	f.DebugPrintf("File uploaded '%s' : %s\n", rawImgType, fullPath)

	// Creating the MediaLink
	mediaID, err := f.AddMediaLink(imgType, uniqueName)
	if err != nil {
		f.DebugPrintf("Error creating media link: %s\n", err)
		http.Error(w, "Erreur de fichier", http.StatusInternalServerError)
		return
	}

	switch imgType {
	case f.UserProfilePicture:
		// Update the user profile picture
		f.DebugPrintf("Changing user \"%s\" profile picture to %s\n", user.Username, uniqueName)
		userConfigs.PfpID = mediaID
		err := f.UpdateUserConfig(userConfigs)
		if err != nil {
			f.DebugPrintf("Error updating user configs: %s\n", err)
			http.Error(w, "Erreur de fichier", http.StatusInternalServerError)
			return
		}
	case f.ThreadBanner:
		// Update the thread banner
		f.DebugPrintf("Changing thread \"%s\" banner to %s\n", thread.ThreadName, uniqueName)
		threadConfigs.ThreadBannerID = mediaID
		err := f.UpdateThreadConfigs(threadConfigs)
		if err != nil {
			f.DebugPrintf("Error updating thread configs: %s\n", err)
			http.Error(w, "Erreur de fichier", http.StatusInternalServerError)
			return
		}
	case f.ThreadIcon:
		// Update the thread icon
		f.DebugPrintf("Changing thread \"%s\" icon to %s\n", thread.ThreadName, uniqueName)
		threadConfigs.ThreadIconID = mediaID
		err := f.UpdateThreadConfigs(threadConfigs)
		if err != nil {
			f.DebugPrintf("Error updating thread configs: %s\n", err)
			http.Error(w, "Erreur de fichier", http.StatusInternalServerError)
			return
		}
	default:
		break
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf(`{"url": "%s", "id": "%d"}`, uniqueName, mediaID)))
	if err != nil {
		f.DebugPrintf("Error writing response: %s\n", err)
		http.Error(w, "Erreur d'écriture", http.StatusInternalServerError)
		return
	}
}
