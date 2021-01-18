package finder

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ymzuiku/hit"
)

const (
	// HotWord  kind
	HotWordBaidu   = "baidu"
	HotWordTouTiao = "toutiao"
	HotWord360So   = "so"
	HotWordSogou   = "sogou"
	HotWordSm      = "sm"

	// Url
	BaiduHotWordsUrl   = "http://top.baidu.com/mobile_v2/buzz/hotspot"
	SoHotWordsUrl      = "https://m.so.com/home/data?types=Hotwords&fmt=json"
	SoGouHotWordsUrl   = "https://wap.sogou.com/data/hotwords/sogou_app.js"
	TouTiaoHotWordsUrl = "https://i.snssdk.com/hot-event/hot-board/?count=50&tab_name=stream&origin=hot_board"
	SmHotWordsUrl      = "http://api.m.sm.cn/rest?method=tools.hot&size=0"

	// Suggest Url
	DefaultBaiduSuggestUrlFmt   = "http://m.baidu.com/from=${from}/s?word=${keyword}&bd_page_type=1"
	DefaultSoGouSuggestUrlFmt   = "https://wap.sogou.com/web/searchList.jsp?keyword=${keyword}&e=${from}"
	Default360SoSuggestUrlFmt   = "https://m.so.com/s?q=${keyword}&src=${from}"
	DefaultTouTiaoSuggestUrlFmt = "https://so.toutiao.com/search?keyword=${keyword}&traffic_source=${from}&original_source=1&source=client"
	DefaultSmSuggestUrlFmt      = "http://m.yz2.sm.cn/s?q=${keyword}&by=hot&from=${from}"

	// Config
	HotWordHttpTimeout = 60
)

/////////////////////////////
// HotWord Result
/////////////////////////////
type HotWordResult struct {
	Keyword    string
	Weight     int
	SuggestUrl string
}

func (wr *HotWordResult) Json() string {
	body, _ := json.Marshal(wr)
	return string(body)
}

/////////////////////////////
// HotWord Finder
/////////////////////////////
type HotWordFinder struct {
	Provider   string
	From       string
	TTL        int
	SuggestUrl string
}

func (hwf *HotWordFinder) Find(from string, suggestUrl string) ([]*HotWordResult, error) {
	switch hwf.Provider {
	case HotWordBaidu:
		return BaiduHotWordFind(hwf, from, suggestUrl)
	case HotWord360So:
		return SoHotWordFind(hwf, from, suggestUrl)
	case HotWordSogou:
		return SoGouHotWordFind(hwf, from, suggestUrl)
	case HotWordTouTiao:
		return TouTiaoHotWordFind(hwf, from, suggestUrl)
	case HotWordSm:
		return SmHotWordFind(hwf, from, suggestUrl)
	default:
		break
	}
	return nil, errors.New(fmt.Sprintf("not found %s provider", hwf.Provider))
}

func (hwf *HotWordFinder) GetSuggestUrl(keyword string, from string, suggestUrl string) string {
	from = hit.If(from, from, hwf.From).(string)
	suggestUrl = hit.If(suggestUrl, suggestUrl, hwf.SuggestUrl).(string)
	switch hwf.Provider {
	case HotWordBaidu:
		suggestUrl = DefaultBaiduSuggestUrlFmt
		break
	case HotWord360So:
		suggestUrl = Default360SoSuggestUrlFmt
		break
	case HotWordSogou:
		suggestUrl = DefaultSoGouSuggestUrlFmt
		break
	case HotWordTouTiao:
		suggestUrl = DefaultTouTiaoSuggestUrlFmt
		break
	case HotWordSm:
		suggestUrl = DefaultSmSuggestUrlFmt
		break
	default:
		break
	}
	// replace
	suggestUrl = strings.ReplaceAll(suggestUrl, "${from}", from)
	suggestUrl = strings.ReplaceAll(suggestUrl, "${keyword}", keyword)
	return suggestUrl
}

func (hwf *HotWordFinder) Close() {
	// nothing
}

func (hwf *HotWordFinder) Version() string {
	return VERSION
}

func NewHotWordFinder(provider string, from string, suggestUrl string, ttl int) *HotWordFinder {
	return &HotWordFinder{
		Provider:   provider,
		From:       from,
		TTL:        ttl,
		SuggestUrl: suggestUrl,
	}
}
