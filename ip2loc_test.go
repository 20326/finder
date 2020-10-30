package finder

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var (
	ip2lf *IP2LocFinder
)

func init() {
	ip2lf = NewIP2LocFinder("./ip2region.db")
}

func genIpaddr() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func BenchmarkMemorySearch(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var i = 0
		for p.Next() {
			i++
			ip := genIpaddr()
			info, err := ip2lf.Find(ip, "memory")
			if err != nil {
				b.Fatal(ip, info, err)
			}
		}
	})
}
