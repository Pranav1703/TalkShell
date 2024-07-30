package main

import (
	"TalkShell/client"
	"TalkShell/server"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main(){
	
	closeSignal := make(chan os.Signal,1)
	stopRoutine := make(chan interface{})
	signal.Notify(closeSignal,syscall.SIGINT,syscall.SIGTERM)

	var wg sync.WaitGroup

	wsServer := server.InitServer()

	fmt.Println("CHOOSE: server or cient")
	var choice string
	fmt.Scanf("%s",choice)
	switch choice {
		case "server":
			wg.Add(1)
			go func(){
				defer wg.Done()
				http.HandleFunc("/",wsServer.HandleWsConn)
				log.Fatal(http.ListenAndServe(":3000", nil))
			}()
		case "client":
			wg.Add(1)
			go func(){
				defer wg.Done()
				client.Start(wsServer,stopRoutine)
			}()
			
		default:
			fmt.Println("Choose between server or client")
	}

	wg.Add(1)
	go func(){
		defer wg.Done()
		wsServer.BroadcastMsg(stopRoutine)
	}()

	<-closeSignal
	fmt.Println("closing all goroutines")
	close(stopRoutine)

	wg.Wait()
}