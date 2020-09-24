package main

import (
	"MafiaTelegram/src/Handler"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func main()  {
	switch OS := os.Getenv("OS"); {
	case OS == "Lambda":
		lambda.Start(Handler.LambdaHandler)

	default:

	}
}

