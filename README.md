Finder (GEO, China Phone, Weather)
----------------------------
### 快速使用
```
package main

import (
	"log"

	fd "github.com/20326/finder"
)

func main() {
	// new engine
	engine := fd.NewEngine()

	// register
	engine.Register(fd.KindPhone, fd.NewPhoneFinder("../phone.db"))

	pf := engine.GetPhoneFinder()
	phoneResult, err := pf.Find("18600228899")
	log.Printf("Phone: %v err: %v", phoneResult, err)

	// unregister
	engine.UnRegister(fd.KindPhone)

	// destroy
	engine.Destroy()
}

```

### 性能测试

```
go test -v --bench="."

```

### 手机号归属地查询

最新最全的中国境内手机号归属地信息库, 基于GO语言实现，使用二分查找法。

 - 归属地信息库文件大小：4,040,893 字节
 - 归属地信息库最后更新：2020年04月
 - 手机号段记录条数：447893

        | 4 bytes |                     <- phone.dat 版本号（如：1701即17年1月份）
        ------------
        | 4 bytes |                     <-  第一个索引的偏移
        -----------------------
        |  offset - 8            |      <-  记录区
        -----------------------
        |  index                 |      <-  索引区
        -----------------------

1. 头部为8个字节，版本号为4个字节，第一个索引的偏移为4个字节；
2. 记录区 中每条记录的格式为"<省份>|<城市>|<邮编>|<长途区号>\0"。 每条记录以'\0'结束；
3. 索引区 中每条记录的格式为"<手机号前七位><记录区的偏移><卡类型>"，每个索引的长度为9个字节；

https://git.oschina.net/oss/phonedata/attach_files

### GEO数据库 ip2region.db

https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.db


## License

Copyright (c) 2020 (brian)

Licensed under the [MIT license](https://opensource.org/licenses/MIT) ([`LICENSE`](LICENSE)).