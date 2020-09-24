package Telegram

import (
	"MafiaTelegram/src/Utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	Config      = Utils.ReadFile("/Config.json").Telegram
	URLTelegram = fmt.Sprintf("https://api.telegram.org/bot%s/", Config.API)
	GroupChat 	= Config.GroupChat
)

type WebhookReqBody struct {
	Message Message	`json:"message"`
}

type WebhookRespBody struct {
	Ok bool `json:"ok"`
	Result Message `json:"result"`
}

type Message struct {
	MessageID int `json:"message_id"`

	From struct{
		Username string `json:"username"`
	} `json:"from"`

	Date int `json:"date"`

	Chat struct {
		ID int64 `json:"id"`
	} `json:"chat"`

	Text string `json:"text"`

	NewChatMember []User `json:"new_chat_members"`

	LeftChatMember User `json:"left_chat_member"`
}

type User struct {
	ID	int `json:"id"`
	IsBot bool `json:"is_bot"`
	Username string `json:"username"`
}

type sendMessageReqBody struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	MessageID int    `json:"message_id"`
}

type UserCount struct {
	Result int `json:"result"`
}

// Returns Amount of Users in chat minus Bot
func GetChatMembersCount() int {
	reqBody := &sendMessageReqBody{
		ChatID: GroupChat,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.Post(URLTelegram + "getChatMembersCount", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalln(errors.New("unexpected status" + res.Status))
	}
	User := new(UserCount)
	json.NewDecoder(res.Body).Decode(User)

	return User.Result - 1
}

// General Method to Send data to Telegram chat
func Send(SendMethod string, reqBody *sendMessageReqBody) *WebhookRespBody {

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.Post(URLTelegram + SendMethod, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalln(errors.New("unexpected status" + res.Status))
	}

	Message := new(WebhookRespBody)
	json.NewDecoder(res.Body).Decode(Message)

	return Message
}

// Method from Send method to send messages to chat
func (body *WebhookReqBody) SendMessage(Text string) *WebhookRespBody {
	reqBody := &sendMessageReqBody{
		ChatID: body.Message.Chat.ID,
		Text: Text,
	}

	Message := Send("sendMessage", reqBody)

	return Message
}

// Method to Delete message from chat
func DeleteMessage(chatID int64, messageID int) {
	reqBody := &sendMessageReqBody{
		ChatID:    chatID,
		MessageID: messageID,
	}

	reqBytes, err := json.Marshal(reqBody)
	Utils.Fatal(err)

	res, err := http.Post(URLTelegram+"deleteMessage", "application/json", bytes.NewBuffer(reqBytes))
	Utils.Fatal(err)

	if res.StatusCode != http.StatusOK {
		Utils.Fatal(errors.New("unexpected status" + res.Status))
	}
}
