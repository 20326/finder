package finder

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	CMCC             byte = iota + 0x01 //中国移动
	CUCC                                //中国联通
	CTCC                                //中国电信
	CTCC_V                              //电信虚拟运营商
	CUCC_V                              //联通虚拟运营商
	CMCC_V                              //移动虚拟运营商
	PhoneIntLen      = 4
	PhoneCharLen     = 1
	PhoneHeadLength  = 8
	PhoneIndexLength = 9
)

var (
	CarrierMap = map[byte]string{
		CMCC:   "中国移动",
		CUCC:   "中国联通",
		CTCC:   "中国电信",
		CTCC_V: "中国电信虚拟运营商",
		CUCC_V: "中国联通虚拟运营商",
		CMCC_V: "中国移动虚拟运营商",
	}
)

/////////////////////////////
// Phone Finder Result
/////////////////////////////
type PhoneResult struct {
	Phone    string `json:"phone"`
	Province string `json:"province"`
	City     string `json:"city"`
	ZipCode  string `json:"zipcode"`
	Region   string `json:"region"`
	Carrier  string `json:"carrier"`
}

/////////////////////////////
// Phone Finder
/////////////////////////////

type PhoneFinder struct {
	content     []byte
	totalLen    int32
	firstOffset int32
}

func NewPhoneFinder(db string) *PhoneFinder {
	var err error
	content, err := ioutil.ReadFile(db)
	if err != nil {
		panic(err)
	}
	totalLen := int32(len(content))
	firstOffset := get4(content[PhoneIntLen : PhoneIntLen*2])

	return &PhoneFinder{
		content:     content,
		totalLen:    totalLen,
		firstOffset: firstOffset,
	}
}

func (pf *PhoneFinder) Find(key string) (*PhoneResult, error) {
	return pf.BtreeSearch(key)
}

func (pf *PhoneFinder) Close() {
	// nothing
}

func (pf *PhoneFinder) Version() string {
	return fmt.Sprintf("Version: %s\nTotalRecord: %d\nFirstRecordOffset: %d\n", pf.version(), pf.totalRecord(), pf.firstRecordOffset())
}

/////////////////////////////
// 二分法查询phone数据
/////////////////////////////
func (pf *PhoneFinder) BtreeSearch(phoneNum string) (*PhoneResult, error) {
	if len(phoneNum) < 7 || len(phoneNum) > 11 {
		return nil, errors.New("illegal phone length")
	}

	var left int32
	phoneSevenInt, err := getN(phoneNum[0:7])
	if err != nil {
		return nil, errors.New("illegal phone number")
	}
	phoneSevenInt32 := int32(phoneSevenInt)
	right := (pf.totalLen - pf.firstOffset) / PhoneIndexLength
	for {
		if left > right {
			break
		}
		mid := (left + right) / 2
		offset := pf.firstOffset + mid*PhoneIndexLength
		if offset >= pf.totalLen {
			break
		}
		curPhone := get4(pf.content[offset : offset+PhoneIntLen])
		recordOffset := get4(pf.content[offset+PhoneIntLen : offset+PhoneIntLen*2])
		carrier := pf.content[offset+PhoneIntLen*2 : offset+PhoneIntLen*2+PhoneCharLen][0]
		switch {
		case curPhone > phoneSevenInt32:
			right = mid - 1
		case curPhone < phoneSevenInt32:
			left = mid + 1
		default:
			cByte := pf.content[recordOffset:]
			endOffset := int32(bytes.Index(cByte, []byte("\000")))
			data := bytes.Split(cByte[:endOffset], []byte("|"))
			carrier, ok := CarrierMap[carrier]
			if !ok {
				carrier = "未知电信运营商"
			}
			return &PhoneResult{
				Phone:    phoneNum,
				Province: string(data[0]),
				City:     string(data[1]),
				ZipCode:  string(data[2]),
				Region:   string(data[3]),
				Carrier:  carrier,
			}, nil
		}
	}
	return nil, errors.New("phone's data not found")
}

func (pf *PhoneFinder) totalRecord() int32 {
	return (int32(len(pf.content)) - pf.firstRecordOffset()) / PhoneIndexLength
}

func (pf *PhoneFinder) firstRecordOffset() int32 {
	return get4(pf.content[PhoneIntLen : PhoneIntLen*2])
}

func (pf *PhoneFinder) version() string {
	return string(pf.content[0:PhoneIntLen])
}

func get4(b []byte) int32 {
	if len(b) < 4 {
		return 0
	}
	return int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
}

func getN(s string) (uint32, error) {
	var n, cutoff, maxVal uint32
	i := 0
	base := 10
	cutoff = (1<<32-1)/10 + 1
	maxVal = 1<<uint(32) - 1
	for ; i < len(s); i++ {
		var v byte
		d := s[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		case 'a' <= d && d <= 'z':
			v = d - 'a' + 10
		case 'A' <= d && d <= 'Z':
			v = d - 'A' + 10
		default:
			return 0, errors.New("invalid syntax")
		}
		if v >= byte(base) {
			return 0, errors.New("invalid syntax")
		}

		if n >= cutoff {
			// n*base overflows
			n = 1<<32 - 1
			return n, errors.New("value out of range")
		}
		n *= uint32(base)

		n1 := n + uint32(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = 1<<32 - 1
			return n, errors.New("value out of range")
		}
		n = n1
	}
	return n, nil
}
