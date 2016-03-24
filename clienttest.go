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
	recvBuf,err:=c.SendTimeOut(10*time.Second,buf)
	if err!=nil {
		fmt.Println("c.Send err|",err)
		return 
	}
	//fmt.Println("recvBuf:",recvBuf)
	
	/*err=c.Send(buf)
	if err!=nil {
		fmt.Println("c.Send err|",err)
		return 
	}*/
	/*var receiveBuf []byte
	err =c.Receive(&receiveBuf)
	if err!=nil {
		fmt.Println("c.Receive err|",err)
		return 
	}
	fmt.Println("receiveBuf:",receiveBuf)
	*/
	var rsp pake.LoginRsp
	p:=mes.Decode(recvBuf)
	json.Unmarshal(p.GetBody(),&rsp)
	
	fmt.Println("rsp:",rsp)
	c.Close()
	//time.Sleep(time.Second*10)
}
