package adapters

import (
	"encoding/json"
	"net/http"
)

type Git struct{}

type GitRepoContent struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Sha          string `json:"sha"`
	Download_Url string `json:"download_url"`
}

func (g GitHub) FetchRepoContents(p string) ([]GitRepoContent, error) {
	resp, err := http.Get("https://api.github.com/users/shavenking/repos")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	GitRepoContents := make([]GitRepoContent, 0)
	json.NewDecoder(resp.Body).Decode(&GitRepoContents)

	return GitRepoContents, nil
}
