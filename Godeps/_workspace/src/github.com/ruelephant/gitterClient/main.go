package gitterClient
import (
	"log"
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"errors"
	"strings"
	"time"
)

type userStruct struct {
	AvatarURLMedium string `json:"avatarUrlMedium"`
	AvatarURLSmall  string `json:"avatarUrlSmall"`
	DisplayName     string `json:"displayName"`
	ID              string `json:"id"`
	URL             string `json:"url"`
	Username        string `json:"username"`
}

type clientStruct struct {
	currentUser userStruct
	apiToken     string
	joinedRoom	 map[string]RoomStruct
	pmCallback   func(room RoomStruct)
}

func Create(apiToken string) (clientStruct) {
	client := new(clientStruct)
	client.apiToken = apiToken
	client.joinedRoom = make(map[string]RoomStruct)
	client.getUser()

	ticker := time.NewTicker(time.Duration(60) * time.Second)
	go func() {
		for range ticker.C {
			if (client.scanPmChannel != nil) {
				client.scanPmChannel()
			}
		}
	}()

	return *client
}

func (client *clientStruct)scanPmChannel() {
	response := client.getApi("/v1/rooms")
	content, _ := ioutil.ReadAll(response)

	channelList := []roomStructSource{}
	err := json.Unmarshal(content, &channelList)

	if err == nil {
		for _, channel := range channelList {
			if (channel.OneToOne) {
				err, room := client.NewRoom(channel.URL, true)
				if (err == nil) {
					client.pmCallback(room)
				}
			}
		}
	}
}

func (client *clientStruct)PmCallback(callback func(room RoomStruct))() {
	client.pmCallback = callback
	client.scanPmChannel()
}

func (client *clientStruct)NewRoom(roomUrl string, monoMode bool)(error, RoomStruct) {
	if _, ok := client.joinedRoom[roomUrl]; ok {
		return errors.New("rooms area joined"), client.joinedRoom[roomUrl]
	} else {
		response := client.getApi("/v1/rooms")
		content, _ := ioutil.ReadAll(response)

		channelList := []roomStructSource{}
		err := json.Unmarshal(content, &channelList)
		room := new(RoomStruct)
		if err == nil {
			for _, channel := range channelList {
				if (strings.ToLower(channel.URL) == strings.ToLower(roomUrl)) {
					room.client = client
					room.monopoleMode = monoMode
					room.Id = channel.ID
					room.Url = channel.URL
					room.MessageChannel = make(chan MessageStruct)

					client.joinedRoom[roomUrl] = *room;
					return nil, client.joinedRoom[roomUrl]
				}
			}
			return errors.New("room not found"), *room
		}
		return errors.New("rooms list json error"), *room
	}
}

// ----- PRIVATE

func(client *clientStruct) getUser(){
	userReader := bufio.NewReader(client.getApi("/v1/user"))
	userJsonSource, _ := userReader.ReadBytes('\n')
	var user = []userStruct{};
	err := json.Unmarshal(userJsonSource, &user)
	if err != nil {
		log.Fatal("Failed to load current user. Bad token?")
	}
	client.currentUser = user[0]
}

func(client *clientStruct) getApi(url string) (io.ReadCloser) {
	return getRequest("https://api.gitter.im"+url+"?access_token="+(client.apiToken))
}
