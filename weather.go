package finder

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

const (
	// he weather
	WeatherHe            = "heweather"
	WeatherKindHeWeather = "weather"
	WeatherKindHeAir     = "air"

	// sina weather
	WeatherSina = "sina"

	// config
	WeatherHttpTimeout = 60
)

/////////////////////////////
// Weather Finder Result
/////////////////////////////
type WeatherResult map[string]interface{}

func (wr *WeatherResult) Json() string {
	body, _ := json.Marshal(wr)
	return string(body)
}

/////////////////////////////
// Weather Finder
/////////////////////////////
type WeatherFinder struct {
	Provider string
	AppKey   string
	TTL      int
}

func (wf *WeatherFinder) Find(city string, kind string) (*WeatherResult, error) {
	switch wf.Provider {
	case WeatherHe:
		return HeWeatherFind(wf, city, kind)
	default:
		break
	}
	return nil, errors.New(fmt.Sprintf("not found %s provider", wf.Provider))
}

func (wf *WeatherFinder) Close() {
	// nothing
}

func (wf *WeatherFinder) Version() string {
	return VERSION
}

func NewWeatherFinder(provider string, appKey string, ttl int) *WeatherFinder {
	return &WeatherFinder{
		Provider: provider,
		AppKey:   appKey,
		TTL:      ttl,
	}
}
