package domains

type Article struct {
	author  string
	title   string
	content string
}

type ArticleProvider interface {
	//publish to s3
	Publish(a Article) error
}
