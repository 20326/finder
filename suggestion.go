package finder

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ymzuiku/hit"
)

const (
	// Suggestion  kind
	SuggestionBaidu = "baidu"

	// Url
	BaiduSuggestionUrl = "https://m.baidu.com/su?type=3&ie=utf-8&json=1&wd=${keyword}"

	// Suggest Url
	DefaultBaiduSuggestionUrlFmt = "http://m.baidu.com/from=${from}/s?word=${keyword}&bd_page_type=1"

	// Config
	SuggestionHttpTimeout = 60
)

/////////////////////////////
// Suggestion Result
/////////////////////////////
type SuggestionResult struct {
	Keyword    string
	Weight     int
	SuggestUrl string
	Query      string
}

func (wr *SuggestionResult) Json() string {
	body, _ := json.Marshal(wr)
	return string(body)
}

/////////////////////////////
// Suggestion Finder
/////////////////////////////
type SuggestionFinder struct {
	Provider   string
	From       string
	TTL        int
	SuggestUrl string
}

func (suf *SuggestionFinder) Find(keyword string, from string, suggestUrl string) ([]*SuggestionResult, error) {
	switch suf.Provider {
	case SuggestionBaidu:
		return BaiduSuggestionFind(suf, keyword, from, suggestUrl)
	default:
		break
	}
	return nil, errors.New(fmt.Sprintf("not found %s provider", suf.Provider))
}

func (suf *SuggestionFinder) GetSuggestUrl(keyword string, from string, suggestUrl string) string {
	from = hit.If(from, from, suf.From).(string)
	suggestUrl = hit.If(suggestUrl, suggestUrl, suf.SuggestUrl).(string)
	switch suf.Provider {
	case SuggestionBaidu:
		suggestUrl = DefaultBaiduSuggestionUrlFmt
		break
	default:
		break
	}
	// replace
	suggestUrl = strings.ReplaceAll(suggestUrl, "${from}", from)
	suggestUrl = strings.ReplaceAll(suggestUrl, "${keyword}", keyword)
	return suggestUrl
}

func (suf *SuggestionFinder) Close() {
	// nothing
}

func (suf *SuggestionFinder) Version() string {
	return VERSION
}

func NewSuggestionFinder(provider string, from string, suggestUrl string, ttl int) *SuggestionFinder {
	return &SuggestionFinder{
		Provider:   provider,
		From:       from,
		TTL:        ttl,
		SuggestUrl: suggestUrl,
	}
}
