package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getSettings() Settings {

	settings := Settings{ListOfServices: []string{"docker", "mysql"},
		NumOfSecCheck: "60s",
		NumOfSecWait:  "10s",
		NumOfAttempts: "4",
	}

	return settings
}
func TestSettingsEquals_Success(t *testing.T) {
	sett := getSettings()

	newSettings := Settings{
		ListOfServices: []string{"docker", "mysql"},
		NumOfSecCheck:  "60s",
		NumOfSecWait:   "10s",
		NumOfAttempts:  "4",
	}

	output := sett.Equals(&newSettings)

	assert.True(t, output)
}

func TestSettingsEquals_Fail_DifferentListOfServices(t *testing.T) {
	sett := getSettings()

	newSettings := Settings{
		ListOfServices: []string{"docker"},
		NumOfSecCheck:  "60s",
		NumOfSecWait:   "10s",
		NumOfAttempts:  "4",
	}
	output := sett.Equals(&newSettings)

	assert.False(t, output)
}

func TestSettingsEquals_Fail_DifferentSecCheck(t *testing.T) {
	sett := getSettings()

	newSettings := Settings{
		ListOfServices: []string{"docker", "mysql"},
		NumOfSecCheck:  "50s",
		NumOfSecWait:   "10s",
		NumOfAttempts:  "4",
	}

	output := sett.Equals(&newSettings)

	assert.False(t, output)
}

func TestSettingsEquals_Fail_DifferentSecWait(t *testing.T) {
	sett := getSettings()

	newSettings := Settings{
		ListOfServices: []string{"docker", "mysql"},
		NumOfSecCheck:  "50s",
		NumOfSecWait:   "9s",
		NumOfAttempts:  "4",
	}

	output := sett.Equals(&newSettings)

	assert.False(t, output)
}
func TestSettingsEquals_Fail_DifferentAttempts(t *testing.T) {
	sett := getSettings()

	newSettings := Settings{
		ListOfServices: []string{"docker", "mysql"},
		NumOfSecCheck:  "50s",
		NumOfSecWait:   "9s",
		NumOfAttempts:  "3",
	}

	output := sett.Equals(&newSettings)

	assert.False(t, output)
}
