package client

import (
	"bufio"
	"fmt"
	"log"

	// "log"
	"net"
	"os"
	"strings"
)

func readFromServer(conn net.Conn,stopRoutine <-chan os.Signal){

	reader := bufio.NewReader(conn) 

	for{
		select{
		case <-stopRoutine:
			fmt.Println("stopped reading from server")
			return
		default:
			msg, err := reader.ReadString('\n')
		
			if err!=nil{
				fmt.Println("cannot read from server",err)
				os.Exit(1)
			}
			fmt.Println(msg)
		}
	}

}

func readAndSendInput(conn net.Conn,username string){
	reader := bufio.NewReader(conn)
	fmt.Print(username,"> ")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		return
	}
	input = strings.TrimSpace(input)

	_,err = conn.Write([]byte(username+"> "+input))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing to server:", err)
		return
	}
}

func StartClient(closeSignal <-chan os.Signal){

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter username: ")
	username,_ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	conn,err := net.Dial("tcp","localhost:8080")
	fmt.Println("connected to: ",conn.LocalAddr())
	if err!=nil{
		log.Fatalln("err when connecting to server: ",err)
	}
	defer conn.Close()
		
	_,err = conn.Write([]byte(username+"\n"))
	if err!= nil{
		log.Fatalln(err)
	}

	go readFromServer(conn,closeSignal)

	for{
		select {
		case <-closeSignal:
			fmt.Println("Closing client ...")
			conn.Close()
			return

		default:
			fmt.Print(username,"> ")

			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
				continue
			}
			input = strings.TrimSpace(input)

			_,err = conn.Write([]byte(username+"> "+input))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing to server:", err)
				return
			}
		}
	}
}

