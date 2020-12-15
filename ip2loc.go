package finder

import (
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/phuslu/log"
)

/////////////////////////////
// IP2Location Finder Result
/////////////////////////////
type IP2LocResult struct {
	IP       string `json:"ip"`
	CityId   int64  `json:"city_id"`
	Country  string `json:"country"`
	Region   string `json:"region"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
}

//func (ip2lr IP2LocResult) Json() string {
//	body, _ := json.Marshal(ip2lr)
//	return string(body)
//}

/////////////////////////////
// IP2Location Finder
/////////////////////////////
type IP2LocFinder struct {
	engine *ip2region.Ip2Region
}

func NewIP2LocFinder(db string) *IP2LocFinder {
	engine, err := ip2region.New(db)
	if err != nil {
		log.Fatal().Msgf("init ip2region db:%s, err: %s", db, err)
	}
	return &IP2LocFinder{engine}
}

func (iplf *IP2LocFinder) Find(ip string, algorithm string) (*IP2LocResult, error) {
	var ipInfo ip2region.IpInfo
	var err error

	switch algorithm {
	case "binary":
		ipInfo, err = iplf.engine.BinarySearch(ip)
		break
	case "btree":
		ipInfo, err = iplf.engine.BtreeSearch(ip)
		break
	case "memory":
		ipInfo, err = iplf.engine.MemorySearch(ip)
		break
	default:
		ipInfo, err = iplf.engine.MemorySearch(ip)
	}
	if err == nil {
		return &IP2LocResult{
			IP:       ip,
			CityId:   ipInfo.CityId,
			Country:  ipInfo.Country,
			Region:   ipInfo.Region,
			Province: ipInfo.Province,
			City:     ipInfo.City,
			ISP:      ipInfo.ISP,
		}, err
	} else {
		return nil, err
	}
}

func (iplf *IP2LocFinder) Close() {
	iplf.engine.Close()
}

func (iplf *IP2LocFinder) Version() string {
	return "ip2region@v2.2.0"
}
