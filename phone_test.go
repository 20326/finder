package finder

import (
	"fmt"
	"testing"
)

var (
	pf *PhoneFinder
)

func init() {
	pf = NewPhoneFinder("./phone.db")
}

func BenchmarkBtreeSearch(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var i = 0
		for p.Next() {
			i++
			_, err := pf.Find(fmt.Sprintf("%s%d%s", "1897", i&10000, "45"))
			if err != nil {
				b.Fatal(err)
			}
		}

	})
}

func TestBtreeSearch1(t *testing.T) {
	_, err := pf.Find("13580198235123123213213")
	if err == nil {
		t.Fatal("错误的结果")
	}
	t.Log(err)
}

func TestBtreeSearch2(t *testing.T) {
	_, err := pf.Find("1300")
	if err == nil {
		t.Fatal("错误的结果")
	}
	t.Log(err)
}

func TestBtreeSearch3(t *testing.T) {
	pr, err := pf.Find("1703576")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pr)
}

func TestBtreeSearch4(t *testing.T) {
	_, err := pf.Find("10074872323")
	if err == nil {
		t.Fatal("错误的结果")
	}
	t.Log(err)
}

func TestBtreeSearch5(t *testing.T) {
	_, err := pf.Find("afsd32323")
	if err == nil {
		t.Fatal("错误的结果")
	}
	t.Log(err)
}
