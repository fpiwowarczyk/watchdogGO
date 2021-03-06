package watchdog

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fpiwowarczyk/watchdogGO/notifier"
	"github.com/fpiwowarczyk/watchdogGO/utils"
)

type Watchdog struct {
	name          string
	numOfSecCheck time.Duration
	numOfSecWait  time.Duration
	numOfAttempts int
	os            utils.FileSystem
}

func NewWatchdog(name, numOfSecCheck, numOfSecWait string, attempts int, os utils.FileSystem) (*Watchdog, error) {
	checkInterval, err := time.ParseDuration(numOfSecCheck)
	if err != nil {
		return nil, err
	}

	retryInterval, err := time.ParseDuration(numOfSecWait)
	if err != nil {
		return nil, err
	}

	if attempts < 1 {
		return nil, errors.New("Attempts number less than 1")
	}

	if len(name) < 1 {
		return nil, errors.New("Empty service name")
	}

	if _, err := os.Stat("/etc/init.d/" + name); err != nil {
		return nil, errors.New("Could not find service: " + name)
	}
	watchdog := new(Watchdog)
	watchdog.name = name
	watchdog.numOfSecCheck = checkInterval
	watchdog.numOfSecWait = retryInterval
	watchdog.numOfAttempts = attempts
	watchdog.os = os

	return watchdog, nil
}

func (watchdog *Watchdog) IsRunning() bool {
	_, err := watchdog.os.ExecAndOutput("service", watchdog.name, "status")
	if err != nil {
		return false
	}
	return true
}

func (watchdog *Watchdog) Start() bool {
	_, err := watchdog.os.ExecAndOutput("service", watchdog.name, "start")
	if err != nil {
		return false
	}
	return true
}

func (watchdog *Watchdog) Watch(notifier *notifier.Notifier, stop chan bool) error {
	watching := true
	checkStatus := make(chan time.Time)
	startService := make(chan time.Time)

	log.Printf("Watchdoggo for sevice %s start running", watchdog.name)
	go func() {
		<-stop
		watching = false
		checkStatus <- time.Now()
		startService <- time.Now()
	}()

	for watching {
		run := watchdog.IsRunning()
		if !run {
			notifier.Notify(fmt.Sprintf("Service %s is down", watchdog.name))
			for i := 1; i <= watchdog.numOfAttempts && watching; i++ {
				if run = watchdog.Start(); run {
					notifier.Notify(fmt.Sprintf("Service %s has been started after %d attempts", watchdog.name, i))
					break
				}
				go func() {
					time.Sleep(watchdog.numOfSecWait)
					startService <- time.Now()
				}()
				<-startService
			}
		}
		if !run {
			notifier.Notify(fmt.Sprintf("Service %s can't be started after %d attempts", watchdog.name, watchdog.numOfAttempts))
			return errors.New("Failed to start service")
		}
		if !watching {
			return nil
		}
		go func() {
			time.Sleep(watchdog.numOfSecCheck)
			checkStatus <- time.Now()
		}()
		<-checkStatus
	}
	return nil

}
