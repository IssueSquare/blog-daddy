package git

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-collections/collections/stack"
)

type (
	User struct {
		User    string `form:"user" json:"user" binding:"required"`
		Repo    string `form:"repo" json:"repo" binding:"required"`
		DocPath string `form:"path" json:"path" binding:"required"`
	}

	Git struct {
		Endpoint string
	}

	GitRepoContent struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Sha          string `json:"sha"`
		Path         string `json:"path"`
		Download_Url string `json:"download_url"`
	}
)

func NewGitHandler(e string) *Git {
	git := new(Git)
	git.Endpoint = e
	return git
}

func (g *Git) githubGetClient(u User, path string) (*http.Response, error) {

	var dst string
	if path == "" {
		dst = g.Endpoint + "/repos" + "/" + u.User + "/" + u.Repo + "/contents/" + u.DocPath
	} else {
		dst = g.Endpoint + "/repos" + "/" + u.User + "/" + u.Repo + "/contents/" + path
	}

	resp, err := http.Get(dst)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *Git) FetchRepoContents(u User) ([]GitRepoContent, error) {

	//first Get
	resp, err := g.githubGetClient(u, "")
	if err != nil {
		panic(err)
	}

	buffer := make([]GitRepoContent, 0)
	json.NewDecoder(resp.Body).Decode(&buffer)
	gitRepoContents := make([]GitRepoContent, 0)

	//craete a stack for store dir
	s := stack.New()

	for _, v := range buffer {
		if v.Type == "file" {
			log.Printf("gitRepoContents %s\n", v)
			gitRepoContents = append(gitRepoContents, v)
		}

		if v.Type == "dir" {
			log.Printf("dir %s\n", v)
			s.Push(v.Path)
		}
	}

	extra := make([]GitRepoContent, 0)
	gitRepoContents = append(gitRepoContents, g.digDir(u, s, extra)...)

	return gitRepoContents, nil
}

func (g *Git) digDir(u User, s *stack.Stack, c []GitRepoContent) []GitRepoContent {
	if s.Len() == 0 {
		return c
	}

	resp, err := g.githubGetClient(u, s.Pop().(string))
	if err != nil {
		panic(err)
	}

	buffer := make([]GitRepoContent, 0)
	json.NewDecoder(resp.Body).Decode(&buffer)

	for _, v := range buffer {
		if v.Type == "file" {
			log.Printf("gitRepoContents %s\n", v)
			c = append(c, v)
		}

		if v.Type == "dir" {
			log.Printf("dir %s\n", v)
			s.Push(v.Path)
		}
	}

	return g.digDir(u, s, c)
}
