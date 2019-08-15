package documents

type PostDocument struct {
	Id              string `bson:"_id,omitempty"`
	Title           string
	ContentHTML     string
	ContentMarkdown string
}
