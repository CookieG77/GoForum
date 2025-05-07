package functions

import (
	"os"
	"path/filepath"
	"strings"
)

// uploadFolder is the folder where all the uploads are stored
var uploadFolder = "uploads/"

// imgUploadSubFolder is the subfolder where all the images uploads are stored this folder is created in the uploadFolder
var imgUploadSubFolder = "img/"

// InitUploadsDirectory initializes the image upload folder
// Also change the upload folder if the environment variable is set
func InitUploadsDirectory() {

	// Get the upload folders from the environment variable
	envUploadFolder := os.Getenv("UPLOAD_FOLDER")
	if envUploadFolder != "" {
		uploadFolder = envUploadFolder
		InfoPrintf("uploadFolder found in .env variable set to %s\n", uploadFolder)
	}

	envImgUploadFolder := os.Getenv("IMG_UPLOAD_FOLDER")
	if envImgUploadFolder != "" {
		imgUploadSubFolder = envImgUploadFolder
		InfoPrintf("imgUploadSubFolder found in .env variable set to %s\n", imgUploadSubFolder)
	}

	// Create the upload folder if it doesn't exist
	if _, err := os.Stat(uploadFolder); os.IsNotExist(err) {
		err := os.Mkdir(uploadFolder, os.ModePerm)
		if err != nil {
			ErrorPrintf("Error creating upload folder: %s\n", err)
			return
		}
	}

	// Create the image upload folder if it doesn't exist
	if _, err := os.Stat(filepath.Join(uploadFolder, imgUploadSubFolder)); os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(uploadFolder, imgUploadSubFolder), os.ModePerm)
		if err != nil {
			ErrorPrintf("Error creating image upload folder: %s\n", err)
			return
		}
	}

	InfoPrintf("Upload folder initialized\n")
}

// GetUploadFolder returns the path to the upload folder
func GetUploadFolder() string {
	return uploadFolder
}

// GetImgUploadSubFolder returns the path to the image upload subfolder from the upload folder
func GetImgUploadSubFolder() string {
	return imgUploadSubFolder
}

// GetImgUploadFolder returns the path to the image upload folder
func GetImgUploadFolder() string {
	return strings.ReplaceAll(filepath.Join(uploadFolder, imgUploadSubFolder), "\\", "/")
}

// RemoveImg removes the image from the upload folder
// It returns true if the image was removed successfully, false otherwise
func RemoveImg(path string) bool {
	// Remove the image from the upload folder
	err := os.Remove(path)
	if err != nil {
		ErrorPrintf("Error removing image: %s\n", err)
		return false
	}
	InfoPrintf("Image removed: %s\n", path)
	return true
}
