package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/auyer/steganography"
	"github.com/Dominikk7/PhotoBomb/backend/pkg/auth"
	"github.com/Dominikk7/PhotoBomb/backend/pkg/database"
	"github.com/Dominikk7/PhotoBomb/backend/pkg/models"
	"github.com/Dominikk7/PhotoBomb/backend/pkg/utils"
)

func ImageCreate(w http.ResponseWriter, r *http.Request) { // uploads image into db
	fmt.Println(r.Cookies())

	//Authenticate request
	userID, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}


	//Parse form data
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageText := r.FormValue("imagetext")

	//Only allow images
	filetype := filepath.Ext(handler.Filename)
	filetype = strings.ToLower(filetype)
	if filetype != ".jpeg" && filetype != ".png" && filetype != ".jpg" {
		//errNew = "The provided file format is not allowed. Please upload a JPEG/JPG or PNG image"
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	fmt.Println(imageText)
	var newImage image.Image

	if filetype == ".jpeg" || filetype == ".jpg" {
		newImage, err = jpeg.Decode(file)
	} else {
		newImage, err = png.Decode(file)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	buf := new(bytes.Buffer)
	if err = steganography.Encode(buf, newImage, []byte(imageText)); err != nil {
		// this will also catch case where the image is too small to hold the text
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// encoded image bytes now contained in byte buffer "buf"
	// now pass this along to AddImage

	utils.AddImage(userID, filetype, buf, w) //Write image file and add to DB

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated) //http_status
}

func ImageDecode(w http.ResponseWriter, r *http.Request) { // takes an image from client, returns decoded text; doesn't upload into db

	// does not require auth


	var imageCode string
	//Parse form data
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Only allow images
	filetype := filepath.Ext(handler.Filename)
	if filetype != ".jpeg" && filetype != ".png" && filetype != ".jpg" {
		//errNew = "The provided file format is not allowed. Please upload a JPEG,JPG or PNG image"
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageCode, err = utils.DecodeImage(&file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/JSON; charset=UTF-8") // CHANGED TO JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(imageCode)
}

func GetImageById(w http.ResponseWriter, r *http.Request) { // returns an image file based on db ID

	//Authenticate request
	userID, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Parse request
	timestamp := r.URL.Query().Get("timestamp")

	var image models.Image
	image.Token = userID
	image.Timestamp = timestamp

	//Get from db
	if err := database.ImageInstance.Where("token = ? AND timestamp = ?", image.Token, image.Timestamp).First(&image).Error; err != nil { //If image does not exist
		w.WriteHeader(http.StatusNotFound)
		return
	}
	filename := image.Token + image.Timestamp + image.Extension
	fileBytes, err := ioutil.ReadFile("../uploads/" + filename)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func ExistingDecode(w http.ResponseWriter, r *http.Request) { // decodes an image from database, returns text to user

	//Authenticate request
	userID, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}


	//Parse request, get details of desired image
	var image models.Image
	timestamp := r.URL.Query().Get("timestamp")
	image.Timestamp = timestamp
	image.Token = userID

	//Get from db
	if err := database.ImageInstance.Where("token = ? AND timestamp = ?", image.Token, image.Timestamp).First(&image).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest) // image doesn't exist
		return
	}

	filename := image.Token + image.Timestamp + image.Extension
	fileBytes, err := ioutil.ReadFile("../uploads/" + filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // image in db but not in filesystem
		return
	}

	// now actually decode image
	imageText, err := utils.DecodeImageBytes(bytes.NewBuffer(fileBytes))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8") // CHANGED TO JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(imageText)

}

func GetAllImages(w http.ResponseWriter, r *http.Request) { // returns all images associated with a user
	//Authenticate request
	userID, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var image models.Image
	image.Token = userID

	//Get from db
	var images []models.Image
	if err := database.ImageInstance.Where("token = ?", image.Token).Find(&images).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if len(images) == 0 { // previous function was not catching all instances with no images
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(images)
}

func DeleteImageById(w http.ResponseWriter, r *http.Request) { // deletes an image from the database

	//Authenticate request
	userID, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Parse request
	var image models.Image
	timestamp := r.URL.Query().Get("timestamp")
	image.Timestamp = timestamp
	image.Token = userID
	
	//First find in db
	if err = database.ImageInstance.Where("token = ? AND timestamp = ?", image.Token, image.Timestamp).First(&image).Error; err != nil {
		// image not found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// now delete from db
	if err = database.ImageInstance.Delete(&image).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now delete from OS
	filename := image.Token + image.Timestamp + image.Extension
	fmt.Println(filename)

	if err = os.Remove("../uploads/" + filename); err != nil {
		// database and OS desynced, or some other issue with deleting file (file open/being used), should not happen
		fmt.Println("ERROR: COULD NOT DELETE FILE FROM OS")
	}

	w.WriteHeader(http.StatusOK)
}
