// net project main.go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/Bulesxz/go/base"
	"github.com/Bulesxz/go/net"
	"github.com/Bulesxz/go/pake"
	log "github.com/Bulesxz/go/logger"
	"time"
)

func main() {

	fmt.Println("start")
	//log.Init()
	log.Error("frhtynh")
	
	login := pake.LoginReq{1, 2, "sxz"}
	ctx := pake.ContextInfo{}
	ctx.SetSess("session")

	ctx.SetUserId("125222")
	mes := &pake.Messages{ctx}

	ctx.Id = pake.LoginId

	b, _ := json.Marshal(login)
	buf := mes.Encode(b)
	//log.Info(login, mes, buf)
	go net.GloablTimingWheel.Run()

	nclient:=10
	connChan := make(chan *net.Client, nclient)
	for i := 0; i < nclient; i++ {
		fmt.Println("NewClient")
		c := net.NewClient("tcp", "127.0.0.1:9000")
		err := c.ConnetcTimeOut(5 * time.Second)
		if err != nil {
			fmt.Println("c.ConnetcTimeOut err|", err)
			log.Error("c.ConnetcTimeOut err|", err)
			return
		}
		connChan <- c
		defer c.Close()
	}
	f := func() bool {
		var c *net.Client
		select {
		case c = <-connChan:
		default:
			return false
		}
		recvBuf, err := c.SendTimeOut(2*time.Second, buf)
		if err != nil {
			log.Error("c.Send err|", err)
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
	var n int32 = 100
	usetime, failedNum := base.BenchmarkFunc(n, 100, f)
	fmt.Println("proxy usetime=", usetime, " failedNum=", failedNum, " qps=", float64(n)/(usetime/1000), " tps=", (usetime/float64(n))/1000)

	//time.Sleep(time.Second*10)
}
