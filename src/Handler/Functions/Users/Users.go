package Users

import (
	"MafiaTelegram/src/MongoDB"
	"MafiaTelegram/src/Telegram"
	"log"
	"strconv"
)

// Adds or edit User added to Group Chat to DB
func UserAddedToChat(newUsers []Telegram.User)  {
	for _, user := range newUsers{

		if user.IsBot != true && MongoDB.FindFromGroupChat(user).Username != user.Username{
			log.Println("Added:", user)
			MongoDB.AddNewUserFromGroupChat(user)

		}else{
			log.Printf("%s is already in GroupDB", user.Username)
			MongoDB.UpdateUserStatusFromGroupChat(user, true)
		}
	}
	MongoDB.DBDisconnect()
}

// Edit User in DB than it was deleted from GroupChat
func UserDeletedFromChat(userDeleted Telegram.User)  {
	if MongoDB.FindFromGroupChat(userDeleted).Username == userDeleted.Username{
		log.Println("Deleted:", userDeleted)
		MongoDB.UpdateUserStatusFromGroupChat(userDeleted, false)

	}else {
		log.Printf("No such user: %s to delete", userDeleted.Username)
	}

	MongoDB.DBDisconnect()
}

// Checks Active Users in DB and GroupChat
func CheckActiveUsers() string {
	var(
		DBactiveUsersNicknames string
		DBactiveUsers = MongoDB.FindActiveChatUsers()
		TelegramActiveUsers = Telegram.GetChatMembersCount()
	)

	for i, data := range DBactiveUsers{
		var comma string
		if i != len(DBactiveUsers) - 1 {
			comma = ","
		}
		DBactiveUsersNicknames += " " + data.Username + comma
	}

	Text := "DB: " + strconv.Itoa(len(DBactiveUsers)) + "\nTelegramChat: " + strconv.Itoa(TelegramActiveUsers) + "\nActivated Users:" + DBactiveUsersNicknames

	return Text
}