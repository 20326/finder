package finder

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
)

func SoHotWordFind(hwf *HotWordFinder, from string, suggestUrl string) ([]*HotWordResult, error) {
	var hwrs []*HotWordResult

	log.Printf("SoHotWordFind url: %s", SoHotWordsUrl)

	// fetch data
	response, body, err := gorequest.New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("Content-Type", "application/json").
		Timeout(HotWordHttpTimeout * time.Second).
		Get(SoHotWordsUrl).End()

	if nil != err || http.StatusOK != response.StatusCode {
		return nil, errors.New("fetch so hotspot failed")
	}

	hotwords := gjson.ParseBytes([]byte(body)).Get("Hotwords")
	for _, hotword := range hotwords.Array() {
		keyword := hotword.Get("title").String()
		weight, _ := strconv.Atoi(hotword.Get("n").String())
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
