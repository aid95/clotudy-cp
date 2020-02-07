package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

const sockBufSize int = 1024

var (
	upgrader = &websocket.Upgrader{
		HandshakeTimeout:  5,
		ReadBufferSize:    sockBufSize,
		WriteBufferSize:   sockBufSize,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		Error:             nil,
		CheckOrigin:       nil,
		EnableCompression: false,
	}
)

func init() {
}

func main() {
	fmt.Println("Compiler service start...")

	router := httprouter.New()
	router.GET("/cp/:cable_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sock, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("ServerHTTP:", err)
			return
		}
		newService(sock, ps.ByName("cable_id"))
	})
}
