package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func (r *http.Request)  bool{
		//originのチェックは行わない
		return true 
	},
}

func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, []byte("hello websocket"))
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
