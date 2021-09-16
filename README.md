# timegroup
这是一个带有超时的并发包。
使用示例如下：

```go
package main

import (
	"log"
	"time"

	"github.com/ekkoxx/timegroup"
)

func main() {
	var g = timegroup.New()
	g.Go(func() error {
		//code 1
		return nil
	})
	g.Go(func() error {
		//code 2
		return nil
	})
	g.Go(func() error {
		//code 3
		return nil
	})
	//timeout and return immediately
	err := g.WaitTimeout(time.Second)
	if err != nil {
		log.Println(err)
	}
}

```