package server

import (
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

type Message struct{
	Text string
}

type WsServer struct{
	Conn *websocket.Conn
	Clients map[*websocket.Conn]bool
	Broadcast chan *Message
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func InitServer() *WsServer{
	return &WsServer{
		Clients: make(map[*websocket.Conn]bool),
		Broadcast: make(chan *Message),
	}
}

func (ws *WsServer)HandleWsConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	ws.Conn = conn
	ws.Clients[conn] = true

	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(">Received message:", string(message))	
		msg := &Message{
			Text: string(message),
		}		
		ws.Broadcast <- msg
	}
}

func (ws *WsServer)BroadcastMsg(stopRoutine <-chan interface{}){
	for{	
		select{
		case msg := <-ws.Broadcast:

			for client := range ws.Clients{
				err:= client.WriteMessage(websocket.TextMessage,[]byte(msg.Text))
				if err!= nil{
					fmt.Println("Write Error: ",err)
					client.Close()
					delete(ws.Clients,client)
				}
			}

		case <-stopRoutine:
			fmt.Println("Closing Broadcast ...")
			return
		}
	}
}

func (ws *WsServer)ReadAndEmitMsg(msg string){
	err:= ws.Conn.WriteMessage(websocket.TextMessage,[]byte(msg))
	if err!= nil{
		fmt.Println("Write Error: ",err)
		ws.Conn.Close()
		delete(ws.Clients,ws.Conn)
	}
}
