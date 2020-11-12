package engine

import (
	"log"
	"stone/go/crawler/fetcher"
)

//Worker：将fetcher与parser组合在一起
func Worker(r Request) (ParserResult, error) {
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Error fetching url %s : %v", r.Url, err)
		return ParserResult{}, err
	}
	return r.Parser.Parse(body, r.Url), nil
}
