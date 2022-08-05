package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/tidwall/gjson"
)

type Message struct {
	Default string `json:"default"`
}

func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) {

	var log = Logger{
		Level: WarnLevel,
	}

	isDebug := os.Getenv("Debug")
	publishToSns := os.Getenv("PublishToSns")

	if strings.ToUpper(isDebug) == "TRUE" {
		log.SetLevel(DebugLevel)
		log.Debug("Debug is on")
	}

	//Outbound sns topic
	publishTopic := os.Getenv("PublishTopicArn")

	log.Debug(publishTopic)

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS

		log.Debug("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)

		//example of how to extract data if the SNS message is a Json payload
		inboundSource := gjson.Get(snsRecord.Message, "source").String()
		log.Debug("Source: %s", inboundSource)

		//gjson makes it super simple to query deep into the json object
		faultCount := gjson.Get(snsRecord.Message, "detail.ClientRequestImpactStatistics.FaultCount").String()
		log.Debug("FaultCount: %s", faultCount)

		//if this env var is true, bomb out early
		//This is just here for testing, if you don't have the outbound sns set up
		if strings.ToUpper(publishToSns) != "TRUE" {
			return
		}

		//Initializes a new AWS session using credentials file. (~/.aws/credentials).
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := sns.New(sess)

		//SNS requires you to define a 'default' message that will be distributed
		//to all subscribers in lieu of a channel specific messages
		//See here : https://docs.aws.amazon.com/sns/latest/api/API_Publish.html
		// Example:
		// {
		// 	"default": "A message.",
		// 	"email": "A message for email.",
		// 	"email-json": "A message for email (JSON).",
		// 	"http": "A message for HTTP.",
		// 	"https": "A message for HTTPS.",
		// 	"sqs": "A message for Amazon SQS."
		// }
		message := Message{
			Default: string(snsRecord.Message),
		}
		messageBytes, _ := json.Marshal(message)
		messageStr := string(messageBytes)

		result, err := svc.Publish(&sns.PublishInput{
			TopicArn:         aws.String(publishTopic),
			Message:          aws.String(messageStr),
			MessageStructure: aws.String("json"),
			MessageAttributes: map[string]*sns.MessageAttributeValue{
				"Sns-Source": {
					DataType:    aws.String("String"),
					StringValue: aws.String(inboundSource),
				},
				"FaultCount": {
					DataType:    aws.String("String"),
					StringValue: aws.String(faultCount),
				},
			},
		})
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		log.Debug(*result.MessageId)

	}
}

func main() {
	lambda.Start(HandleRequest)
}
