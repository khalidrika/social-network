package utils

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// Patht to default image location
const defaultImage = "imageUpload/default.svg"

// Creates new file and reads image bytes into it
// returns path to new image
// returns default avatar or provided one
func SaveAvatar(r *http.Request) string {
	// Read data from request
	file, fileHeader, errRead := r.FormFile("avatar")
	if errRead != nil {
		return defaultImage
	}
	defer file.Close()
	// get content type -> png, gif or jpeg
	contentType := fileHeader.Header["Content-Type"][0]
	// Create empty local file with correct file extension
	localFile, err := createTempFile(contentType)
	// If type not recognized return default image path
	if err != nil {
		return defaultImage
	}
	defer localFile.Close()

	// read data in new file
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return defaultImage
	}
	localFile.Write(fileData)
	return strings.Replace(localFile.Name(), "\\", "/", -1)
}

/* ------------------------- for posts and comments ------------------------- */
// Creates new file and reads image bytes into it
// returns path to new image
// returns no path if not exist
func SaveImage(r *http.Request) string {
	// Read data from request
	file, fileHeader, errRead := r.FormFile("image")
	if errRead != nil {
		return ""
	}
	defer file.Close()
	// get content type -> png, gif or jpeg
	contentType := fileHeader.Header["Content-Type"][0]
	// Create empty local file with correct file extension
	localFile, err := createTempFile(contentType)
	// If type not ecognized return default image path
	if err != nil {
		return ""
	}
	defer localFile.Close()

	// read data in new file
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return ""
	}
	localFile.Write(fileData)
	return strings.Replace(localFile.Name(), "\\", "/", -1)
}

// creates empty local file based on filt type
func createTempFile(fileType string) (*os.File, error) {
	var localFile *os.File
	var err error

	switch fileType {
	case "image/jpeg":
		localFile, err = ioutil.TempFile("imageUpload", "*.jpg")
	case "image/png":
		localFile, err = ioutil.TempFile("imageUpload", "*.png")
	case "image/gif":
		localFile, err = ioutil.TempFile("imageUpload", "*.gif")
	default:
		return localFile, err
	}
	return localFile, nil
}

// SaveUploadedFile saves an uploaded file to the specified directory
// returns the file path relative to the project root
func SaveUploadedFile(file multipart.File, fileHeader *multipart.FileHeader, directory string) (string, error) {
	// get content type -> png, gif or jpeg
	contentType := fileHeader.Header["Content-Type"][0]

	// Create empty local file with correct file extension in specified directory
	localFile, err := createTempFileInDir(contentType, directory)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	// read data in new file
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	_, err = localFile.Write(fileData)
	if err != nil {
		return "", err
	}

	return strings.Replace(localFile.Name(), "\\", "/", -1), nil
}

// creates empty local file based on file type in specified directory
func createTempFileInDir(fileType, directory string) (*os.File, error) {
	var localFile *os.File
	var err error

	switch fileType {
	case "image/jpeg":
		localFile, err = ioutil.TempFile(directory, "*.jpg")
	case "image/png":
		localFile, err = ioutil.TempFile(directory, "*.png")
	case "image/gif":
		localFile, err = ioutil.TempFile(directory, "*.gif")
	default:
		return localFile, err
	}
	return localFile, nil
}
