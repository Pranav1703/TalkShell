package client

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"TalkShell/server"
)

func Start(wsServer *server.WsServer,stopRoutine <-chan interface{}){
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
			wsServer.ReadAndEmitMsg(input)
		}

	}

}