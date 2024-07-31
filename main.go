package main

import (
	"TalkShell/client"
	"TalkShell/server"
	"bufio"
	"strings"

	"fmt"
	"os"
)

func main(){
	

	fmt.Println("CHOOSE: 1.server or 2.client")
	reader := bufio.NewReader(os.Stdin)
	choice,_ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
		case "1":
			
			server.StartServer()
				
	
			
		case "2":

			client.StartClient()
			
			
		default:
			fmt.Println("Choose between server or client")
	}


}