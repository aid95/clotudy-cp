package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

const sockBufSize int = 1024

var (
	upgrader = &websocket.Upgrader{
		ReadBufferSize:    sockBufSize,
		WriteBufferSize:   sockBufSize,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		Error:             nil,
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: false,
	}
)

func init() {
	fmt.Println("Compiler service start...")
}

func main() {
	router := httprouter.New()
	router.GET("/cp/:cable_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("Websocket access request.")
		sock, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("ServerHTTP:", err)
			return
		}
		newService(sock, ps.ByName("cable_id"))
	})

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
