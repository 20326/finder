package finder

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/phuslu/log"
	"github.com/tidwall/gjson"
)

func SoGouHotWordFind(hwf *HotWordFinder, from string, suggestUrl string) ([]*HotWordResult, error) {
	var hwrs []*HotWordResult

	log.Debug().Str("provider", hwf.Provider).Str("url", SoGouHotWordsUrl).Msgf("SoGouHotWordFind")

	// fetch data
	response, body, err := gorequest.New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("Content-Type", "application/json").
		Timeout(HotWordHttpTimeout * time.Second).
		Get(SoGouHotWordsUrl).End()

	if nil != err || http.StatusOK != response.StatusCode {
		return nil, errors.New("fetch sogou hotspot failed")
	}

	jsonp, cutErr := CutJsonp(`\[.*\]`, body)
	if nil != cutErr {
		return nil, errors.New("fetch jsonp failed")
	}

	hotwords := gjson.ParseBytes([]byte(jsonp))
	for _, hotword := range hotwords.Array() {
		keyword := hotword.Get("word").String()
		weight, _ := strconv.Atoi(hotword.Get("weight").String())
		suggestUrl := hwf.GetSuggestUrl(keyword, from, suggestUrl)

		hwr := &HotWordResult{
			Keyword:    keyword,
			Weight:     weight,
			SuggestUrl: suggestUrl,
		}
		hwrs = append(hwrs, hwr)
	}
	return hwrs, nil
}
