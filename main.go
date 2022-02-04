package main

import (
	"fmt"
	"mygithubapp/domain"
	"mygithubapp/handleapis"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
)

// var (
// 	key   = []byte("super-secret-key")
// 	Store = sessions.NewCookieStore(key)
// )

func initialize() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/github/listRepo", auth(handleapis.GithubListRepository)).Methods("GET")
	// r.HandleFunc("github/createRepo", handleapis.GithubCreateRepositoryController).Methods("POST")
	// r.HandleFunc("github/getRepo", handleapis.GithubGetRepositoryController).Methods("GET")
	// r.HandleFunc("github/createBranch", handleapis.GithubCreateBranchController).Methods("POST")
	// r.HandleFunc("github/createPullRequest", handleapis.GithubCreatePullRequestController).Methods("POST")
	// r.HandleFunc("github/createContent", handleapis.GithubCreateContentRequestController).Methods("POST")

	r.HandleFunc("/login", handleapis.GithubLogin).Methods("GET")
	r.HandleFunc("/callback/handler", handleapis.GithubLoginCallbackHandler).Methods("GET")

	// t.GET("login", infrastructure.GithubLoginController)
	// t.GET("callback/handler", infrastructure.GithubCallbackHandlerController)

	// t := e.Group("/")
	// t.Use(SessionAuth())
	// t.GET("github/listRepo", handleapis.GitListRepositoryController)
	// t.POST("github/createRepo", handleapis.GithubCreateRepositoryController)
	// t.GET("github/getRepo", handleapis.GithubGetRepositoryController)
	// t.POST("github/createBranch", handleapis.GithubCreateBranchController)
	// t.POST("github/createPullRequest", handleapis.GithubCreatePullRequestController)
	// t.POST("github/createContent", handleapis.GithubCreateContentRequestController)

	// 	Serve = r
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
			fmt.Println("acces token nahi mila in main55, so redirecting to login")
			http.Redirect(w, r, "/login", 302)
			return
		}
		fmt.Println("acces token mila in main59, so redirecting to listrepo")
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
