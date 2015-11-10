package main

import (
	"github.com/ruelephant/gitterClient"
	"github.com/subosito/gotenv"
	"os"
	/*"strings"
	"strconv"
	"math/rand"
	 */
	"net/http"
	"fmt"
	"html"
	"sync"
	"log"
	"time"
)

func init() {
	gotenv.Load()
}

var results map[string]int;

/*
func defaultMessageHandler(room gitterClient.RoomStruct, message gitterClient.MessageStruct) {
	if (strings.Contains(message.Text, "результат") || strings.Contains(message.Text, "крутить") || strings.Contains(message.Text, "крути") || strings.Contains(message.Text, "крутите барабан") || strings.Contains(message.Text, "крутить барабан")) {
		room.SendMessage("@"+message.FromUser.Username+" Рулетка доступна только на [специальном канале](https://gitter.im/GigoBot/rulegame) ")
		return
	}

	room.SendMessage("@"+message.FromUser.Username+"  Не пойму о чем вы :(")
}

func ruleGameHandler(room gitterClient.RoomStruct, message gitterClient.MessageStruct) {
	if (strings.Contains(message.Text, "крутить") || strings.Contains(message.Text, "крути") || strings.Contains(message.Text, "крутите барабан") || strings.Contains(message.Text, "крутить барабан")) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		roulette := [...]int{0,10,100,30,50,25,-1,0,50,0,10,25,5,150,-1,15,100,30,0,25,5,50,0,1000,150,25,10,30,-1,5,250,25,0,10,50,30,100,-1,25,50 }

		localSpore := roulette[r.Intn(len(roulette)-1)]
		if (localSpore == -1) {
			room.SendMessage("@"+message.FromUser.Username+" Вы банкрот!")
			delete(results, message.FromUser.ID)
		} else {
			room.SendMessage("@"+message.FromUser.Username+"  У вас "+strconv.Itoa(localSpore)+" очков")
			if oldValue, ok := results[message.FromUser.ID]; ok {
				results[message.FromUser.ID] = localSpore+oldValue
			} else {
				results[message.FromUser.ID] = localSpore
			}
		}
		return
	}

	if (strings.Contains(message.Text, "результат")) {
		count := 0
		if oldValue, ok := results[message.FromUser.ID]; ok {
			count = oldValue;
		}
		room.SendMessage("@"+message.FromUser.Username+" Всего вы заработали: "+strconv.Itoa(count)+" очков")
		return
	}
	
	room.SendMessage("@"+message.FromUser.Username+"  Не пойму о чем вы :(")
}

*/

type gigoBot struct {
	rooms []*gitterClient.RoomStruct;
}


func (bot *gigoBot)InfoMessage(room *gitterClient.RoomStruct, Message string, second int) (*time.Ticker) {
	ticker := time.NewTicker(time.Duration(second) * time.Second)
	go func() {
		for range ticker.C {
			room.SafeSendMessage(Message)
		}
	}()
	return ticker
}

func (bot *gigoBot)AddLisner(room *gitterClient.RoomStruct) {
	bot.rooms = append(bot.rooms, room)
	go room.Join()
}

func (bot *gigoBot)messageHandler(room *gitterClient.RoomStruct, message gitterClient.MessageStruct) {
	room.SendMessage("@"+message.FromUser.Username+"  Не пойму о чем вы :(")
}

//	gameRoom.InfoMessage("Работает \"Барабан\", напишите \"крутить\" или \"результат\"", 900)
func (bot *gigoBot)ChatLister() {
	defer wg.Done()
	for {
		for _, room := range bot.rooms  {
			select {
				case message := <-room.GetChannel():
					bot.messageHandler(room, message)
				default:

			}
		}
	}
}

func (bot *gigoBot)WebInterfaceLisner(webserverPort string) {
	defer wg.Done()
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	if (webserverPort == "") {
		webserverPort = "8080"
	}
	http.ListenAndServe(":"+webserverPort, serverMux)
}

var wg sync.WaitGroup

func main() {
	results = make(map[string]int)
	token := os.Getenv("GITTER_API_TOKEN")
	webserverPort := os.Getenv("PORT")

	gitter := gitterClient.Create(token)

	bot := &gigoBot{}
	if err,room:=gitter.NewRoom("LaravelRUS/GitterBot");err == nil {
		bot.AddLisner(&room)
	} else {
		log.Fatal(err)
	}

	if err,room:=gitter.NewRoom("GigoBot/RuleGame");err == nil {
		bot.InfoMessage(&room, "Работает \"Барабан\", напишите \"крутить\" или \"результат\"", 900)
		bot.AddLisner(&room)
	} else {
		log.Fatal(err)
	}

	/*
		- Add Leave method:

		Request URL:https://gitter.im/api/v1/rooms/563b92da16b6c7089cb99c97/users/54ddf4ce15522ed4b3dbf9e0
		Request Method:DELETE

		- LS support
		- Монопольные каналы
	 */
	wg.Add(2)
	go bot.WebInterfaceLisner(webserverPort)
	go bot.ChatLister()
	wg.Wait()
}