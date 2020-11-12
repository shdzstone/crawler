package model

//前端model
type SearchResult struct {
	Hits     int64
	Start    int
	Items    []interface{}
	Query    string
	PreFrom  int
	NextFrom int
}
