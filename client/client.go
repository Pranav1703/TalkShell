package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func StartClient(stopRoutine <-chan interface{}){
	conn,err := net.Dial("tcp","localhost:3000")
	
	if err!=nil{
		log.Fatalln("err when connecting to server: ",err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)	
	for{
		select {
		case <-stopRoutine:
			fmt.Println("Closing client ...")
			return
		default:
			fmt.Print(">")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
				continue
			}
			input = strings.TrimSpace(input)
			_,err = conn.Write([]byte(input))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing to server:", err)
				return
			}
		}
	}
}