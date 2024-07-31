package client

import (
	"bufio"
	"fmt"
	"log"
	"os/signal"
	
	"strings"
	"syscall"

	// "log"
	"net"
	"os"
)

func readFromServer(conn net.Conn) {
    scanner := bufio.NewScanner(conn)

    for scanner.Scan() {
        msg := scanner.Text()
        fmt.Printf("%s\n",msg)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading from server:", err)
    }
}


func readAndSendInput(conn net.Conn,username string){
	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan(){
	// 	fmt.Print(username,">")
	// 	input := scanner.Text()
	// 	_,err := conn.Write([]byte(username+">"+input+"\n"))
	// 	if err!=nil{
	// 		fmt.Println("error writing to server:",err)
	// 	}
	// }

	// if err := scanner.Err(); err != nil {
	// 	fmt.Println("Error reading from server:", err)
	// 	return
	// }

	reader := bufio.NewReader(os.Stdin)
	for{
		input,err := reader.ReadString('\n')
		if err!=nil{
			fmt.Println("error reading from terminal")
		}
		_,err = conn.Write([]byte(username+">"+input))
		if err!=nil{
			fmt.Println("error writing to server:",err)
		}
	}

}

func StartClient(){

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter username: ")
	username,_ := reader.ReadString('\n')

	conn,err := net.Dial("tcp","localhost:8080")
	fmt.Println("connected to:",conn.LocalAddr())
	if err!=nil{
		log.Fatalln("err when connecting to server:",err)
	}
		
	_,err = conn.Write([]byte(username))
	if err!= nil{
		log.Fatalln(err)
	}

	closeSignal := make(chan os.Signal,1)

	go func(){
		closeSignal := make(chan os.Signal,1)
	
		signal.Notify(closeSignal,syscall.SIGINT,syscall.SIGTERM)
		<-closeSignal
		conn.Close()
		os.Exit(0)
	}()
	username = strings.TrimSpace(username)
	go readFromServer(conn)
	go readAndSendInput(conn,username)

	<-closeSignal
	fmt.Println("Client closing...")
	
}

