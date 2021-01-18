package finder

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/phuslu/log"
	"github.com/tidwall/gjson"
)

func SmHotWordFind(hwf *HotWordFinder, from string, suggestUrl string) ([]*HotWordResult, error) {
	var hwrs []*HotWordResult

	log.Debug().Str("provider", hwf.Provider).Str("url", SmHotWordsUrl).Msgf("SmHotWordFind")

	// fetch data
	response, body, err := gorequest.New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("Content-Type", "application/json").
		Timeout(HotWordHttpTimeout * time.Second).
		Get(SmHotWordsUrl).End()

	if nil != err || http.StatusOK != response.StatusCode {
		return nil, errors.New("fetch sm hotspot failed")
	}

	hotwords := gjson.ParseBytes([]byte(body))
	// format json
	lists := map[string]int64{}
	for _, hotword := range hotwords.Array() {
		keyword := hotword.Get("title").String()
		weight := hotword.Get("hot_flag").Int()
		suggestUrl := hwf.GetSuggestUrl(keyword, from, suggestUrl)

		if lists[keyword] == 0 {
			lists[keyword] = weight
			hwr := &HotWordResult{
				Keyword:    keyword,
				Weight:     int(weight),
				SuggestUrl: suggestUrl,
			}
			hwrs = append(hwrs, hwr)
		}
	}
	return hwrs, nil
}
