package utils

func GetMenuOptions(searchType int) []string {
	if searchType == SearchRepo {
		return []string{"Open Repo", "Clone Repo", "Copy Link", "Cancel"}
	}
	return []string{"View Profile", "Show Repositories", "Follow User", "Cancel"}
}
