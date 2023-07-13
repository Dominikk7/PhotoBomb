package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Dominikk7/PhotoBomb/backend/pkg/database"
	"github.com/Dominikk7/PhotoBomb/backend/pkg/models"
	"github.com/auyer/steganography"
)

// Add image to file and db
func AddImage(token string, filetype string, file *bytes.Buffer, w http.ResponseWriter) {
	//filename = token + ID -> ID = timestamp
	timestamp := Timestamp()
	fileName := token + timestamp + filetype

	//Save file
	WriteFile(fileName, file)

	//Image object to DB
	w.Header().Set("Content-Type", "application/json")
	var image models.Image
	image.Token = token
	image.Timestamp = timestamp
	image.Extension = filetype
	database.ImageInstance.Create(&image)
	json.NewEncoder(w).Encode(image)
}

func DecodeImage(file *multipart.File) (string, error) {
	newImage, _, err := image.Decode(*file)
	if err != nil {
		return "", errors.New("failed converting multipartfile to image.Image")
	}
	
	size := steganography.GetMessageSizeFromImage(newImage)

	msg := steganography.Decode(size, newImage)

	return string(msg), nil
}

func DecodeImageBytes(file *bytes.Buffer) (string, error) {
	newImage, _, err := image.Decode(file)
	if err != nil {
		return "", errors.New("failed converting byte buffer to image.Image")
	}

	size := steganography.GetMessageSizeFromImage(newImage)

	msg := steganography.Decode(size, newImage)

	return string(msg), nil
}

func WriteFile(fileName string, file *bytes.Buffer) {
	//Write image file
	f, err := os.Create("../uploads/"+fileName)

	if err != nil {
		fmt.Println("error writing file")
		return
	}
	file.WriteTo(f) // write buffer to file
	f.Close()
}

func Timestamp() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
