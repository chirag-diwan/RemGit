package githubapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}

type Owner struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	NodeID    string `json:"node_id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Fork     bool   `json:"fork"`

	Owner Owner `json:"owner"`

	Description *string `json:"description"`

	HTMLURL    string `json:"html_url"`
	APIURL     string `json:"url"`
	CommitsURL string `json:"commits_url"`

	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
	GitURL   string `json:"git_url"`
	SVNURL   string `json:"svn_url"`

	Homepage *string `json:"homepage"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PushedAt  time.Time `json:"pushed_at"`

	Size     int     `json:"size"`
	Language *string `json:"language"`

	StargazersCount int `json:"stargazers_count"`
	WatchersCount   int `json:"watchers_count"`
	ForksCount      int `json:"forks_count"`
	OpenIssuesCount int `json:"open_issues_count"`

	HasIssues      bool `json:"has_issues"`
	HasProjects    bool `json:"has_projects"`
	HasWiki        bool `json:"has_wiki"`
	HasPages       bool `json:"has_pages"`
	HasDiscussions bool `json:"has_discussions"`

	Archived bool `json:"archived"`
	Disabled bool `json:"disabled"`

	License *License `json:"license"`

	Topics []string `json:"topics"`

	Visibility    string `json:"visibility"`
	DefaultBranch string `json:"default_branch"`

	Score float64 `json:"score,omitempty"`
}

type Verification struct {
	Verified   bool   `json:"verified"`
	Reason     string `json:"reason"`
	Signature  string `json:"signature"`
	Payload    string `json:"payload"`
	VerifiedAt string `json:"verified_at"`
}

type Tree struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type CommitAuthor struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type Commit struct {
	Author    CommitAuthor `json:"author"`
	Committer CommitAuthor `json:"committer"`

	Message string `json:"message"`

	Tree Tree `json:"tree"`

	URL          string `json:"url"`
	CommentCount int    `json:"comment_count"`

	Verification Verification `json:"verification"`
}

type CommitItem struct {
	SHA    string `json:"sha"`
	NodeID string `json:"node_id"`

	Commit Commit `json:"commit"`
}

type UserSummary struct {
	Login    string `json:"login"`
	HTMLURL  string `json:"html_url"`
	ReposURL string `json:"repos_url"`
}

type UserSearchResponse struct {
	TotalCount int           `json:"total_count"`
	Items      []UserSummary `json:"items"`
}

type RepoSearchResponse struct {
	TotalCount        int          `json:"total_count"`
	IncompleteResults bool         `json:"incomplete_results"`
	Items             []Repository `json:"items"`
}

func GetRepoUser(username string) []Repository {
	baseUrl := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	resp, err := http.Get(baseUrl)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic("Failed to fetch repos")
	}
	defer resp.Body.Close()
	var repos []Repository
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		panic(err)
	}
	return repos
}

func GetCommits(repo *Repository) []Commit {
	commitUrl := (*repo).CommitsURL
	resp, err := http.Get(commitUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var commits []Commit
	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		panic(err)
	}
	return commits
}

func GetUsers(userName string) UserSearchResponse {
	url := fmt.Sprintf("https://api.github.com/search/users?q=%s&sort=followers&order=desc&per_page=10&page=1", userName)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var users UserSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		panic(err)
	}
	return users
}

func GetRepos(query string) RepoSearchResponse {
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+sort:stars&per_page=10", query)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var repos RepoSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		panic(err)
	}
	return repos
}
