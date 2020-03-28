package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"

	//"time"
	"time"
)

const (
	MSG_TYPE_MESSAGE = "message"
	MSG_TYPE_STATUS  = "status"
)

var (
	onlineInfo *Online
	rwlock     sync.RWMutex
)

type Online struct {
	UserList  map[string]*websocket.Conn
	UserCount int64
}

type Msg struct {
	User        string `json:"user"`
	Messag      string `json:"message"`
	OnlineTotal int64  `json:"online"`
	PostTime    string `json:"time"`
	Type        string `json:"msgType"` //status：用户状态更新；msg：新消息
}

func init() {
	onlineInfo = &Online{
		UserList:  make(map[string]*websocket.Conn),
		UserCount: 0,
	}
}

// 上线、下线用户更新信息，通知在线用户
func updateOnlineUser(user string, ws *websocket.Conn, isOnline bool) {
	var status string
	rwlock.Lock()

	if isOnline {
		status = "上线"
		onlineInfo.UserList[user] = ws
		onlineInfo.UserCount++
	} else {
		status = "下线"
		delete(onlineInfo.UserList, user)
		onlineInfo.UserCount--
	}
	rwlock.Unlock()

	log.Infof("用户:%s ,%s了", user, status)
	msg := Msg{
		User:        user,
		Messag:      fmt.Sprintf("用户%s %s了", user, status),
		OnlineTotal: onlineInfo.UserCount,
		PostTime:    time.Now().Format("2006-01-02 15:04:05"),
		Type:        MSG_TYPE_STATUS,
	}
	SendMsg(msg)
}

func SendMsg(msg Msg) {
	if len(msg.Messag) < 1 && msg.Type != MSG_TYPE_STATUS {
		return
	}
	message, _ := json.Marshal(msg)
	for _, v := range onlineInfo.UserList {
		if err := websocket.Message.Send(v, string(message)); err != nil {
			log.Errorf("send msg[%s] error: %s", msg, err)
			continue
		}
	}
}

func Server(ws *websocket.Conn) {
	var err error
	userName := ws.Config().Location.Query().Get("userName")
	if len(userName) < 1 {
		return
	}

	//新用户上线
	updateOnlineUser(userName, ws, true)

	//监听消息
	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			if err == io.EOF {
				//用户下线或断开提醒
				updateOnlineUser(userName, nil, false)
				break
			}
			log.Error("websocket 接收消息失败 ", err)
			time.Sleep(time.Second * 10)
		}
		if len(reply) < 1 {
			log.Errorf("接收消息空 ", reply)
			continue
		}
		msg := htmlDecode(strings.TrimRight(strings.TrimLeft(reply, "\n"), "\n"))
		newMsg := Msg{
			User:        userName,
			Messag:      msg,
			OnlineTotal: onlineInfo.UserCount,
			PostTime:    time.Now().Format("2006-01-02 15:04:05"),
			Type:        MSG_TYPE_MESSAGE,
		}
		SendMsg(newMsg)
	}

}

func htmlDecode(str string) string {
	str = strings.Replace(str, "<", "\\u003c", -1)
	str = strings.Replace(str, ">", "\\u003e", -1)
	str = strings.Replace(str, "\n", "&#13;", -1)
	return str
}

func Index(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("./html/index.html")
	t, err := template.ParseFiles("./html/wechat.html")
	if err != nil {
		log.Errorf("parse template error: %s", err)
		return
	}
	if err = t.Execute(w, nil); err != nil {
		log.Error(err)
	}
}

func main() {
	h := http.FileServer(http.Dir("./html"))
	http.Handle("/images/", h)

	http.HandleFunc("/index", Index)

	http.Handle("/", websocket.Handler(Server))
	log.Info("start server :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListendAndServer " + err.Error())
	}
}
