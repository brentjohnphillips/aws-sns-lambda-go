package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) {
	fmt.Println(snsEvent)
}

func main() {
	lambda.Start(HandleRequest)
}
