package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"strings"
	"sync"
	"time"
)

//var origin = "http://127.0.0.1:7777/"
//var url = "ws://127.0.0.1:7777/ws"

var connectsuccess int
var connectfail int
var sendsuccess int
var sendfail int
var acceptsuccess int
var acceptfail int
var maxtime int64
var mintime int64

var spend = make(chan int64)
var interval = 5
var waitgroup sync.WaitGroup

func main() {
	var n1 = flag.Int("n",0,"总连接数")
	var c1 = flag.Int("c",0,"并发连接数")
	var origin1 = flag.String("origin","","源地址")
	var url1 = flag.String("url","","ws链接地址")
	flag.Parse()

	if *n1 == 0{
		log.Fatal("缺少 -n 值")
	}
	n := *n1
	if *c1 == 0{
		log.Fatal("缺少 -c 值")
	}
	c := *c1

	if n <= 0 || c <= 0 {
		log.Fatalln("-n -c must > 0")
	}

	if *url1 == ""{
		log.Fatal("缺少 -url 值")
	}
	url := *url1

	var origin string
	if *origin1 == ""{
		//log.Fatal("缺少 -origin 值")
		wssurl := strings.ReplaceAll(*url1,"wss","")
		wsurl  := strings.ReplaceAll(wssurl,"ws","")
		origin = "http"+ wsurl
	}else{
		origin = *origin1
	}

	msg:="test"

	begin := time.Now().UnixNano()

	for i:=1;i<=n/c ;i++  {
		for j:=1;j<=c ;j++  {
			waitgroup.Add(1)
			go connect(url,origin,msg)
		}
		time.Sleep(time.Duration(interval)*time.Millisecond)

	}

	var t int64
	t = 0
	total := 0
	for k:=1;k<=n ;k++  {
		go func() {
			select {
			case t = <-spend:
				total++
				if maxtime < t {
					maxtime = t
				}
				if mintime == 0 || mintime > t {
					mintime = t
				}
			}
		}()

	}

	waitgroup.Wait()
	finish := time.Now().UnixNano()

	fmt.Println("total request : ",total)
	totalTime := float64(finish-begin)/1000000000-float64(interval*(n/c)/1000)
	fmt.Println("spend time : ",totalTime,"s")
	fmt.Printf("qps : %.2f [#/sec] \r\n",float64(total)/totalTime)

	fmt.Println("per request min time : ",float64(mintime)/1000000,"ms")
	fmt.Println("per request max time : ",float64(maxtime)/1000000,"ms")

	fmt.Println("connect success: ",connectsuccess)
	fmt.Println("connect fail: ",connectfail)
	fmt.Println("send message success: ",sendsuccess)
	fmt.Println("send message fail: ",sendfail)
	fmt.Println("accept message success: ",acceptsuccess)
	fmt.Println("accept message fail: ",acceptfail)
}


func connect(url string,origin string,msg string){
	start := time.Now().UnixNano()

	ws, err := websocket.Dial(url, "", origin)
	defer ws.Close()
	if err != nil {
		connectfail++
		log.Fatal(err)
	}else{
		connectsuccess++
	}
	//fmt.Println("connect: ",time.Now())

	message := []byte(msg)
	_, err = ws.Write(message)
	if err != nil {
		sendfail++
		log.Fatal(err)
	}else {
		sendsuccess++
	}
	//fmt.Printf("Send: %s\n", message)

	var data = make([]byte, 512)
	_, err2 := ws.Read(data)
	if err2 != nil {
		acceptfail++
		log.Fatal(err2)
	}else {
		acceptsuccess++
	}

	end := time.Now().UnixNano()
	spend <- end-start
	waitgroup.Done()
	//fmt.Printf("Receive: %s\n", msg[:m])

	//ws.Close()//关闭连接
}