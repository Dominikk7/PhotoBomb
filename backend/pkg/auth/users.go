package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Dominikk7/PhotoBomb/backend/pkg/database"
	"github.com/Dominikk7/PhotoBomb/backend/pkg/models"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
)

type Credentials struct {
	Password string `json:"password" validate:"required,ascii,min=8,max=20"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=20"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()

	// Parse and decode the request body into a new `Credentials` instance
	creds := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	
	// check against username and password requirements
	p := bluemonday.UGCPolicy()
	if errr := validate.Struct(creds); errr != nil {
		// username/password does not meet requirements
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	creds.Username = p.Sanitize(creds.Username)
	creds.Password = p.Sanitize(creds.Password)

	// check for dupe username
	var foundUser models.User
	if err := database.UserInstance.Where("username = ?", creds.Username).First(&foundUser).Error; err == nil {
		// username already exists
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	if err != nil {
		// something wrong
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Next, create a new user entry for the database
	defaultTime := time.Now()

	// ADD IN SESSION AND SESHEXP
	newUser := models.User{Username: creds.Username, HashWord: string(hashedPassword), SeshExp: defaultTime}

	// Now add to DB

	result := database.UserInstance.Create(&newUser)

	// check for error
	if result.Error != nil {
		// something wrong
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	w.WriteHeader(http.StatusOK)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check for missing user/password
	if creds.Username == "" || creds.Password == "" {
		// missing username or password
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// First look for user
	var foundUser models.User

	if err := database.UserInstance.First(&foundUser, "username = ?", creds.Username).Error; err != nil {
		// username not found
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	storedPassword := foundUser.HashWord

	// attempt to decrypt password

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		// wrong password
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new random session token
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(2419200 * time.Second) // expiration date is 28 days

	// set token and expiration in database
	database.UserInstance.Model(&foundUser).Update("Session", sessionToken)
	database.UserInstance.Model(&foundUser).Update("SeshExp", expiresAt)

	//Construct object
	var currentCookie models.CookieModel
	currentCookie.Name = "session_token"
	currentCookie.Value = sessionToken
	currentCookie.Expires = expiresAt.String()

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
  
	json.NewEncoder(w).Encode(currentCookie)


	w.WriteHeader(http.StatusOK)
}



func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// find session token in DB
	var foundUser models.User
	if err := database.UserInstance.First(&foundUser, "session = ?", c.Value).Error; err != nil {
		// bad token
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// remove session from DB
	database.UserInstance.Model(&foundUser).Update("Session", "")
	database.UserInstance.Model(&foundUser).Update("SeshExp", "")

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
}

func GetUser(r *http.Request) (string, error) {
	cooky, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "-1", errors.New("no cookie")
		}
		return "-1", errors.New("bad request")
	}
	var foundUser models.User

	if err := database.UserInstance.First(&foundUser, "session = ?", cooky.Value).Error; err != nil {
		// user id doesn't exist
		return "-1", errors.New("userID doesn't exist")
	}
	//Convert to string
	userIDs := strconv.FormatUint(foundUser.ID, 10)

	return userIDs, nil
}
