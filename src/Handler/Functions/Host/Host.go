package Host

import (
	"MafiaTelegram/src/MongoDB"
	"MafiaTelegram/src/Telegram"
	"MafiaTelegram/src/Utils"
	"errors"
	"log"
	"strconv"
	"strings"
)

func RoleDivider(TelegramBody *Telegram.WebhookReqBody, message string) error{
	var (
		TelegramChatActiveUsers = Telegram.GetChatMembersCount()
		DBActiveUsers = MongoDB.FindActiveChatUsers()

		Message= SplitIntoRoles(message)
		NumOfActiveRoles = CountActiveRoles(Message)

		err error
		)

	switch {

	// Works if GroupChat Users minus 1 equals active users from DB
	case TelegramChatActiveUsers == len(DBActiveUsers):

		switch  {

		// Works if Active Roles less then Active players to check if all roles will be distributed
		case NumOfActiveRoles < TelegramChatActiveUsers:

			// Creates Active Roles []string to distribute among Private chats
			Roles := CreateRolesSlice(Message, TelegramChatActiveUsers)

			// Clears old messages in Private chats
			ClearOldChatMessages()

			// Distributes messages among users and records them to DB
			SendMessagesInPrivateChats(DBActiveUsers, TelegramBody, Roles)

			// Disconnects from DB
			MongoDB.DBDisconnect()

		// If Active Roles equals or more then Active Chat players
		default:
			err = errors.New("Not enough players for given roles")
			TelegramBody.SendMessage("Not enough players for given roles")
		}

	// If GroupChat Users minus 1 not equal active users from DB
	default:
		err = errors.New("Wrong amount Active users in DB and Chat")
		TelegramBody.SendMessage("Wrong amount Active users in DB and Chat")
	}

	return err
}

// Split message string into Roles map
func SplitIntoRoles(message string) map[string]int {
	var(
		MessageTelegram = strings.Split(message, ",")
		Message = make(map[string]int)
		err error
	)

	for _, split := range MessageTelegram {
		if strings.Contains(split, ":"){
			Splitter := strings.Split(split, ":")
			Message[Splitter[0]], err = strconv.Atoi(Splitter[1])
		}else{
			err = errors.New("Wrong Input")
		}

	}

	log.Println(err)

	return Message
}

// Count active roles from message converted to map
func CountActiveRoles(message map[string]int) int {
	var(
		NumOfActiveRoles int
	)

	for _, numberOfPlayersInRole := range message {
		NumOfActiveRoles += numberOfPlayersInRole
	}

	return NumOfActiveRoles
}

func CreateRolesSlice(message map[string]int, telegramChatActiveUsers int) []string {
	var(
		Roles []string
	)

	for rolesname, rolesnum := range message {
		for iterator := 0; iterator < rolesnum; iterator ++ {
			Roles = append(Roles, rolesname)
		}
	}

	if len(Roles) < telegramChatActiveUsers - 1 {
		for iterator := 0; iterator < telegramChatActiveUsers - 1 - len(Roles); iterator ++ {
			Roles = append(Roles, "Мирный житель")
		}
	}

	return Roles
}

func ClearOldChatMessages()  {
	var(
		OldPrivateChatMessages = MongoDB.FindOldMessages()
	)

	if len(OldPrivateChatMessages) > 0 {
		for _, d := range OldPrivateChatMessages {
			Telegram.DeleteMessage(d.ChatID, d.MessageID)
		}
	}
}

func SendMessagesInPrivateChats(dBActiveUsers []MongoDB.GroupChatUser, telegramBody *Telegram.WebhookReqBody, roles []string)  {
	for _, user := range dBActiveUsers {

		// Works only for matching Active user from DB that is not a /host message sender
		if  telegramBody.Message.From.Username != user.Username {

			// Generates random role and deletes in from initial slice
			Role, index := Utils.Rand(roles)
			roles = Utils.RemoveIndexFromSlice(roles, index)

			// Creates new User from type to send Message
			User := new(Telegram.WebhookReqBody)
			User.Message.Chat.ID = user.ChatID

			// Sends message to User and records its data to DB
			Message := User.SendMessage(Role)

			MongoDB.AddNewPrivateMessage(Message.Result)
		}
	}
}