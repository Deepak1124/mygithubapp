package main

import (
	"fmt"
	"mygithubapp/domain"
	"mygithubapp/handleapis"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
)

func initialize() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/github/listRepo", auth(handleapis.GithubListRepository)).Methods("GET")
	r.HandleFunc("/github/createRepo", auth(handleapis.GithubCreateRepo)).Methods("POST")
	r.HandleFunc("/github/getRepo", auth(handleapis.GithubGetRepo)).Methods("GET")
	r.HandleFunc("/github/createBranch", auth(handleapis.GithubCreateBranch)).Methods("POST")
	r.HandleFunc("/github/createPullRequest", auth(handleapis.CreateGithubPullRequest)).Methods("POST")
	r.HandleFunc("/github/createContent", auth(handleapis.CreateRepositoryContent)).Methods("POST")

	r.HandleFunc("/login", handleapis.GithubLogin).Methods("GET")
	r.HandleFunc("/callback/handler", handleapis.GithubLoginCallbackHandler).Methods("GET")

	return r
}

func auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := domain.Get(r)
		if err != nil {
			panic(err)
		}
		_, ok := sess.Values["access_token"]
		if !ok {
			fmt.Println("acces token found, so redirecting to login")
			http.Redirect(w, r, "/login", 302)
			return
		}
		fmt.Println("acces token not found, so redirecting to listrepo")
		HandlerFunc.ServeHTTP(w, r)
	}
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.Info("Panic in go-git initialization.", "error", err)
		}
	}()

	r := initialize()
	log.Info("Server Started on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
