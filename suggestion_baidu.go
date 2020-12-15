package finder

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/phuslu/log"
	"github.com/tidwall/gjson"
)

func BaiduSuggestionFind(suf *SuggestionFinder, keyword string, from string, suggestUrl string) ([]*SuggestionResult, error) {
	var surt []*SuggestionResult

	baiduSuggestionUrl := strings.ReplaceAll(BaiduSuggestionUrl, "${keyword}", keyword)
	log.Debug().Str("provider", suf.Provider).Str("url", baiduSuggestionUrl).Msgf("BaiduSuggestionFind")

	// fetch data
	response, body, err := gorequest.New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("Content-Type", "application/json").
		Timeout(SuggestionHttpTimeout * time.Second).
		Get(baiduSuggestionUrl).End()

	if nil != err || http.StatusOK != response.StatusCode {
		return nil, errors.New("fetch baidu suggestion failed")
	}

	jsonp, cutErr := CutJsonp(`\{.*\}`, body)
	if nil != cutErr {
		return nil, errors.New("fetch jsonp failed")
	}

	// format json
	suggestions := gjson.ParseBytes([]byte(jsonp))

	for _, suggestion := range suggestions.Get("s").Array() {
		suggestUrl := suf.GetSuggestUrl(suggestion.String(), from, suggestUrl)

		sur := &SuggestionResult{
			Keyword:    suggestion.String(),
			Weight:     0,
			Query:      keyword,
			SuggestUrl: suggestUrl,
		}
		surt = append(surt, sur)
	}
	return surt, nil
}
