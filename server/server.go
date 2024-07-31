package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"net"
)

type WsServer struct{
	// Conn net.Conn
	Clients map[net.Conn]string
	mu 		sync.Mutex
}

func InitServer() *WsServer{
	return &WsServer{
		Clients: make(map[net.Conn]string),
	}
}

func StartServer(){
	wsServer := InitServer()
	wsServer.ListenWsConns()
}

func (ws *WsServer)Register(conn net.Conn)string{
	reader := bufio.NewReader(conn)
	username,err := reader.ReadString('\n')
	if err!=nil{
		fmt.Println("Error while registering: ",err)
	}
	username = strings.TrimSpace(username)
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.Clients[conn] = username
	fmt.Println("regeisted client: ",conn.RemoteAddr(), " with username '",username,"'")
	return username
}


func (ws *WsServer)ListenWsConns() {
	listener,err := net.Listen("tcp","localhost:8080")
	
	if err != nil {
		log.Fatalln(err)
		return
	}
	
	fmt.Println("listneing on :",listener.Addr())
	go func(){
		closeSignal := make(chan os.Signal,1)
	
		signal.Notify(closeSignal,syscall.SIGINT,syscall.SIGTERM)
	
		<-closeSignal
		listener.Close()
		fmt.Println("tcp connection closed.")
		os.Exit(0)
	}()

	for{
		conn,err := listener.Accept()
		if err!=nil{
			fmt.Println("counldn't accept conn: ",err)
			continue
		}

		go ws.HandleWsConn(conn)
	
	}

}

func (ws *WsServer)HandleWsConn(conn net.Conn){
	defer conn.Close()
	//register conn with username
	_= ws.Register(conn)

	//read ms from client and broadcast to other clients
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		ws.BroadcastMsg(fmt.Sprintf("%s\n",msg),conn) //check here
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from client:", err)
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()
	delete(ws.Clients, conn)
	
}

func (ws *WsServer)BroadcastMsg(msg string, sender net.Conn ){
	ws.mu.Lock()
	defer ws.mu.Unlock()
	fmt.Println(ws.Clients,"in server")
	for client := range ws.Clients{
		if client == sender{
			continue
		}
		fmt.Println("write to ->",client.RemoteAddr())
		_,err:=client.Write([]byte(msg))
		if err!= nil{
			fmt.Println("Write Error: ",err)
			fmt.Println("proceeding to remove client")
			client.Close()
			delete(ws.Clients,client)
		}
		// bufferedWriter := bufio.NewWriter(client)
		// _, err := bufferedWriter.Write([]byte(msg))
		// if err != nil {
		//     fmt.Println("Write Error:", err)
		// }
		// bufferedWriter.Flush()

	}
}