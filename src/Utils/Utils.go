package Utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)
type Creds struct {
	Telegram Telegram
	MongoDB MongoDB
}

type Telegram struct {
	API string
	GroupChat int64
}

type MongoDB struct {
	Username string
	Password string
	DBName string
}

// Get file name from working directory
func GetFile(newFile string) string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.Buffer{}
	buf.WriteString(dir)
	buf.WriteString(newFile)

	result := buf.String()

	return result

}

// Read file to type Credis
func ReadFile(newFile string) Creds {

	f, err := os.Open(GetFile(newFile))
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	Credis := Creds{}
	err = json.NewDecoder(f).Decode(&Credis)

	if err != nil {
		fmt.Println("error:", err)
	}

	return Credis
}

// Fatals error
func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Prints error
func Println(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Get random string and index from []string
func Rand(slice []string)  (string, int){
	rand.Seed(time.Now().Unix())
	randomIndex := rand.Intn(len(slice))

	return slice[randomIndex], randomIndex
}

// Remove object by its index from []string
func RemoveIndexFromSlice(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}
