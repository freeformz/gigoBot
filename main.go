package main

import (
	"github.com/ruelephant/gitterClient"
	"github.com/subosito/gotenv"
	"os"
	"log"
	"strings"
	"strconv"
	"math/rand"
	"time"
)

func init() {
	gotenv.Load()
}
var results map[string]int;

func messageHandler(chat gitterClient.ChatStruct, message gitterClient.MessageStruct) {
	if (strings.Contains(message.Text, "крутите барабан") || strings.Contains(message.Text, "крутить барабан")) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		roulette := [...]int{0,10,100,30,50,25,-1,0,50,0,10,25,5,150,-1,15,100,30,0,25,5,50,0,1000,150,25,10,30,-1,5,250,25,0,10,50,30,100,-1,25,50 }

		localSpore := roulette[r.Intn(len(roulette)-1)]
		if (localSpore == -1) {
			chat.SendMessage("@"+message.FromUser.Username+" Вы банкрот!")
			delete(results, message.FromUser.ID)
		} else {
			chat.SendMessage("@"+message.FromUser.Username+"  У вас "+strconv.Itoa(localSpore)+" очков")
			if oldValue, ok := results[message.FromUser.ID]; ok {
				results[message.FromUser.ID] = localSpore+oldValue
			} else {
				results[message.FromUser.ID] = localSpore
			}
		}
	}

	if (strings.Contains(message.Text, "приз")) {
		count := 0
		if oldValue, ok := results[message.FromUser.ID]; ok {
			count = oldValue;
		}
		chat.SendMessage("@"+message.FromUser.Username+" Всего вы заработали: "+strconv.Itoa(count)+" очков")
	}
}

func main() {
	results = make(map[string]int)
	token := os.Getenv("GITTER_API_TOKEN")


	gitter := gitterClient.Create(token)


	log.Print("Join channel debugChannel")
	debugChannel := gitterClient.ChatStruct{ TokenApi: token, RoomId:"563b92da16b6c7089cb99c97", Channel: make(chan gitterClient.MessageStruct)  }
	debugChannel.InfoMessage("Новая игра! Вы можете \"крутить барабан\" или получить \"приз\"", 60)
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
				messageHandler(debugChannel, message)
			/*
			case message:=<-secondChannel.Channel:
				messageHandler(secondChannel.Channel, message)
				*/
			default:
		}
	}
}