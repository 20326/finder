package finder

import (
	"sync"
)

const (
	KindIP2Loc     = "ip2location"
	KindPhone      = "phone"
	KindWeather    = "weather"
	KindHotword    = "hotword"
	KindSuggestion = "suggestion"

	VERSION = "v1.0.12"
)

var (
	once     sync.Once
	instance Engine
)

/////////////////////////////
// finder interface
/////////////////////////////
type Finder interface {
	Close()
}

/////////////////////////////
// finder engine
/////////////////////////////
// Engine is default finder implement.
type Engine struct {
	finderMap map[string]Finder
}

// NewEngine creates a new Engine
func NewEngine() Engine {
	once.Do(func() {
		instance = Engine{
			finderMap: map[string]Finder{},
		}
	})
	return instance
}

// Register set @finder with @key to local memory.
func (engine *Engine) Register(key string, finder Finder) {
	engine.finderMap[key] = finder
}

// Register set @finder with @kind @key to local memory.
func (engine *Engine) RegisterByKind(kind, key string, finder Finder) {
	engine.finderMap[kind+"_"+key] = finder
}

// Get with @key return @finder.
func (engine *Engine) Get(key string) Finder {
	return engine.finderMap[key]
}

// UnRegister gets finder map.
func (engine *Engine) UnRegister(key string) {
	// destroy finder
	finder := engine.finderMap[key]
	if finder != nil {
		finder.(Finder).Close()
	} else {
		delete(engine.finderMap, key)
	}
}

// Destroy will destroy all finder, so it only is called once.
func (engine *Engine) Destroy() {
	// destroy finder
	for key := range engine.finderMap {
		finder := engine.finderMap[key]
		if finder != nil {
			finder.(Finder).Close()
		} else {
			delete(engine.finderMap, key)
		}
	}
}

// GetPhoneFinder with @key return @finder.
func (engine *Engine) GetPhoneFinder() *PhoneFinder {
	fd := engine.finderMap[KindPhone]
	if fd == nil {
		panic("not found finder")
	}
	return fd.(*PhoneFinder)
}

// GetIP2LocFinder with @key return @finder.
func (engine *Engine) GetIP2LocFinder() *IP2LocFinder {
	fd := engine.finderMap[KindIP2Loc]
	if fd == nil {
		panic("not found finder")
	}
	return fd.(*IP2LocFinder)
}

// GetWeatherFinder with @key return @finder.
func (engine *Engine) GetWeatherFinder(key string) *WeatherFinder {
	fd := engine.finderMap[KindWeather+"_"+key]
	if fd == nil {
		panic("not found finder")
	}
	return fd.(*WeatherFinder)
}

// GetHowWordFinder with @key return @finder.
func (engine *Engine) GetHowWordFinder(key string) *HotWordFinder {
	fd := engine.finderMap[KindHotword+"_"+key]
	if fd == nil {
		panic("not found finder")
	}
	return fd.(*HotWordFinder)
}

// GetSuggestionFinder with @key return @finder.
func (engine *Engine) GetSuggestionFinder(key string) *SuggestionFinder {
	fd := engine.finderMap[KindSuggestion+"_"+key]
	if fd == nil {
		panic("not found finder")
	}
	return fd.(*SuggestionFinder)
}
