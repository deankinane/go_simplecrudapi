package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func main() {
	router := mux.NewRouter()

	initMoviesDb()

	router.HandleFunc("/movies", getAllMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies", deleteMovie).Methods("DELETE")

	fmt.Print("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func initMoviesDb() {
	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "123456789",
		Title: "Movie One",
		Director: &Director{
			FirstName: "Dean",
			LastName:  "Kinane",
		},
	})

	movies = append(movies, Movie{
		ID:    "2",
		Isbn:  "78945613",
		Title: "Movie Two",
		Director: &Director{
			FirstName: "Natasha",
			LastName:  "Kinane",
		},
	})
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(404)
	fmt.Fprintf(w, "No movie found for ID: %s", params["id"])
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var deleted Movie
	for index, item := range movies {
		if item.ID == params["id"] {
			deleted = movies[index]
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(deleted)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var newMovie Movie
	if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid request body")
		return
	}
	newMovie.ID = strconv.Itoa(rand.Intn(9999))
	movies = append(movies, newMovie)

	w.Header().Set("Location", fmt.Sprintf("http://localhost:8000/movies/%s", newMovie.ID))
	w.WriteHeader(201)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid request body")
		return
	}

	movie.ID = id
	found := false
	for idx, item := range movies {
		if item.ID == id {
			movies = append(movies[:idx], movies[idx+1:]...)
			movies = append(movies, movie)
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Movie not found for ID: %s", movie.ID)
	}
}
