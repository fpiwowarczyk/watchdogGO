package types

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Connection struct {
	sess *session.Session
	svc  *dynamodb.DynamoDB
}

type Service struct {
	id             string
	ListOfServices []string
	NumOfSecCheck  string
	NumOfSecWait   string
	NumOfAttempts  string
}
