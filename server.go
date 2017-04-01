package main

import (
	"net/http"
	"code.google.com/p/go.net/websocket"
	"fmt"
	//"time"
	"time"
)

var userList = make(map[string]*websocket.Conn, 10)

func Server(ws *websocket.Conn)  {

	var err error
	userName := ws.Config().Location.Query().Get("userName")
	fmt.Println("user",userName)
	userList[userName] = ws



	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			time.Sleep(time.Second * 10)
		}


		message := fmt.Sprintf("{\"user\":\"%s\",\"message\":\"%s\"}",userName,reply);
		fmt.Println(message)

		for _,v := range userList {
			if len(message) < 1 || message == "" {
				continue
			}
			if err = websocket.Message.Send(v, message); err != nil {
				fmt.Println("Can't send")
				continue
			}

		}
	}

}

func main()  {
	http.Handle("/",websocket.Handler(Server))

	err := http.ListenAndServe(":8080",nil)

	if err != nil {
		panic("ListendAndServer " + err.Error() )
	}
}
