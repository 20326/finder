package finder

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/phuslu/log"
	"github.com/tidwall/gjson"
)

func HeWeatherFind(wf *WeatherFinder, city string, kind string) (*WeatherResult, error) {
	var wr *WeatherResult

	targetUrl := fmt.Sprintf("https://free-api.heweather.net/s6/%s/now?key=%s&location=%s", kind, wf.AppKey, city)
	log.Debug().Str("city", city).Str("kind", kind).Str("provider", wf.Provider).Str("url", targetUrl).Msgf("HeWeatherFind")

	// fetch data
	response, body, err := gorequest.New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("Content-Type", "application/json").
		Timeout(WeatherHttpTimeout * time.Second).
		Get(targetUrl).EndStruct(&wr)

	if nil != err || http.StatusOK != response.StatusCode {
		return nil, errors.New("fetch weather now failed")
	}
	status := gjson.ParseBytes(body).Get("HeWeather6.0.status")
	if status.String() != "ok" {
		return nil, errors.New(status.String())
	}

	return wr, nil
}
