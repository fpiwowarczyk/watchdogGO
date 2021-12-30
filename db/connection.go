package db

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)


type connection struct  {
	sess  
	svc 
}
func New() (*connection) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)



	// input := &dynamodb.PutItemInput{
	// 	Item: map[string]*dynamodb.AttributeValue{
	// 		"id": {
	// 			S: aws.String("1"),
	// 		},
	// 		"LIstOfServices": {
	// 			S: aws.String("mysql"),
	// 		},
	// 		"NumOfSecCheck": {
	// 			S: aws.String("60"),
	// 		},
	// 		"NumOfSecWait": {
	// 			S: aws.String("10"),
	// 		},
	// 		"NumOfAttempts": {
	// 			S: aws.String("4"),
	// 		},
	// 	},
	// 	ReturnConsumedCapacity: aws.String("TOTAL"),
	// 	TableName:              aws.String("watchdog-table-wnwagnaimzikzanzhcdwyddo"),
	// }

	// result, err := svc.PutItem(input)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(result)

}

func (db *connection) getItem(id int){
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(strconv.Itoa(id)),
			},
		},
		TableName: aws.String("watchdog-table-wnwagnaimzikzanzhcdwyddo"),
	}

	result, err := svc.GetItem(input)

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
		return
	}
	fmt.Println(result)
}