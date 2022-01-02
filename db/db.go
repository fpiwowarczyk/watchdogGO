package db

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fpiwowarczyk/watchdogGO/utils"
)

type DynamoDB struct {
	sess *session.Session
	svc  *dynamodb.DynamoDB
}

type Settings struct {
	id             string
	ListOfServices []string
	NumOfSecCheck  string
	NumOfSecWait   string
	NumOfAttempts  string
}

var watchdogTable string
var err error

func New() *DynamoDB {
	db := new(DynamoDB)
	db.sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db.svc = dynamodb.New(db.sess)

	watchdogTable, err = utils.GetConfig("tables/watchdog")
	if err != nil {
		log.Println(err)
	}

	return db

}

// only added for testing
func (db *DynamoDB) PutItem() {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("1"),
			},
			"ListOfServices": {
				SS: aws.StringSlice([]string{"bluetooth", "lvm2"}),
			},
			"NumOfSecCheck": {
				S: aws.String("60s"),
			},
			"NumOfSecWait": {
				S: aws.String("10s"),
			},
			"NumOfAttempts": {
				S: aws.String("4"),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(watchdogTable),
	}

	result, err := db.svc.PutItem(input)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(result)
}

func (db *DynamoDB) GetItem(id string) (*Settings, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(watchdogTable),
	}

	result, err := db.svc.GetItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
		return nil, err
	}

	if result.Item == nil {
		log.Println("Could not find '" + id + "'")
	}

	settings := Settings{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &settings)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return &settings, nil

}

func (sett *Settings) Equals(newSett *Settings) bool {
	if newSett.NumOfAttempts != sett.NumOfAttempts ||
		newSett.NumOfSecCheck != sett.NumOfSecCheck ||
		newSett.NumOfSecWait != sett.NumOfSecWait ||
		utils.Equals(newSett.ListOfServices, sett.ListOfServices) {
		return true
	}

	return false
}
