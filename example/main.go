package main

import (
	"encoding/json"

	fd "github.com/20326/finder"
	"github.com/phuslu/log"
)

func main() {
	// new engine
	engine := fd.NewEngine()

	// register
	engine.Register(fd.KindPhone, fd.NewPhoneFinder("../phone.2004.dat"))
	engine.Register(fd.KindIP2Loc, fd.NewIP2LocFinder("../ip2region.db"))
	engine.RegisterByKind(fd.KindWeather, fd.WeatherHe, fd.NewWeatherFinder(fd.WeatherHe, "appKey", 300))
	engine.RegisterByKind(fd.KindHotword, fd.HotWordTouTiao, fd.NewHotWordFinder(fd.HotWordTouTiao, "appFrom", "", 300))
	engine.RegisterByKind(fd.KindSuggestion, fd.SuggestionBaidu, fd.NewSuggestionFinder(fd.SuggestionBaidu, "appFrom", "", 300))

	pf := engine.GetPhoneFinder()
	phoneResult, err := pf.Find("18600228899")
	log.Info().Msgf("Phone: %v err: %v", phoneResult, err)

	lf := engine.GetIP2LocFinder()
	locResult, err := lf.Find("175.25.21.39", "memory")
	log.Info().Msgf("Location: %v err: %v", locResult, err)

	wf := engine.GetWeatherFinder(fd.WeatherHe)
	weatherResult, err := wf.Find("beijing", fd.WeatherKindHeAir)
	weatherData, err := json.Marshal(weatherResult)
	log.Info().Msgf("Weather: %s err: %v", string(weatherData), err)

	hwf := engine.GetHowWordFinder(fd.HotWordTouTiao)
	hotwordResult, err := hwf.Find("", "")
	hotwordData, err := json.Marshal(hotwordResult)
	log.Info().Msgf("HowWord: %s err: %v", string(hotwordData), err)

	suf := engine.GetSuggestionFinder(fd.SuggestionBaidu)
	suggestResult, err := suf.Find("苹果", "", "")
	suggestData, err := json.Marshal(suggestResult)
	log.Info().Msgf("Suggestion: %s err: %v", string(suggestData), err)

	// unregister
	engine.UnRegister(fd.KindIP2Loc)

	// destroy
	engine.Destroy()
}
