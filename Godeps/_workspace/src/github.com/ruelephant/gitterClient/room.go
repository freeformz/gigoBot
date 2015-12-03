package gitterClient
import (
	"net/http"
	"net/url"
	"bufio"
	"encoding/json"
	"strings"
)

type MessageStruct struct {
	FromUser struct {
				 AvatarURLMedium string `json:"avatarUrlMedium"`
				 AvatarURLSmall  string `json:"avatarUrlSmall"`
				 DisplayName     string `json:"displayName"`
				 Gv              string `json:"gv"`
				 ID              string `json:"id"`
				 URL             string `json:"url"`
				 Username        string `json:"username"`
			 } `json:"fromUser"`
	HTML     string        `json:"html"`
	ID       string        `json:"id"`
	Issues   []interface{} `json:"issues"`
	Mentions []struct {
		ScreenName string        `json:"screenName"`
		UserID     string        `json:"userId"`
		UserIds    []interface{} `json:"userIds"`
	} `json:"mentions"`
	Meta   []interface{} `json:"meta"`
	ReadBy int           `json:"readBy"`
	Sent   string        `json:"sent"`
	Text   string        `json:"text"`
	Unread bool          `json:"unread"`
	Urls   []interface{} `json:"urls"`
	V      int           `json:"v"`
}

type RoomStruct struct {
	client 		*clientStruct
	blockSendingMessage  bool
	monopoleMode bool
	Id     	string
	Url string
	MessageChannel  chan MessageStruct
}

type roomStructSource struct {
	GithubType  string        `json:"githubType"`
	ID          string        `json:"id"`
	Lurk        bool          `json:"lurk"`
	Mentions    int           `json:"mentions"`
	Name        string        `json:"name"`
	Noindex     bool          `json:"noindex"`
	OneToOne    bool          `json:"oneToOne"`
	RoomMember  bool          `json:"roomMember"`
	Tags        []interface{} `json:"tags"`
	Topic       string        `json:"topic"`
	UnreadItems int           `json:"unreadItems"`
	URL         string        `json:"url"`
	User        struct {
					AvatarURLMedium string `json:"avatarUrlMedium"`
					AvatarURLSmall  string `json:"avatarUrlSmall"`
					DisplayName     string `json:"displayName"`
					Gv              string `json:"gv"`
					ID              string `json:"id"`
					URL             string `json:"url"`
					Username        string `json:"username"`
					V               int    `json:"v"`
				} `json:"user"`
	UserCount int `json:"userCount"`
}

func (room *RoomStruct)Join() {
	response := getRequest("https://stream.gitter.im/v1/rooms/"+(room.Id)+"/chatMessages?access_token="+(room.client.apiToken))

	room.blockSendingMessage = false
	//room.SendMessage("Доброе время суток %username%!")

	reader := bufio.NewReader(response)
	for {
		line, _ := reader.ReadBytes('\n')
		if (string(line) != " \n") {
			var s = MessageStruct{}
			err := json.Unmarshal(line, &s)
			if err == nil {
				if (room.client.currentUser.ID != s.FromUser.ID) { // it is not my message
					room.blockSendingMessage = false
					if (room.monopoleMode) {
						s.Text = strings.ToLower(s.Text)
						room.MessageChannel <- s
					} else {
						// Search
						if (len(s.Mentions) > 0) {
							for _, mention := range s.Mentions {
								if (mention.UserID == room.client.currentUser.ID) { // Message to bot
									s.Text = strings.ToLower(s.Text)
									room.MessageChannel <- s
									break
								}
							}
						}
					}
				}
			}
		}
	}
}

func (room *RoomStruct)SafeSendMessage(text string) {
	if (!room.blockSendingMessage) {
		room.SendMessage(text)
	}
}

func (room *RoomStruct)SendMessage(text string) {
	room.blockSendingMessage = true
	http.PostForm("https://api.gitter.im/v1/rooms/"+(room.Id)+"/chatMessages?access_token="+(room.client.apiToken), url.Values{"text": {text}})
}
