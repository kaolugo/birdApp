package main

// responsible for instantiating + configuring the birdApp object + its dependencies
// also starts the webserver listening

//"github.com/"
import (
	"fmt"
	"log"
	"net/http"

	"github.com/HENNGE/kaoru-BirdApp/internal/birdApp"
	"github.com/HENNGE/kaoru-BirdApp/redisDB"

	"github.com/gorilla/mux"
)

const listenPort = 8080
const redisPort = 6379

//const listenAddr = "localhost:8080"
//const redisAddr = "localhost:6379"

/* CORS Middleware */
func forCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		// プリフライトリクエストの対応
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}

func main() {
	// initialize new instances of the birdApp here and its dependencies

	// create a new instance of a mux router
	router := mux.NewRouter().StrictSlash(true)
	router.Use(forCORS)

	// the redis database initialization
	redisDatabase, err := redisDB.New("localhost:6379")

	// make sure that the connection to database was created successfully
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	// API initialization
	api := birdApp.New(redisDatabase)

	// call http.HandleFunc for all endpoints here
	router.HandleFunc("/test", api.Test)

	// user related http requests
	router.HandleFunc("/user", api.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", api.GetUser)
	router.HandleFunc("/allUsers", api.AllUsers)

	//router.HandleFunc("/follow/{id}", api.FollowUser).Methods("PUT")
	//router.HandleFunc("/unfollow/{id}", api.UnfollowUser).Methods("PUT")

	router.HandleFunc("/follow/{id}", api.FollowUser).Methods("POST")
	router.HandleFunc("/follow/{id}/{friendID}", api.VerifyFollow)

	router.HandleFunc("/unfollow/{id}", api.UnfollowUser).Methods("POST")
	router.HandleFunc("/unfollow", api.Test)

	// display personal timeline
	router.HandleFunc("/personal/{id}", api.ShowPersonal)
	// display global timeline
	router.HandleFunc("/global", api.ShowGlobal).Methods("GET", "OPTIONS")

	// tweet related http requests
	router.HandleFunc("/tweet/{id}", api.DeleteTweet).Methods("DELETE")
	router.HandleFunc("/tweet/{id}", api.EditTweet).Methods("PUT")
	router.HandleFunc("/tweet/{id}", api.NewTweet).Methods("POST")
	router.HandleFunc("/tweet/{id}", api.GetTweet)

	http.Handle("/", router)

	// add log.Printf statements and log.Fatal statements here
	log.Printf("Listening on %d\n", listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil))
}
