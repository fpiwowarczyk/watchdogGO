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

type NotifierI interface {
	publish(messagePtr *string) error
}

func New() (*Notifier, error) {
	notifier := new(Notifier)
	notifier.sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	notifier.svc = sns.New(notifier.sess)

	topic, err := utils.GetConfig("sns/watchdog")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	notifier.topicPtr = &topic

	return notifier, nil
}

func (notifier *Notifier) publish(messagePtr *string) error {
	_, err := notifier.svc.Publish(&sns.PublishInput{
		Message:  messagePtr,
		TopicArn: notifier.topicPtr,
	})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func Notify(notifier NotifierI, message string) error {
	err := notifier.publish(&message)
	if err != nil {
		return err
	}
	log.Println(message)
	return nil
}
