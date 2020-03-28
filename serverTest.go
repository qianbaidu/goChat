package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"golang.org/x/net/websocket"
	"math/rand"
	"sync"
	"time"
)

const (
	origin   = "http://127.0.0.1:8080/"
	testTime = 60 * 10 // 测试时长分钟
	Users    = 200    //测试websocket用户数
)

func client(user string, ctx *sync.WaitGroup) {
	startTime := time.Now().Unix()
	url := fmt.Sprintf("ws://127.0.0.1:8080/chat?userName=%s", user)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Error("建立socker失败,",err)
		return
	}
	message := []byte("hello, world")
	_, err = ws.Write(message)
	if err != nil {
		log.Error("发送消息失败",err)
	}
	//log.Info("Send: %s\n", message)

	var msg = make([]byte, 512)

	for {
		if (time.Now().Unix())-startTime >= testTime {
			log.Infof("user[%s]测试时长[%d]秒完成退出任务", user, testTime)
			break
		}
		_, err := ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Receive: %s\n", msg[:m])
		time.Sleep(time.Second * 10)

	}

	ws.Close()
	ctx.Done()
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	res := fmt.Sprintf("%d-%s", time.Now().Nanosecond(), string(result))
	return res
}

func main() {
	ctx := &sync.WaitGroup{}
	randNum := 4
	for i := 0; i < Users; i++ {
		if i % 50 == 0 {
			randNum ++
		}
		user := GetRandomString(randNum)
		//time.Sleep(time.Second)
		ctx.Add(i)
		go client(user, ctx)
		fmt.Println("client ", i)
	}
	log.Infof("等待webosokcet线程运行 %d s", testTime)
	ctx.Wait()
}
