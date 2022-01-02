package watchdog

import (
	"errors"
	"testing"
	"time"

	"github.com/fpiwowarczyk/watchdogGO/utils"
	"github.com/stretchr/testify/assert"
)

const (
	initd          = "/etc/init.d/"
	workingName    = "bluetooth"
	notWorkingName = "lvm2"
	incorrectName  = "xyz"
	checkInterval  = "60s"
	waitInterval   = "10s"
	attepmts       = 4
)

func getWorkingWatchdog() Watchdog {
	checkI, _ := time.ParseDuration(checkInterval)
	waitI, _ := time.ParseDuration(waitInterval)
	watchdog := Watchdog{name: workingName, numOfSecCheck: checkI, numOfSecWait: waitI, numOfAttempts: attepmts, os: &utils.OsFS{}}
	return watchdog
}

// Only intervals and attepmts
func getNewWatchdog() Watchdog {
	checkI, _ := time.ParseDuration(checkInterval)
	waitI, _ := time.ParseDuration(waitInterval)
	watchdog := Watchdog{numOfSecCheck: checkI, numOfSecWait: waitI, numOfAttempts: attepmts}
	return watchdog
}

func TestNewWatchdog_Success(t *testing.T) {
	fsOsMock := new(utils.MockOsFs)
	fsOsMock.On("Stat", "/etc/init.d/bluetooth").Return(nil, nil) // First arg isnt important, only no error matters
	expectedCheckTime, _ := time.ParseDuration("60s")
	expectedWaitTime, _ := time.ParseDuration("10s")

	tested, err := NewWatchdog("bluetooth", "60s", "10s", 4, fsOsMock)

	assert.Nil(t, err, "Should return no errors")
	assert.NotNil(t, tested, "Returned watchdog shoudn't not be nil")
	assert.Equal(t, "bluetooth", tested.name, "Should return same name as given in args")
	assert.Equal(t, expectedCheckTime, tested.numOfSecCheck)
	assert.Equal(t, expectedWaitTime, tested.numOfSecWait)
	assert.Equal(t, 4, tested.numOfAttempts)
}

func TestNewWatchdog_WrongCheckInterval(t *testing.T) {
	_, err := NewWatchdog("bluetooth", "wrongValue", "10s", 4, &utils.OsFS{})
	assert.Equal(t, errors.New("time: invalid duration wrongValue"), err, "Incorrect check interval should return error")
}

func TestNewWatchdog_WrongWaitInterval(t *testing.T) {
	_, err := NewWatchdog("bluetooth", "60s", "wrongVal", 4, &utils.OsFS{})

	assert.Equal(t, errors.New("time: invalid duration wrongVal"), err, "Incorrect wait inteval should return error")
}

func TestNewWatchdog_WrongAttemptsVal(t *testing.T) {
	_, err := NewWatchdog("bluetooth", "60s", "10s", 0, &utils.OsFS{})

	assert.Equal(t, errors.New("Attempts number less than 1"), err, "Attenpts value should be bigger than 0")
}

func TestNewWatchdog_Success_EmptyServiceName(t *testing.T) {
	_, err := NewWatchdog("", "60s", "10s", 4, &utils.OsFS{})

	assert.Equal(t, errors.New("Empty service name"), err, "Service name cannot be empty and should return error")
}

func TestNewWatchdog_Success_IncorrectServiceName(t *testing.T) {
	fsOsMock := new(utils.MockOsFs)
	fsOsMock.On("Stat", "/etc/init.d/xyz").Return(nil, errors.New("Could not find service: xyz"))
	_, err := NewWatchdog("xyz", "60s", "10s", 4, fsOsMock)
	assert.Equal(t, errors.New("Could not find service: xyz"), err, "Service should exist function should return error")
}

func TestIsRunning_Success(t *testing.T) {
	watchdog := getNewWatchdog()
	osFsMock := new(utils.MockOsFs)
	osFsMock.On("ExecAndOutput", "service", []string{"bluetooth", "status"}).Return([]byte{}, nil)
	watchdog.name = "bluetooth"
	watchdog.os = osFsMock
	output := watchdog.IsRunning()
	assert.True(t, output, "Service should run and function return true")
}

func TestIsRunning_Fail(t *testing.T) {
	watchdog := getNewWatchdog()
	osFsMock := new(utils.MockOsFs)
	osFsMock.On("ExecAndOutput", "service", []string{"bluetooth", "status"}).Return([]byte{}, errors.New("Some error"))
	watchdog.name = "bluetooth"
	watchdog.os = osFsMock

	output := watchdog.IsRunning()
	assert.False(t, output, "If service isnt running return false ")
}

func TestStart_Success(t *testing.T) {
	watchdog := getNewWatchdog()
	watchdog.name = "bluetooth"
	osFsMock := new(utils.MockOsFs)
	osFsMock.On("ExecAndOutput", "service", []string{"bluetooth", "start"}).Return([]byte{}, nil)
	watchdog.os = osFsMock

	output := watchdog.Start()

	assert.True(t, output, "Return true if servie run after start")
}

func TestStart_fail(t *testing.T) {
	watchdog := getNewWatchdog()
	watchdog.name = "bluetooth"
	osFsMock := new(utils.MockOsFs)
	osFsMock.On("ExecAndOutput", "service", []string{"bluetooth", "start"}).Return([]byte{}, errors.New("Some error"))
	watchdog.os = osFsMock

	output := watchdog.Start()

	assert.False(t, output, "Return false if service did not run after start")
}
