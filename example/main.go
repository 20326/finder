package main

import (
	"encoding/json"
	"log"

	fd "github.com/20326/finder"
)

func main() {
	// new engine
	engine := fd.NewEngine()

	// register
	engine.Register(fd.KindPhone, fd.NewPhoneFinder("../phone.2004.dat"))
	engine.Register(fd.KindIP2Loc, fd.NewIP2LocFinder("../ip2region.db"))
	engine.Register(fd.KindWeather, fd.NewWeatherFinder(fd.WeatherHe, "appKey", 300))
	engine.Register(fd.KindHotword, fd.NewHotWordFinder(fd.HotWordTouTiao, "appFrom", "", 300))
	engine.Register(fd.KindSuggestion, fd.NewSuggestionFinder(fd.SuggestionBaidu, "appFrom", "", 300))

	pf := engine.GetPhoneFinder()
	phoneResult, err := pf.Find("18600228899")
	log.Printf("Phone: %v err: %v", phoneResult, err)

	lf := engine.GetIP2LocFinder()
	locResult, err := lf.Find("175.25.21.39", "memory")
	log.Printf("Location: %v err: %v", locResult, err)

	wf := engine.GetWeatherFinder()
	weatherResult, err := wf.Find("beijing", fd.WeatherKindHeAir)
	weatherData, err := json.Marshal(weatherResult)
	log.Printf("Weather: %v err: %v", string(weatherData), err)

	hwf := engine.GetHowWordFinder()
	hotwordResult, err := hwf.Find("", "")
	hotwordData, err := json.Marshal(hotwordResult)
	log.Printf("HowWord: %v err: %v", string(hotwordData), err)

	suf := engine.GetSuggestionFinder()
	suggestResult, err := suf.Find("苹果", "", "")
	suggestData, err := json.Marshal(suggestResult)
	log.Printf("Suggestion: %v err: %v", string(suggestData), err)

	// unregister
	engine.UnRegister(fd.KindIP2Loc)

	// destroy
	engine.Destroy()
}
