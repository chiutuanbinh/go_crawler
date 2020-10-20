package xtype

//Article represent an article from news sites
type Article struct {
	ID        string         `json:"id,omitempty"`
	Title     string         `json:"title,omitempty"`
	Publisher string         `json:"publisher,omitempty"`
	URI       string         `json:"uri,omitempty"`
	Meta      ArticleMeta    `json:"meta,omitempty"`
	Content   ArticleContent `json:"content,omitempty"`
	Nlp       ArticleNlp     `json:"nlp_infos,omitempty"`
}

//ArticleContent contain the body of the article, include image, video and text
type ArticleContent struct {
	Parts []interface{}
}

type Paragraph struct {
	Content string `json:"content,omitempty"`
}

type Image struct {
	ID      string `json:"id,omitempty"`
	Caption string `json:"caption,omitempty"`
	URI     string `json:"uri,omitempty"`
}

type ArticleMeta struct {
	Text      string `json:"text,omitempty"`
	SourceID  string `json:"sourceid,omitempty"`
	PublishTs uint64 `json:"publishts,omitempty"`
}

type ArticleNlp struct {
	NamedEntities []string `json:"named_entities,omitempty"`
	KeyWords      []string `json:"keywords,omitempty"`
}
