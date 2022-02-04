package handleapis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"mygithubapp/domain"

	"github.com/google/go-github/github"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
)

//githubRepositoryListOutput struct for listing repository
type githubRepositoryListOutput struct {
	Name        string `json:"name"`
	FullName    string `json:"fullName"`
	Description string `json:"description"`
}

//createBranchParam struct for creating branch
type createBranchParam struct {
	RepoName              string `json:"repoName" validate:"required"`
	SourceBranchName      string `json:"sourceBranchName" default:"main"`
	DestinationBranchName string `json:"destinationBranchName" validate:"required"`
	Private               bool   `json:"private"`
}

//createRepoParam struct for creating repo
type createRepoParam struct {
	RepoName    string `json:"repoName" validate:"required"`
	Private     bool   `json:"private"`
	Description string `json:"description, omitempty"`
}

//createPullRequestParam struct for creating pull request
type createPullRequestParam struct {
	RepoName  string `json:"repoName" validate:"required"`
	PrSubject string `json:"prSubject" validate:"required"`
	Head      string `json:"head" validate:"required"`
	Base      string `json:"base" validate:"required"`
}

//createFileContentParam struct for creating a file in github
type createFileContentParam struct {
	RepoName    string `json:"repoName" validate:"required"`
	BranchName  string `json:"branchName" validate:"required"`
	FileName    string `json:"fileName" validate:"required"`
	FileContent string `json:"fileContent" validate:"required"`
	Path        string `json:"path" validate:"required"`
	Message     string `json:"message" validate:"required"`
}

var ctx = context.Background()

// createGithubClient return a github client
func createGithubClient(r *http.Request) (*github.Client, error) {
	// sess, err := session.Get("session", c)
	sess, err := domain.Get(r)
	if err != nil {
		panic(err)
	}
	accessToken := sess.Values["access_token"]
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken.(string)},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

// GithubListRepository endpoint for listing repository
func GithubListRepository(w http.ResponseWriter, r *http.Request) {
	client, err := createGithubClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	repos, _, err := client.Repositories.List(ctx, "", nil)
	// TODO get branch and add to output list
	//repo, resp, err := client.Repositories.GetBranch(ctx, userName, param.Name)
	repoList := make([]githubRepositoryListOutput, len(repos))
	i := 0
	for j := range repos {
		repoList[i].Name = repos[j].GetName()
		repoList[i].FullName = repos[j].GetFullName()
		repoList[i].Description = repos[j].GetDescription()
		i++
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	jsonResponse, _ := json.Marshal(repoList)
	// return c.JSON(http.StatusOK, repoList)
	fmt.Println(repoList)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GithubCreateRepo create a repo under auth user
func GithubCreateRepo(w http.ResponseWriter, r *http.Request) error {
	var param createRepoParam
	client, err := createGithubClient(r)
	if err != nil {
		return err
	}
	err = c.Bind(&param)
	if err != nil {
		return err
	}
	err = c.Validate(param)
	if err != nil {
		return err
	}
	repo := github.Repository{
		Name: &param.RepoName,
	}
	repos, _, err := client.Repositories.Create(ctx, "", &repo)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, repos)
}

// GithubLogin is to login in github using oauth
func GithubLogin(w http.ResponseWriter, r *http.Request) {
	githubClientID, err := domain.GetGithubClientID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=repo", githubClientID,
	)
	log.Info("Redirect URL", redirectURL)
	// http.Redirect(res, req, fmt.Sprintf("https://%s%s", req.Host, req.URL), http.StatusPermanentRedirect)
	http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
	// return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// // GithubLoginCallbackHandler to store access_token in session
func GithubLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	githubAccessToken, err := domain.GetGithubAccessToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	log.Info("Github access token mil gayaa in GithubLoginCallbackHandler 135", githubAccessToken)
	sess, err := domain.Get(r)
	if err != nil {
		panic(err)
	}

	sess.Values["access_token"] = githubAccessToken
	// sess.Save()
	err = sess.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// sess.Save(c.Request(), c.Response())
	log.Info("Call back handler done")
	http.Redirect(w, r, "/github/listRepo", http.StatusTemporaryRedirect)
}
