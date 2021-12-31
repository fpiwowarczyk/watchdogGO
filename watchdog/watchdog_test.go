package watchdog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	workingName    = "bluetooth"
	notWorkingName = "lvm2"
	waitInterval   = "10s"
	checkInterval  = "5s"
	attepmts       = 4
)

func getDummyServiceWorking() Watchdog {
	checkI, _ := time.ParseDuration(checkInterval)
	waitI, _ := time.ParseDuration(waitInterval)
	service := Watchdog{name: workingName, numOfSecCheck: checkI, numOfSecWait: waitI, numOfAttempts: attepmts}
	return service
}

func getDummyServiceNotWorking() Watchdog {
	checkI, _ := time.ParseDuration(checkInterval)
	waitI, _ := time.ParseDuration(waitInterval)
	service := Watchdog{name: notWorkingName, numOfSecCheck: checkI, numOfSecWait: waitI, numOfAttempts: attepmts}
	return service
}
func TestIsRunning_Success(t *testing.T) {
	service := getDummyServiceWorking()

	output := service.IsRunning()
	assert.True(t, output, "Service should run and return true")

}

func TestIsRunning_Fail(t *testing.T) {
	service := getDummyServiceNotWorking()

	output := service.IsRunning()
	assert.False(t, output, "Service should run and return true")

}
