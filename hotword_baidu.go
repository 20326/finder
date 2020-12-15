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

func BaiduHotWordFind(hwf *HotWordFinder, from string, suggestUrl string) ([]*HotWordResult, error) {
	var hwrs []*HotWordResult

	log.Debug().Str("provider", hwf.Provider).Str("url", BaiduHotWordsUrl).Msgf("BaiduHotWordFind")

	// fetch data
	response, body, err := gorequest.New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("Content-Type", "application/json").
		Timeout(HotWordHttpTimeout * time.Second).
		Get(BaiduHotWordsUrl).End()

	if nil != err || http.StatusOK != response.StatusCode {
		return nil, errors.New("fetch baidu hotspot failed")
	}

	result := gjson.ParseBytes([]byte(body))

	status := result.Get("success").Int()
	if status != 1 {
		return nil, errors.New("fetch baidu response status is" + strconv.FormatInt(status, 10))
	}

	// format json
	hotwords := result.Get("result.topwords")
	for _, hotword := range hotwords.Array() {
		keyword := hotword.Get("keyword").String()
		weight, _ := strconv.Atoi(hotword.Get("searches").String())
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
