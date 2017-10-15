package git

import (
	"encoding/json"
	"net/http"
)

type (
	User struct {
		User string `form:"user" json:"user" binding:"required"`
		Repo string `form:"repo" json:"repo" binding:"required"`
		Path string `form:"path" json:"path" binding:"required"`
	}

	Git struct {
		Endpoint string
	}

	GitRepoContent struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Sha          string `json:"sha"`
		Download_Url string `json:"download_url"`
	}
)

func NewGitHandler(e string) *Git {
	git := new(Git)
	git.Endpoint = e
	return git
}

func (g *Git) FetchRepoContents(u User) ([]GitRepoContent, error) {

	dst := g.Endpoint + "/repos" + "/" + u.User + "/" + u.Repo + "/contents" + u.Path
	resp, err := http.Get(dst)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	GitRepoContents := make([]GitRepoContent, 0)
	json.NewDecoder(resp.Body).Decode(&GitRepoContents)

	return GitRepoContents, nil
}
