package server

import (
	"fmt"
	"log"
	
	"net/http"

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
	CheckOrigin: func(r *http.Request) bool {
		// Allow any origin (or add logic to validate origin)
		return true
	},
}

func InitServer() *WsServer{
	return &WsServer{
		Clients: make(map[*websocket.Conn]bool),
		Broadcast: make(chan *Message),
	}
}



func (ws *WsServer)HandleWsConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	ws.Conn = conn
	ws.Clients[conn] = true
	defer func(){
		conn.Close()
		delete(ws.Clients,conn)
	}()
	for{
	
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatalln("couldnt read: ",err)
			return
		}

		fmt.Println(">Received message:", string(message))	
		msg := &Message{
			Text: string(message),
		}		

		for client := range ws.Clients{
			if client == ws.Conn{
				continue
			}
			err:= client.WriteMessage(messageType,[]byte(msg.Text))
			if err!= nil{
				fmt.Println("Write Error: ",err)
				client.Close()
				delete(ws.Clients,client)
			}
		}
	}
	
}

func (ws *WsServer)BroadcastMsg(stopRoutine <-chan interface{}){
	outer:
	for{	
		select{
		case msg := <-ws.Broadcast:

			for client := range ws.Clients{
				if client == ws.Conn{
					continue
				}
				err:= client.WriteMessage(websocket.TextMessage,[]byte(msg.Text))
				if err!= nil{
					fmt.Println("Write Error: ",err)
					client.Close()
					delete(ws.Clients,client)
				}
			}

		case <-stopRoutine:
			fmt.Println("Closing Broadcast ...")
			break outer;
		}
	}
}


