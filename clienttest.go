// net project main.go
package main



import (
	"time"
	"github.com/Bulesxz/go/pake"
	"encoding/json"
	"fmt"
	"github.com/Bulesxz/go/net"
)


func main() {
	 
	fmt.Println("start")
	
	login:=pake.LoginReq{1,2,"sxz"}
	ctx :=pake.ContextInfo{}
	ctx.SetSess(nil)
	
	ctx.SetUserId(9999)
	mes :=&pake.Messages{ctx}
	
	
	id:=pake.PakeId(1)
	b,_:=json.Marshal(login)
	buf:=mes.Encode(id,b)
	fmt.Println(login,mes,buf)
	go net.GloablTimingWheel.Run()
	c :=net.NewClient("tcp","127.0.0.1:9000")
	err :=c.ConnetcTimeOut(5*time.Second)
	if err!=nil {
		fmt.Println("c.ConnetcTimeOut err|",err)
		return 
	}
	begin := time.Now().UnixNano()
	times:=200000
	for i:=0;i<times;i++{
		
		recvBuf,err:=c.SendTimeOut(2*time.Second,buf)
		if err!=nil {
			fmt.Println("c.Send err|",err)
			return 
		}
		if recvBuf!=nil{
			var rsp pake.LoginRsp
			p:=mes.Decode(recvBuf)
			json.Unmarshal(p.GetBody(),&rsp)
			
			fmt.Println("rsp:",rsp)
		}
		
	}
	end := time.Now().UnixNano()
	usetime := float64(end-begin) / float64(1000000)
	fmt.Println("qps ",float64(times)/(usetime/1000))
	fmt.Println("tps ",usetime/float64(times),"ms")
	c.Close()
	//time.Sleep(time.Second*10)
}
