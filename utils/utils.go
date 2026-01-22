package utils

import (
	"github.com/chirag-diwan/RemGit/githubapi"
)

func PerformSearch(query string, searchType int) SearchResult {
	if query == "" {
		return SearchResult{
			Users: githubapi.UserSearchResponse{},
			Repos: githubapi.RepoSearchResponse{},
		}
	}

	if searchType == SearchRepo {
		repos := githubapi.GetRepos(query)
		return SearchResult{
			Repos: repos,
			Users: githubapi.UserSearchResponse{},
		}

	} else {
		users := githubapi.GetUsers(query)
		return SearchResult{
			Repos: githubapi.RepoSearchResponse{},
			Users: users,
		}
	}
}

func GetMenuOptions(searchType int) []string {
	if searchType == SearchRepo {
		return []string{"Open Repo", "Clone Repo", "Copy Link", "Cancel"}
	}
	return []string{"View Profile", "Show Repositories", "Follow User", "Cancel"}
}
