package MongoDB

import (
	"MafiaTelegram/src/Telegram"
	"MafiaTelegram/src/Utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var(
	Config 			= Utils.ReadFile("/Config.json").MongoDB
	DBloginURL  	= "mongodb+srv://%v:%v@testapi.dhzzs.mongodb.net/%v?retryWrites=true&w=majority"
	Login  			= fmt.Sprintf(DBloginURL, Config.Username, Config.Password, Config.DBName)
	MafiaDB, clientDB = DBconnect()
)


type GroupChatUser struct {
	Username string `json:”username”`
	ChatID	 int64	`json:”chatid”`
	IsActive bool	`json:”isactive”`

}

type Message struct {
	ChatID 	  int64  `json:"chatid"`
	MessageID int 	 `json:"messageid"`
}

// Find If User From Group Chat in DB
func FindFromGroupChat(user Telegram.User) GroupChatUser {
	var (
		filter = bson.M{"username": user.Username}
		User GroupChatUser
	)

	MafiaDB.Collection("GroupUsers").FindOne(context.TODO(), filter).Decode(&User)

	return User
}

// Adds new User to DB GroupUsers from Group Chat
func AddNewUserFromGroupChat(user Telegram.User)  {
	var (
		User = GroupChatUser{Username: user.Username, IsActive: true,
		}
	)

	_, err := MafiaDB.Collection("GroupUsers").InsertOne(context.TODO(), User)
	Utils.Fatal(err)
}

// Updates User status in DB GroupUsers from Group Chat
func UpdateUserStatusFromGroupChat(user Telegram.User, status bool)  {
	var (
		filter = bson.D{primitive.E{Key: "username", Value: user.Username}}

		update = bson.D{primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "isactive", Value: status},
		}}}
	)

	MafiaDB.Collection("GroupUsers").FindOneAndUpdate(context.TODO(), filter, update)
}

// Updates ChatID in DB GroupUsers
func UpdateUserChatIDFromPrivateChat(username string, chatID int64) error {
	var (
		filter = bson.D{primitive.E{Key: "username", Value: username}}

		update = bson.D{primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "chatid", Value: chatID},
		}}}
	)

	err := MafiaDB.Collection("GroupUsers").FindOneAndUpdate(context.TODO(), filter, update).Err()

	return err
}

// Find All Users with status Active from DB
func FindActiveChatUsers() []GroupChatUser {
	var (
		filter = bson.D{{"isactive", true}}
		Users  = []GroupChatUser{}
	)

	results, err := MafiaDB.Collection("GroupUsers").Find(context.TODO(), filter)
	Utils.Fatal(err)

	defer results.Close(context.TODO())
	for results.Next(context.TODO()) {
		var user GroupChatUser
		if err = results.Decode(&user); err != nil {
			log.Fatal(err)
		}
		Users = append(Users, user)

	}
	return Users
}

// Adds new message Data to DB RolesPrivateMessages
func AddNewPrivateMessage(message Telegram.Message)  {
	var (
		User = Message{message.Chat.ID, message.MessageID}
	)

	_, err := MafiaDB.Collection("RolesPrivateMessages").InsertOne(context.TODO(), User)
	Utils.Fatal(err)
}

// Finds all previous message from DB RolesPrivateMessages and Delete them in DB
func FindOldMessages() []Message {
	var (
		Messages []Message
	)

	results, err := MafiaDB.Collection("RolesPrivateMessages").Find(context.TODO(), bson.D{})

	Utils.Fatal(err)

	defer results.Close(context.TODO())

	for results.Next(context.TODO()) {
		var message Message
		if err = results.Decode(&message); err != nil {
			log.Fatal(err)
		}
		Messages = append(Messages, message)

	}

	if len(Messages) > 0 {
		MafiaDB.Collection("RolesPrivateMessages").DeleteMany(context.TODO(), bson.D{})
	}

	return Messages
}

// Connects to DB and returns client and DB
func DBconnect() (*mongo.Database, *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientDB, err := mongo.Connect(ctx, options.Client().ApplyURI(Login))
	Utils.Fatal(err)

	Utils.Fatal(clientDB.Ping(context.TODO(), nil))

	log.Println("Successfully connected to MongoDB")

	MafiaDB := clientDB.Database(Config.DBName)

	return MafiaDB, clientDB
}

// Disconnects from DB
func DBDisconnect()  {
	Utils.Fatal(clientDB.Disconnect(context.TODO()))
	fmt.Println("DB was disconnected")
}