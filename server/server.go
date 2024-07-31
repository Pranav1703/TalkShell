package server

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"net"
)

type WsServer struct{
	// Conn net.Conn
	Clients map[net.Conn]string
}

func InitServer() *WsServer{
	return &WsServer{
		Clients: make(map[net.Conn]string),
	}
}

func (ws *WsServer)Register(conn net.Conn)string{
	reader := bufio.NewReader(conn)
	username,err := reader.ReadString('\n')
	if err!=nil{
		fmt.Println("Error while registering: ",err)
	}
	ws.Clients[conn] = username
	fmt.Println("regeisted client: ",conn.LocalAddr())
	return username
}


func (ws *WsServer)ListenWsConns(closeSignal <-chan os.Signal) {
	listener,err := net.Listen("tcp","localhost:8080")
	
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer listener.Close()
	fmt.Println("listneing on :8080 locally, details:",listener.Addr())
	go func(){
		for{
			conn,err := listener.Accept()
			if err!=nil{
				fmt.Println("counldn't accept conn: ",err)
				continue
			}
	
			go ws.HandleWsConn(conn)
		
		}
	}()

	<-closeSignal
	fmt.Println("server closed.")

}

func (ws *WsServer)HandleWsConn(conn net.Conn){
	defer conn.Close()
	//register conn with username
	username := ws.Register(conn)
	BroadcastMsg(fmt.Sprintf("%s has joined the chat.",username),ws.Clients,conn)

	
}

func BroadcastMsg(msg string, clients map[net.Conn]string, sender net.Conn ){
	for client := range clients{
		if client == sender{
			continue
		}
		_,err:= client.Write([]byte(msg))
		if err!= nil{
			fmt.Println("Write Error: ",err)
			fmt.Println("proceeding to remove client")
			client.Close()
			delete(clients,client)
		}
	}
}