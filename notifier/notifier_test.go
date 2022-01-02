package notifier

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotify_Fail_EmptyMessage(t *testing.T) {
	notifier, _ := New()
	err := notifier.Notify("")

	assert.Equal(t, errors.New("Message need to have body"), err)
}
