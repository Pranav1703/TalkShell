package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func readFromServer(conn *websocket.Conn){


	for{
		_, msg, err := conn.ReadMessage()
		
		if err!=nil{
			fmt.Println("cannot read from server",err)
			os.Exit(1)
		}
		fmt.Println(string(msg))
	}

}

func StartClient(stopRoutine <-chan interface{}){
	conn, _,err := websocket.DefaultDialer.Dial("ws://localhost:3000",nil)

	if err!=nil{
		log.Fatalln("err when connecting to server: ",err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)	
	
	go readFromServer(conn)

	for{
		select {
		case <-stopRoutine:
			fmt.Println("Closing client ...")
			conn.Close()
			return

		default:
			fmt.Print(">")

			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
				continue
			}
			trimmedInput := strings.TrimSpace(input)

			err = conn.WriteMessage(websocket.TextMessage,[]byte(trimmedInput))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing to server:", err)
				return
			}
		}
	}
}

