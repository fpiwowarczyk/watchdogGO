package notifier

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type NotifierMock struct{}

func (mock *NotifierMock) publish(messagePtr *string) error {
	if *messagePtr == "someMessage" {
		return nil
	}
	return errors.New("Incorrect message")
}

func TestNotify_ForwardMessageToPublish(t *testing.T) {
	mock := new(NotifierMock)
	correctMessage := "someMessage"
	err := mock.publish(&correctMessage)

	assert.Nil(t, err, "If message was forwarded correct should return no errors")

}
