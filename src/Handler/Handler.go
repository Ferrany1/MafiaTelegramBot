package Handler

import (
	"MafiaTelegram/src/Handler/Functions/Host"
	"MafiaTelegram/src/Handler/Functions/Users"
	"MafiaTelegram/src/MongoDB"
	"MafiaTelegram/src/Telegram"
	"MafiaTelegram/src/Utils"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"strings"
)

var (
	Config = Utils.ReadFile("/Config.json")
	GroupChatID = Config.Telegram.GroupChat
)

func LambdaHandler(event events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error)  {
	var (
		eventReader = strings.NewReader(event.Body)
		TelegramBody = &Telegram.WebhookReqBody{}
	)

	// Decode event from Telegram webhook from AWS APIGateway
	switch err := json.NewDecoder(eventReader).Decode(TelegramBody); {
	case err == nil:
		var (
			TUser = TelegramBody.Message.From.Username
			TChatID = TelegramBody.Message.Chat.ID
			TMessage = TelegramBody.Message.Text
			NewUsers = TelegramBody.Message.NewChatMember
			UserDeleted = TelegramBody.Message.LeftChatMember
			)

		switch {

		// Working with GroupChatID
		case TChatID == GroupChatID:

			switch  {

			// Handle /admin command from GroupChat
			case strings.Contains(TMessage, "/host"):

				switch  {

				// Sends help message if message equals Host
				case TMessage == "/host":
					TelegramBody.SendMessage("Type: Role:Number,Role:Number to send roles to players")
					log.Println("Telegram messages with host rules was sent")

				// Distributes roles among players
				default:
					if strings.Contains(TMessage, "/host ") && strings.Contains(TMessage, ":") {
						SplitMessage := strings.Split(TMessage, "/host ")

						if err := Host.RoleDivider(TelegramBody, SplitMessage[1]); err == nil{
							TelegramBody.SendMessage("Roles were given!")
							log.Println("Telegram messages with roles were sent")
						}else {
							log.Println("Error on hosting message data")
						}

					}else {
						TelegramBody.SendMessage("Type: Role:Number,Role:Number to send roles to players")
					}
				}


			// Add User or edit status on GroupChat newUser Join Chat on Adding new User
			case len(NewUsers) > 0:
				Users.UserAddedToChat(NewUsers)

				TelegramBody.SendMessage("If joined first time or changed username write `/start` to @HomeMafia_Bot")

			// Command to check Active Users on DB and Telegram
			case TMessage == "/activeusers@HomeMafia_Bot":

				TelegramBody.SendMessage(Users.CheckActiveUsers())

			// Make User from GroupChat inactive after User Left Chat
			case UserDeleted.Username != "":
				Users.UserDeletedFromChat(UserDeleted)

			// If Nothing Happens in cases just prints body
			default:
				log.Println("TelegramBody: ", TelegramBody)
			}

		// Add ChatID to DB if he sends Message from private chat
		case TChatID != GroupChatID:
			var (
				TelegramUser = Telegram.User{Username: TUser}
				)

			// Checks if user in DB
			switch  DBUsername := MongoDB.FindFromGroupChat(TelegramUser).Username; {
			case DBUsername == TUser:

				// Chat text processor
				switch  {

				// Works for command /start
				case strings.Contains(TMessage, "/start"):

					// Completes if User was added to DB from GroupChat
					switch err := MongoDB.UpdateUserChatIDFromPrivateChat(TUser, TChatID); {
					case err == nil:
						MongoDB.UpdateUserStatusFromGroupChat(TelegramUser, true)
						TelegramBody.SendMessage(fmt.Sprintf("%s was added to game", TUser))
					default:
						log.Printf("No such user to add from Private chat: %s", TUser)
					}
				// Works for command /stop
				case strings.Contains(TMessage, "/stop"):
					MongoDB.UpdateUserStatusFromGroupChat(TelegramUser, false)
					TelegramBody.SendMessage(fmt.Sprintf("%s was deleted from game", TUser))
				}
			}
			MongoDB.DBDisconnect()
		}

	// APIGateway error handling
	default:
		log.Println("Error:", err, eventReader)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}