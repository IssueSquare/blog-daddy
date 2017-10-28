package main

import (
	"io"
	"regexp"

	"github.com/IssueSquare/blog-daddy/adapters/git"
	"github.com/IssueSquare/blog-daddy/providers/s3"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/russross/blackfriday.v2"
)

type MarkdownParser struct {
	reader io.Reader
}

func NewMarkdownParser(r io.Reader) *MarkdownParser {
	return &MarkdownParser{r}
}

func (m *MarkdownParser) Read(p []byte) (n int, err error) {
	n, err = m.reader.Read(p)
	copy(p, blackfriday.Run(p))

	return n, err
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	//viper.AddConfigPath("/var/run/secret")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}

	router := gin.Default()
	//router.LoadHTMLGlob("templates/*")

	router.POST("/setup", func(c *gin.Context) {
		var u git.User
		if c.BindJSON(&u) == nil {
			//minio cannot create bucket name with uppercase
			//and github is not case sensitive
			//trans to lowercase is not dangerous
			u.User = strings.ToLower(u.User)
			GithubHandler := git.NewGitHandler("https://api.github.com")
			mds := make([]git.GitRepoContent, 0)
			mds, err := GithubHandler.FetchRepoContents(u)
			if err != nil {
				panic(err)
			}

			//create user's bucket
			S3Provider := s3.NewS3Provider(viper.GetString("S3Endpoint"), viper.GetString("S3AccessKey"), viper.GetString("S3SecretKey"))

			err = S3Provider.CreateBucket(u.User)
			if err != nil {
				panic(err)
			}

			for _, md := range mds {

				resp, err := http.Get(md.Download_Url)

				if err != nil {
					panic(err)
				}

				defer resp.Body.Close()

				m := NewMarkdownParser(resp.Body)

				//TODO upload really index.html
				err = S3Provider.Upload(u.User, regexp.MustCompile("\\.[^.]+$").ReplaceAllString(md.Path, ".html"), m)
			}

			c.JSON(http.StatusOK, gin.H{"url": "https://s3.arthurma.com.tw/" + u.User + "/index.html"})
		}
	})

	//Webhook Endpint used to update/modify S3 doc contents
	router.POST("/webhook", func(c *gin.Context) {
		//payload from push event
		payload, err := c.GetRawData()
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"payload": payload})
	})

	router.Run(":8080")
}
