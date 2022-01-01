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

func New() (*Notifier, error) {
	conn := new(Notifier)
	conn.sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	conn.svc = sns.New(conn.sess)

	topic, err := utils.GetConfig("sns/watchdog")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	conn.topicPtr = &topic

	return conn, nil
}

func (conn *Notifier) publish(messagePtr *string) error {
	_, err := conn.svc.Publish(&sns.PublishInput{
		Message:  messagePtr,
		TopicArn: conn.topicPtr,
	})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (notifier *Notifier) Notify(message string) error {
	err := notifier.publish(&message)
	if err != nil {
		return err
	}
	log.Println(message)
	return nil
}
