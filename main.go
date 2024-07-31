package main

import (
	"TalkShell/client"
	"TalkShell/server"
	"bufio"
	"strings"

	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main(){
	
	closeSignal := make(chan os.Signal,1)
	
	signal.Notify(closeSignal,syscall.SIGINT,syscall.SIGTERM)

	var wg sync.WaitGroup

	wsServer := server.InitServer()

	fmt.Println("CHOOSE: server or cient")
	reader := bufio.NewReader(os.Stdin)
	choice,_ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
		case "server":
			wg.Add(1)
			go func(){
				defer wg.Done()
				wsServer.ListenWsConns(closeSignal)
				
			}()
			
		case "client":
			wg.Add(1)
			go func(){
				defer wg.Done()
				client.StartClient(closeSignal)
			}()
			
		default:
			fmt.Println("Choose between server or client")
	}


	<-closeSignal
	fmt.Println("closing all goroutines")

	wg.Wait()
}