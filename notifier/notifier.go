package notifier

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/fpiwowarczyk/watchdogGO/utils"
)

type Connection struct {
	sess  *session.Session
	svc   *sns.SNS
	topic *string
}

func New() {
	conn := new(Connection)
	conn.sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	conn.svc = sns.New(conn.sess)

	topic := utils.GetConfig("sns/watchdog")
	conn.topic = &topic
}

func (conn *Connection) Publish(message *string) {
	_, err := conn.svc.Publish(&sns.PublishInput{
		Message:  message,
		TopicArn: conn.topic,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
