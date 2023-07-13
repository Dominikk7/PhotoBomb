package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Dominikk7/PhotoBomb/backend/pkg/database"

	"encoding/json"
	"math/rand"
	"time"
)


// For these tests a user exists in DB: username: "banana" password: "pass"

func TestMain(m *testing.M) {
	// have to start user db first
	database.Connect()
	database.Migrate()

	// seed rng
	rand.Seed(time.Now().UnixNano())

	// run all tests

	os.Exit(m.Run())
}

func genrandUsername() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	n := 9 //length
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestSignup(t *testing.T) {

	// generate ran username for creating new user
	userPayload, _ := json.Marshal(map[string]string{"password": "banana12", "username": genrandUsername()})
	payload2, _ := json.Marshal(map[string]string{"password": "banana*#&", "username": genrandUsername()})

	tt := []struct {
		name       string
		method     string
		body       string
		statusCode int
	}{
		{
			name:       "with empty request",
			method:     http.MethodPost,
			body:       ``,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with bad json",
			method:     http.MethodPost,
			body:       `{"pass":"banana12","user":"jsidj"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with empty user",
			method:     http.MethodPost,
			body:       `{"password":"banana12","username":""}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with empty password",
			method:     http.MethodPost,
			body:       `{"password":"","username":"banana"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with too short user",
			method:     http.MethodPost,
			body:       `{"password":"banana12","username":"n"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with too short password",
			method:     http.MethodPost,
			body:       `{"password":"ba","username":"tanner"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with symbols in user",
			method:     http.MethodPost,
			body:       `{"password":"banana12","username":"&&#*#%@&"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with symbols in password",
			method:     http.MethodPost,
			body:       string(payload2),
			statusCode: http.StatusOK,
		},
		{
			name:       "normal signup",
			method:     http.MethodPost,
			body:       string(userPayload),
			statusCode: http.StatusOK,
		},
		{
			name:       "already taken username",
			method:     http.MethodPost,
			body:       string(userPayload),
			statusCode: http.StatusConflict,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/signup", strings.NewReader(tc.body))
			responseRecorder := httptest.NewRecorder()

			Signup(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}
		})
	}
}

func TestSignin(t *testing.T) {

	// unknown username
	noUserPayload, _ := json.Marshal(map[string]string{"password": "banana12", "username": genrandUsername()})

	// wrong password, existing user
	badPassPayload, _ := json.Marshal(map[string]string{"password": "WRONG", "username": "banana"})

	// right password, existing user
	goodPassPayload, _ := json.Marshal(map[string]string{"password": "pass", "username": "banana"})

	tt := []struct {
		name       string
		method     string
		body       string
		statusCode int
	}{
		{
			name:       "with empty request",
			method:     http.MethodPost,
			body:       ``,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with bad json",
			method:     http.MethodPost,
			body:       `{"pass":"banana","user":"jsidj"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with empty user",
			method:     http.MethodPost,
			body:       `{"password":"banana","username":""}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "with empty password",
			method:     http.MethodPost,
			body:       `{"password":"","username":"banana"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "normal signin",
			method:     http.MethodPost,
			body:       string(goodPassPayload),
			statusCode: http.StatusOK,
		},
		{
			name:       "user doesn't exist",
			method:     http.MethodPost,
			body:       string(noUserPayload),
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "wrong password, user exists",
			method:     http.MethodPost,
			body:       string(badPassPayload),
			statusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/signin", strings.NewReader(tc.body))
			responseRecorder := httptest.NewRecorder()

			Signin(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}
		})
	}
}

func TestLogout(t *testing.T) {

	t.Run("no cookie", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/logout", nil)
		responseRecorder := httptest.NewRecorder()

		Logout(responseRecorder, request)

		if responseRecorder.Code != http.StatusUnauthorized {
			t.Errorf("Want status '%d', got '%d'", http.StatusUnauthorized, responseRecorder.Code)
		}
	})

	t.Run("bad cookie", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/logout", nil)
		responseRecorder := httptest.NewRecorder()

		http.SetCookie(responseRecorder, &http.Cookie{Name: "session_token", Value: "banana"})

		Logout(responseRecorder, request)

		if responseRecorder.Code != http.StatusUnauthorized {
			t.Errorf("Want status '%d', got '%d'", http.StatusUnauthorized, responseRecorder.Code)
		}
	})

	t.Run("good cookie", func(t *testing.T) {

		// first sign in to get proper cookie
		// using testing user

		goodPassPayload, _ := json.Marshal(map[string]string{"password": "pass", "username": "banana"}) // testing user

		signinRequest :=httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(string(goodPassPayload)))
		signinResponseRecorder := httptest.NewRecorder()

		Signin(signinResponseRecorder, signinRequest)

		// now extract cookie from response and use in new request

		request := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(""))
		responseRecorder := httptest.NewRecorder()
		
		request.AddCookie(signinResponseRecorder.Result().Cookies()[0])

		Logout(responseRecorder, request)

		if responseRecorder.Code != http.StatusOK {
			t.Errorf("Want status '%d', got '%d'", http.StatusOK, responseRecorder.Code)
		}
	})

}
