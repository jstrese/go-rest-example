package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jstrese.net/lib/respond"
	"log"
	"net/http"
	"strconv"
)

var ListenPort int = 10000

// Note: First letters must be capitalized or the struct won't be marshalled
type Post struct {
	Id    int    `json:"Id"`
	Title string `json:"Title"`
	Body  string `json:"Body"`
	User  string `json:"User"`
}

// Generic request body struct for all incoming requests
// Requests can use 1 or all of these attributes
type RequestBody struct {
	Limit int
	Page  int
	Query string
}

var Posts []Post

func testCall(w http.ResponseWriter, r *http.Request) {

}

// TODO Support adding optional limit parameter
func getPosts(w http.ResponseWriter, r *http.Request) {
	var body RequestBody

	if r.ContentLength > 0 {
		json.NewDecoder(r.Body).Decode(&body)
	}

	if posts, postsErr := GetPostList(body.Limit); postsErr == nil && posts != nil {
		json.NewEncoder(w).Encode(posts)
	} else {
		// Posts will be nil if there were 0 results
		if postsErr != nil {
			respond.Response(w, 500, postsErr)
		} else {
			respond.Response(w, 204)
		}
	}
}

// func addPost(w http.ResponseWriter, r *http.Request) {
// 	nextId := len(Posts) + 1
// 	Posts = append(Posts, Post{Id: nextId, Title: "Added", Content: "Dynamic content baby!", User: "RestUser"})
// 	getPosts(w, r)
// }

func getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil || id <= 0 {
		respond.BadRequest(w, "Malformed ID")
		return
	}

	if post, postErr := GetPostById(id); postErr == nil {
		json.NewEncoder(w).Encode(post)
	} else {
		if postErr == sql.ErrNoRows {
			respond.NotFound(w, "Post not found")
		} else {
			respond.Response(w, 500, postErr)
		}
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	// TEST GET
	router.HandleFunc("/test", testCall).Methods("GET")

	// GET
	router.HandleFunc("/get/all", getPosts).Methods("GET")
	router.HandleFunc("/get/{id}", getPost).Methods("GET")

	// POST
	//router.HandleFunc("/add", addPost).Methods("POST")

	println("Server listening on port " + fmt.Sprintf("%d", ListenPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ListenPort), router))
}

func main() {
	defer handleRequests()
	defer connect()

	// Add any other setup code that we may want here
	// The calls above are deferred and will execute after code below
}
