package watchdog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fpiwowarczyk/watchdogGO/notifier"
)

type Watchdog struct {
	name          string
	numOfSecCheck time.Duration
	numOfSecWait  time.Duration
	numOfAttempts int
}

func NewWatchdog(name, numOfSecCheck, numOfSecWait string, attempts int) (*Watchdog, error) {
	checkInterval, err := time.ParseDuration(numOfSecCheck)
	if err != nil {
		return nil, err
	}

	retryInterval, err := time.ParseDuration(numOfSecWait)
	if err != nil {
		return nil, err
	}

	if attempts < 1 {
		return nil, errors.New("Incorrect number of attempts to run service")
	}

	if len(name) < 1 {
		return nil, errors.New("Service name cannot be empty")
	}

	if _, err := os.Stat("/etc/init.d/" + name); err != nil {
		return nil, errors.New("Could not find service: " + name)
	}
	watchdog := new(Watchdog)
	watchdog.name = name
	watchdog.numOfSecCheck = checkInterval
	watchdog.numOfSecWait = retryInterval
	watchdog.numOfAttempts = attempts

	return watchdog, nil
}

func (service *Watchdog) IsRunning() bool {
	_, err := exec.Command("service", service.name, "status").Output()
	if err != nil {
		return false
	}
	return true
}

func (service *Watchdog) Start() bool {
	_, err := exec.Command("service", service.name, "start").Output()
	if err != nil {
		return false
	}
	return true
}

func (watchdog *Watchdog) Watch(notifier *notifier.Notifier, stop chan bool) error {
	watching := true
	checkStatus := make(chan time.Time)
	startService := make(chan time.Time)

	log.Printf("Watch dog for sevice %s start running", watchdog.name)
	go func() {
		<-stop
		watching = false
		checkStatus <- time.Now()
		startService <- time.Now()
	}()

	for watching {
		run := watchdog.IsRunning()
		if !run {
			notifier.Notify(fmt.Sprintf("%s Service %s is down", time.Now().String(), watchdog.name))
			for i := 1; i <= watchdog.numOfAttempts && watching; i++ {
				if run = watchdog.Start(); run {
					notifier.Notify(fmt.Sprintf("%s Service %s has been started after %d attempts", time.Now().String(), watchdog.name, i))
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
			notifier.Notify(fmt.Sprintf("%s Service %s can't be started after %d attempts", time.Now().String(), watchdog.name, watchdog.numOfAttempts))
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
