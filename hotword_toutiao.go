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

func TouTiaoHotWordFind(hwf *HotWordFinder, from string, suggestUrl string) ([]*HotWordResult, error) {
	var hwrs []*HotWordResult

	log.Printf("TouTiaoHotWordFind url: %s", TouTiaoHotWordsUrl)

	// fetch data
	response, body, err := gorequest.New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("Content-Type", "application/json").
		Timeout(HotWordHttpTimeout * time.Second).
		Get(TouTiaoHotWordsUrl).End()

	if nil != err || http.StatusOK != response.StatusCode {
		return nil, errors.New("fetch toutiao hotspot failed")
	}

	// format json
	hotwords := gjson.ParseBytes([]byte(body)).Get("data")
	for _, hotword := range hotwords.Array() {
		keyword := hotword.Get("Title").String()
		weight, _ := strconv.Atoi(hotword.Get("HotValue").String())
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
