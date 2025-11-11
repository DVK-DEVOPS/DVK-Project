package db

type PageRepoInterface interface {
	FindSearchResults(searchStr, language string) ([]Result, error)
}
