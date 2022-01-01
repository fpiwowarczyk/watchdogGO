package watchdog

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	workingName    = "bluetooth"
	notWorkingName = "lvm2"
	checkInterval  = "60s"
	waitInterval   = "10s"
	attepmts       = 4
)

func getDummyWatchdogForWorkingService() Watchdog {
	checkI, _ := time.ParseDuration(checkInterval)
	waitI, _ := time.ParseDuration(waitInterval)
	watchdog := Watchdog{name: workingName, numOfSecCheck: checkI, numOfSecWait: waitI, numOfAttempts: attepmts}
	return watchdog
}

func getDummyWatchdogForNotWorkingService() Watchdog {
	checkI, _ := time.ParseDuration(checkInterval)
	waitI, _ := time.ParseDuration(waitInterval)
	watchdog := Watchdog{name: notWorkingName, numOfSecCheck: checkI, numOfSecWait: waitI, numOfAttempts: attepmts}
	return watchdog
}

func TestNewWatchdog_Success(t *testing.T) {
	expected := getDummyWatchdogForWorkingService()

	tested, err := NewWatchdog("bluetooth", "60s", "10s", 4)
	assert.Nil(t, err)
	assert.Equal(t, &expected, tested, "Should return valid watchdog after getting correct values")
}

func TestNewWatchdog_WrongCheckInterval(t *testing.T) {
	_, err := NewWatchdog("bluetooth", "wrongValue", "10s", 4)
	assert.Equal(t, errors.New("time: invalid duration wrongValue"), err, "Incorrect check interval should return error")
}

func TestNewWatchdog_WrongWaitInterval(t *testing.T) {
	_, err := NewWatchdog("bluetooth", "60s", "wrongVal", 4)

	assert.Equal(t, errors.New("time: invalid duration wrongVal"), err, "Incorrect wait inteval should return error")
}

func TestNewWatchdog_WrongAttemptsVal(t *testing.T) {
	_, err := NewWatchdog("bluetooth", "60s", "10s", 0)

	assert.Equal(t, errors.New("Incorrect number of attempts to run service"), err, "Attenots value should be bigger than 0")
}

func TestNewWatchdog_Success_EmptyServiceName(t *testing.T) {
	_, err := NewWatchdog("", "60s", "10s", 4)

	assert.Equal(t, errors.New("Service name cannot be empty"), err, "Service name cannot be empty and should return error")
}

func TestNewWatchdog_Success_IncorrectServiceName(t *testing.T) {
	_, err := NewWatchdog("xyz", "60s", "10s", 4)

	assert.Equal(t, errors.New("Could not find service: xyz"), err, "Service should exist function should return error")
}

func TestIsRunning_Success(t *testing.T) {
	service := getDummyWatchdogForWorkingService()

	output := service.IsRunning()
	assert.True(t, output, "Service should run and function return true")

}

func TestStart_success(t *testing.T) {
	service := getDummyWatchdogForWorkingService()

	output := service.Start()

	assert.True(t, output)
}

func TestStart_fail(t *testing.T) {
	service := getDummyWatchdogForNotWorkingService()

	output := service.Start()

	assert.False(t, output)
}

func TestIsRunning_Fail(t *testing.T) {
	service := getDummyWatchdogForNotWorkingService()

	output := service.IsRunning()
	assert.False(t, output, "Service shoud not run and function return false")
}
