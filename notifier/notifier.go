package notifier

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/fpiwowarczyk/watchdogGO/utils"
)

type Notifier struct {
	sess     *session.Session
	svc      *sns.SNS
	topicPtr *string
}

func New() *Notifier {
	conn := new(Notifier)
	conn.sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	conn.svc = sns.New(conn.sess)

	topic, err := utils.GetConfig("sns/watchdog")
	if err != nil {
		log.Println(err)
	}
	conn.topicPtr = &topic

	return conn
}

func (conn *Notifier) publish(messagePtr *string) {
	_, err := conn.svc.Publish(&sns.PublishInput{
		Message:  messagePtr,
		TopicArn: conn.topicPtr,
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func (notifier *Notifier) Notify(message string) {
	notifier.publish(&message)
	log.Println(message)
}
