package models

type Post struct {
	Id              string
	Title           string
	ContentHTML     string
	ContentMarkdown string
}

func NewPost(id, title, contentHTML, contentMarkdown string) *Post {
	return &Post{
		Id:              id,
		Title:           title,
		ContentHTML:     contentHTML,
		ContentMarkdown: contentMarkdown,
	}
}
