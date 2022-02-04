package domain

import (
	"errors"
	"log"
	"net/http"
	"os"

	gsessions "github.com/gorilla/sessions"
)

var (
	key   = []byte(os.Getenv("SESSION_KEY"))
	store = gsessions.NewCookieStore(key)
)

func Get(req *http.Request) (*gsessions.Session, error) {
	return store.Get(req, "session")
}

func GetNamed(req *http.Request, name string) (*gsessions.Session, error) {
	return store.Get(req, name)
}

// GetGithubClientSecret is to fetch GITHUB_CLIENT_SECRET env variable
func GetGithubClientSecret() (string, error) {

	githubClientSecret, exists := os.LookupEnv("GITHUB_CLIENT_SECRET")
	if !exists {
		log.Fatal("Github ClientSecret not defined in .env file")
		return "", errors.New("GITHUB_CLIENT_SECRET is not define")
	}

	return githubClientSecret, nil
}

// GetGithubClientID is to fetch GITHUB_CLIENT_ID env variable
func GetGithubClientID() (string, error) {

	githubClientID, exists := os.LookupEnv("GITHUB_CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
		return "", errors.New("GITHUB_CLIENT_ID is not define")
	}

	return githubClientID, nil
}

// GetGithubClientID is to fetch GITHUB_CLIENT_ID env variable
func GetGithubAccessToken() (string, error) {

	githubAccessToken, exists := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
		return "", errors.New("GITHUB_CLIENT_ID is not define")
	}

	return githubAccessToken, nil
}
