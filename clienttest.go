// net project main.go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/Bulesxz/go/base"
	log "github.com/Bulesxz/go/logger"
	"github.com/Bulesxz/go/net"
	"github.com/Bulesxz/go/pake"
	"sync"
	"sync/atomic"
	"time"
)

var seq uint64

func NewPake() (*pake.Messages, []byte, pake.LoginReq, uint64) {
	id := atomic.AddUint64(&seq, 1)
	login := pake.LoginReq{int32(id), 2, "sxz"}
	ctx := pake.ContextInfo{}
	ctx.SetSess("session")
	ctx.SetId(pake.LoginId)
	ctx.SetUserId("125222")
	ctx.SetSeq(id)
	mes := &pake.Messages{ctx}
	b, _ := json.Marshal(login)
	buf := mes.Encode(b)
	return mes, buf, login, id
}

func testBenchmark() {
	nclient := 20
	connChan := make(chan *net.Client, nclient)
	for i := 0; i < nclient; i++ {
		c := net.NewClient("tcp", "127.0.0.1:9000")
		err := c.ConnetcTimeOut(1 * time.Second)
		if err != nil {
			//fmt.Println("c.ConnetcTimeOut err|", err)
			log.Error("c.ConnetcTimeOut err|", err)
			return
		}
		connChan <- c
		defer c.Close()
	}
	//fmt.Println("newclient end.............buf",buf,"len",len(buf))
	//time.Sleep(10*time.Minute)
	f := func() bool {
		var c *net.Client
		select {
		case c = <-connChan:
		default:
			fmt.Println("no client")
			return false
		}
		mes, buf, _, id := NewPake()

		recvBuf, err := c.SendTimeOut(3*time.Second, buf)
		if err != nil {
			//fmt.Println("c.Send err|", err)
			log.Error("c.Send err|", err)
			return false
		}

		//fmt.Println("login:",login,"id",id)
		if recvBuf != nil {
			var rsp pake.LoginRsp
			p := mes.Decode(recvBuf)
			err = json.Unmarshal(p.GetBody(), &rsp)
			if err != nil {
				fmt.Println(err)
				return false
			}

			//fmt.Println("login:",login,"rsp:",rsp,"session seq:",p.GetSession(),"id",id)
			if p.GetSession().Seq != id || rsp.A != int32(id) {
				fmt.Println("rsp:", rsp)
				return false
			}
		} else {
			fmt.Println("f() false")
			return false
		}
		connChan <- c
		return true
	}

	var n int32 = 100000
	usetime, failedNum := base.BenchmarkFunc(n, nclient, f)
	fmt.Println("proxy usetime=", usetime, " failedNum=", failedNum, " qps=", float64(n)/(usetime/1000), " tps=", (usetime/float64(n))/1000)
}

func testerrpack() {

	//同一个连接并发发包
	c := net.NewClient("tcp", "127.0.0.1:9000")
	err := c.ConnetcTimeOut(1 * time.Second)
	if err != nil {
		//fmt.Println("c.ConnetcTimeOut err|", err)
		log.Error("c.ConnetcTimeOut err|", err)
		return
	}
	defer c.Close()

	fmt.Println("************同一个连接并发发包**************")
	var failNum int32 = 0
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			fail := errpack(&wg, c)
			if fail == false {
				atomic.AddInt32(&failNum, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("failNum ", failNum)
}

func errpack(wg *sync.WaitGroup, c *net.Client) bool {
	defer wg.Done()
	mes, buf, login, id := NewPake()
	recvBuf, err := c.SendTimeOut(3*time.Second, buf)
	if err != nil {
		fmt.Println("c.Send err|", err)
		log.Error("c.Send err|", err)
		return false
	}
	if recvBuf != nil {
		var rsp pake.LoginRsp
		p := mes.Decode(recvBuf)
		err = json.Unmarshal(p.GetBody(), &rsp)
		if err != nil {
			fmt.Println(err)
			return false
		}
		if p.GetSession().Seq != id {
			fmt.Println("login:", login, "rsp:", rsp, "id:", id, "rsp seq:", p.GetSession().Seq)
			//fmt.Println("login:",login,"|rsp:",p,id)
			return false
		}

	} else {
		fmt.Println("recvBuf == nil")
		return false
	}
	return true
}

func testRcClient() {
	rcclient := net.NewRcClient("tcp", "127.0.0.1:9000")
	err := rcclient.ConnetcTimeOut(time.Second * 5)
	if err != nil {
		fmt.Println("err", err)
		log.Error(err)
		return
	}

	for i := 1; i <10000; i++ {
		go func(j int32) {
			//fmt.Println(j)
			mes := &pake.Messages{}
			login := pake.LoginReq{j, j, "1"}
			id := net.InitPake(mes, pake.LoginId)

			//
			/*msg, err := json.Marshal(login)
			if err != nil {
				log.Error(err)
				return
			}*/
			
			//sendData := mes.Encode(msg)
			//fmt.Println("mes", mes.Context)
			//fmt.Println("id", id, "pake", sendData)

			rsp,err:=rcclient.Call(mes,login,time.Millisecond*1000)
			//fmt.Println("rsp:",rsp,"err",err)
			logrsp :=pake.LoginRsp{}
			//fmt.Println("req:",msg,"body:",rsp.(*pake.Pake).GetBody())
			if err!=nil {
				fmt.Println("rsp:",rsp,"err:",err)
				return
			}
			//fmt.Println("rrrrrrrsp:",rsp,"err:",err)
			json.Unmarshal(rsp.(*pake.Pake).GetBody(),&logrsp)
			if login.A != logrsp.A {
				fmt.Println("---------------err--------------")
			}
			fmt.Println(id,",",login,",",logrsp,err)
		}(int32(i))
	}

}
func main() {

	fmt.Println("start")

	go net.GloablTimingWheel.Run()

	//testBenchmark()
	//testerrpack()
	testRcClient()
	time.Sleep(time.Second * 10)
}
