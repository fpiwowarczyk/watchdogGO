package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fpiwowarczyk/watchdogGO/utils"
)

type Connection struct {
	sess *session.Session
	svc  *dynamodb.DynamoDB
}

type Settings struct {
	id             string
	ListOfServices string
	NumOfSecCheck  string
	NumOfSecWait   string
	NumOfAttempts  string
}

var watchdogTable string

func New() *Connection {
	conn := new(Connection)
	conn.sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	conn.svc = dynamodb.New(conn.sess)

	watchdogTable = utils.GetConfig("tables/watchdog")

	return conn

}

// To remove, only added for testing
func (conn *Connection) PutItem() {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("1"),
			},
			"ListOfServices": {
				// SS: aws.StringSlice([]string{"bluetooth", "lvm2"}),
				S: aws.String("bluetooth"),
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

	result, err := conn.svc.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}

func (conn *Connection) GetItem(id string) (*Settings, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(watchdogTable),
	}

	result, err := conn.svc.GetItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil, err
	}

	if result.Item == nil {
		fmt.Println("Could not find '" + id + "'")
	}

	settings := Settings{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &settings)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return &settings, nil

}
