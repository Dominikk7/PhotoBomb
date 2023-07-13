package main

import (
	"log"
	"net/http"

	"github.com/Dominikk7/PhotoBomb/backend/pkg/controllers"
	"github.com/Dominikk7/PhotoBomb/backend/pkg/database"

	"github.com/Dominikk7/PhotoBomb/backend/pkg/auth"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"os"
)

const PORT = "8080"

func main() {
	//DB password from cli
	var sqlPass string
	if (len(os.Args) < 2){
		sqlPass = ""
	} else{
		sqlPass = os.Args[1]
	}
	
	// Initialize Database
	database.Connect(sqlPass)
	database.Migrate()

	//Init router
	r := mux.NewRouter()

	//Routes
	RegisterRoutes(r)

	//Start server
	log.Printf("Server is running on http://localhost:%s", PORT)

	//Cors allow all orgins
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	log.Println(http.ListenAndServe(":"+PORT, handler))

}
func RegisterRoutes(r *mux.Router) {
	//User login
	r.HandleFunc("/signin1", auth.Signin).Methods("POST")
	r.HandleFunc("/signup1", auth.Signup).Methods("POST")
	r.HandleFunc("/logout", auth.Logout).Methods("POST")

	//Image upload
	r.HandleFunc("/upload/encode", controllers.ImageCreate).Methods("POST")
	r.HandleFunc("/upload/decode", controllers.ImageDecode).Methods("POST") //Decode new
	r.HandleFunc("/decode/", controllers.ExistingDecode).Methods("GET")     //Decode existing
	r.HandleFunc("/download/", controllers.GetImageById).Methods("GET")
	r.HandleFunc("/download/list/", controllers.GetAllImages).Methods("GET")
	r.HandleFunc("/delete/", controllers.DeleteImageById).Methods("DELETE")
}
