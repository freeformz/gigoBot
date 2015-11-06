package main

import (
	"github.com/ruelephant/gitterClient"
	"github.com/subosito/gotenv"
	"os"
	"log"
)

func init() {
	gotenv.Load()
}

func messageHandler(channel chan string, message string) {
	channel<-"Echo: "+message
}

func main() {
	token := os.Getenv("GITTER_API_TOKEN")

	log.Print("Join channel debugChannel")
	debugChannel := gitterClient.Chat{ TokenApi: token, RoomId:"563b92da16b6c7089cb99c97", Channel: make(chan string)  }
	go debugChannel.JoinRoom()

	/*
		// Example - Second channel
		// You can get channel id in (open in browser) https://api.gitter.im/v1/rooms?access_token={YOU_TOKEN}
		// Previously join to the channel with the gitter client (gitter.im)

		log.Print("Join channel myChannel")
		secondChannel := GitterClient.Chat{ TokenApi: token, RoomId:"{PASTE YOU CHANNEL ID}", Channel: make(chan string)  }
		go secondChannel.JoinRoom()
	 */

	for {
		select {
			case message:=<-debugChannel.Channel:
				messageHandler(debugChannel.Channel, message)
			/*
			case message:=<-secondChannel.Channel:
				messageHandler(secondChannel.Channel, message)
				*/
			default:
		}
	}
}