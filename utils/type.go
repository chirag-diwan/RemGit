package utils

import (
	"github.com/chirag-diwan/RemGit/githubapi"
)

type SearchResult struct {
	Users githubapi.UserSearchResponse
	Repos githubapi.RepoSearchResponse
}

type Menu struct {
	Active  bool
	Options []string
	Cursor  int
}
