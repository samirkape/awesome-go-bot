package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samirkape/awesome-go-bot/parser"
)

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(parser.Sync)
}
