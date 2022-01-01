package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig_Success(t *testing.T) {
	output, err := GetConfig("test/testConfig")

	assert.Nil(t, err)
	assert.Equal(t, "Some value", output)
}

func TestGetConfig_Fail_MissingProperty(t *testing.T) {
	_, err := GetConfig("missingProp")

	assert.Equal(t, errors.New("Missing property: missingProp"), err)
}
