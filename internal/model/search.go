package model

type SearchNode struct {
	Parent string `json:"parent"`
	Name   string `json:"name"`
	IsDir  bool   `json:"is_dir"`
	Size   int64  `json:"size"`
}

type SearchReq struct {
	Parent             string `json:"parent"`
	Keywords           string `json:"keywords"`
	Scope              int    `json:"scope"` // 0=all, 1=folder, 2=file
	Page               int    `json:"page"`
	PerPage            int    `json:"per_page"`
	Password           string `json:"password"`
	SeparateWordSearch bool   `json:"separate_word_search"`
}

func (r SearchReq) Validate() error {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PerPage <= 0 {
		r.PerPage = 100
	}
	return nil
}
