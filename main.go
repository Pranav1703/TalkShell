package main

import (
	"TalkShell/client"
	"TalkShell/server"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main(){
	
	closeSignal := make(chan os.Signal,1)
	stopRoutine := make(chan interface{})
	signal.Notify(closeSignal,syscall.SIGINT,syscall.SIGTERM)

	var wg sync.WaitGroup

	wsServer := server.InitServer()

	fmt.Println("CHOOSE: server or cient")
	var choice string
	fmt.Scanf("%s",&choice)

	handler := http.NewServeMux()
	
	s := http.Server{
		Addr: "localhost:3000",
		Handler: handler,
	}

	switch choice {
		case "server":
			wg.Add(1)
			go func(){
				defer wg.Done()

				handler.HandleFunc("/",wsServer.HandleWsConn)
				
				if err:= s.ListenAndServe(); err!= nil && !errors.Is(err,http.ErrServerClosed){
					fmt.Println("SERVER ERROR: ",err)
				}
				
			}()
			
		case "client":
			wg.Add(1)
			go func(){
				defer wg.Done()
				client.StartClient(stopRoutine)
			}()
			
		default:
			fmt.Println("Choose between server or client")
	}


	<-closeSignal
	fmt.Println("closing all goroutines")
	close(stopRoutine)

	if(choice == "server"){
		fmt.Println("closing server with 2sec timeout")
		ctx,cancel := context.WithTimeout(context.Background(),2*time.Second)
		defer cancel()
		err := s.Shutdown(ctx)
		if err!=nil{
			fmt.Println("error while trying to shutdown server:",err)
		}
		fmt.Println("server closed")
	}

	wg.Wait()
}