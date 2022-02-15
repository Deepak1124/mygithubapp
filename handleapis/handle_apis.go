package handleapis

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GithubCreateRepo create a repo under auth user
func GithubCreateRepo(w http.ResponseWriter, r *http.Request) {
	var param createRepoParam
	client, err := createGithubClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	repo := github.Repository{
		Name: &param.RepoName,
	}
	log.Info("params in creating repo", param)
	repos, _, err := client.Repositories.Create(ctx, "", &repo)
	log.Info("repos in creating repos", repos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	jsonResponse, _ := json.Marshal(repos)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

//Github GetRepo endpoint to get repo
func GithubGetRepo(w http.ResponseWriter, r *http.Request) {
	client, err := createGithubClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	keys, ok := r.URL.Query()["name"]

	if !ok || len(keys[0]) < 1 {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	repoName := keys[0]

	if repoName == "" {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	userName, _, _ := client.Users.Get(ctx, "")
	repo, resp, err := client.Repositories.Get(ctx, *userName.Login, repoName)
	if err != nil {
		log.Error(err)
		if resp.StatusCode == 404 {
			http.Error(w, "Repo not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	}
	jsonResponse, _ := json.Marshal(repo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GithubCreateBranch endpoint to create branch
func GithubCreateBranch(w http.ResponseWriter, r *http.Request) {
	var param createBranchParam
	client, err := createGithubClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if param.SourceBranchName == "" {
		param.SourceBranchName = "main"
	}
	userName, _, _ := client.Users.Get(ctx, "")
	sourceBranchRefString := fmt.Sprintf("refs/heads/%s", param.SourceBranchName)
	sourceBranchRef, resp, err := client.Git.GetRef(ctx, *userName.Login, param.RepoName, sourceBranchRefString)
	if err != nil {
		log.Error(err)
		if resp.StatusCode == 404 {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "", http.StatusOK)
		}
	}
	newBranchRefString := fmt.Sprintf("refs/heads/%s", param.DestinationBranchName)
	*sourceBranchRef.Ref = newBranchRefString
	newBranchRef, resp, err := client.Git.CreateRef(ctx, *userName.Login, param.RepoName, sourceBranchRef)
	if err != nil {
		log.Error(err)
		if resp.StatusCode == 404 {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "", http.StatusOK)
		}
	}
	log.Info(fmt.Sprintf("Branch %s Created Successfully", param.DestinationBranchName))

	jsonResponse, _ := json.Marshal(newBranchRef)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// CreateGithubPullRequest endpoint to create pull request
func CreateGithubPullRequest(w http.ResponseWriter, r *http.Request) {
	var param createPullRequestParam
	client, err := createGithubClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	createPullRequest := &github.NewPullRequest{
		Title: &param.PrSubject,
		Head:  &param.Head,
		Base:  &param.Base,
	}
	userName, _, _ := client.Users.Get(ctx, "")
	repo, _, err := client.PullRequests.Create(ctx, *userName.Login, param.RepoName, createPullRequest)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	log.Info("PullRequest Created Successfully")

	jsonResponse, _ := json.Marshal(repo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// CreateRepositoryContent creates a file with base64 content and commit using the message
func CreateRepositoryContent(w http.ResponseWriter, r *http.Request) {
	var createFileParam createFileContentParam
	client, err := createGithubClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &createFileParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	repositoryFileContent := &github.RepositoryContentFileOptions{
		Message: &createFileParam.Message,
		Content: []byte(createFileParam.FileContent),
		Branch:  &createFileParam.BranchName,
	}
	if createFileParam.Path == "" {
		createFileParam.Path = createFileParam.FileName
	}
	userName, _, _ := client.Users.Get(ctx, "")
	content, resp, err := client.Repositories.CreateFile(ctx, *userName.Login, createFileParam.RepoName, createFileParam.Path,
		repositoryFileContent)
	if err != nil {
		log.Error(err.Error())
		if resp.StatusCode == http.StatusNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if resp.StatusCode == http.StatusConflict {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
	}
	jsonResponse, _ := json.Marshal(content)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
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
	http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
}

// // GithubLoginCallbackHandler to store access_token in session
func GithubLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	githubAccessToken, err := domain.GetGithubAccessToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	log.Info("Github access token found", githubAccessToken)
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

	log.Info("Call back handler done")
	http.Redirect(w, r, "/github/listRepo", http.StatusTemporaryRedirect)
}
