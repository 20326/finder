package finder

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
)

func SoGouHotWordFind(hwf *HotWordFinder) ([]*HotWordResult, error) {
	var hwrs []*HotWordResult

	log.Printf("SoGouHotWordFind url: %s", SoGouHotWordsUrl)

	// fetch data
	response, body, err := gorequest.New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("Content-Type", "application/json").
		Timeout(HotWordHttpTimeout * time.Second).
		Get(SoGouHotWordsUrl).End()

	if nil != err || http.StatusOK != response.StatusCode {
		return nil, errors.New("fetch sogou hotspot failed")
	}

	jsonp, cutErr := cutJsonp(body)
	if nil != cutErr {
		return nil, errors.New("fetch jsonp failed")
	}
	log.Printf("SoGouHotWordFind jsonp: %s", jsonp)

	hotwords := gjson.ParseBytes([]byte(jsonp))
	for _, hotword := range hotwords.Array() {
		keyword := hotword.Get("word").String()
		weight, _ := strconv.Atoi(hotword.Get("weight").String())
		suggestUrl := hwf.GetSuggestUrl(keyword)

		hwr := &HotWordResult{
			Keyword:    keyword,
			Weight:     weight,
			SuggestUrl: suggestUrl,
		}
		hwrs = append(hwrs, hwr)
	}
	return hwrs, nil
}

func cutJsonp(jsonp string) (string, error) {
	infoRegex := regexp.MustCompile(`\[.*\]`)
	slices := infoRegex.FindStringSubmatch(jsonp)
	if len(slices) < 1 {
		return "", errors.New("invalid jsonp")
	}
	return slices[0], nil
}
