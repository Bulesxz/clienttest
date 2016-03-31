// net project main.go
package main

import (
	"benchmark"
	"encoding/json"
	"fmt"
	"github.com/Bulesxz/go/net"
	"github.com/Bulesxz/go/pake"
	"time"
)

func main() {

	fmt.Println("start")

	login := pake.LoginReq{1, 2, "sxz"}
	ctx := pake.ContextInfo{}
	ctx.SetSess(nil)

	ctx.SetUserId(9999)
	mes := &pake.Messages{ctx}

	id := pake.PakeId(1)
	b, _ := json.Marshal(login)
	buf := mes.Encode(id, b)
	fmt.Println(login, mes, buf)
	go net.GloablTimingWheel.Run()

	connChan := make(chan *net.Client, 1000)
	for i := 0; i < 1000; i++ {
		c := net.NewClient("tcp", "127.0.0.1:9000")
		err := c.ConnetcTimeOut(5 * time.Second)
		if err != nil {
			fmt.Println("c.ConnetcTimeOut err|", err)
			return
		}
		connChan <- c
		defer c.Close()
	}
	f := func() bool {
		var c *net.Client
		select {
			case c = <-connChan:
			default: return false
		}
		recvBuf, err := c.SendTimeOut(2*time.Second, buf)
		if err != nil {
			fmt.Println("c.Send err|", err)
			return false
		}
		if recvBuf != nil {
			var rsp pake.LoginRsp
			p := mes.Decode(recvBuf)
			err = json.Unmarshal(p.GetBody(), &rsp)
			if err != nil {
				return false
			}
			//fmt.Println("rsp:", rsp)
		} else {
			return false
		}
		connChan <- c
		return true
	}
	n := 10000
	usetime, failedNum := benchmark.BenchmarkFunc(n, 100, f)
	fmt.Println("proxy usetime=", usetime, " failedNum=", failedNum, " qps=", float64(n)/(usetime/1000), " tps=", (usetime/float64(n))/1000)

	
	//time.Sleep(time.Second*10)
}
