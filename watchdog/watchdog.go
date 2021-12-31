package watchdog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/fpiwowarczyk/watchdogGO/db"
	"github.com/fpiwowarczyk/watchdogGO/notifier"
)

const (
	checkForSettingsTime = time.Minute * 15
	serviceDown          = iota
	serviceStart         = iota
	serviceCannotStart   = iota
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
		return nil, errors.New("Incorrect number of tries to run service")
	}

	if len(name) < 1 {
		return nil, errors.New("Service name cannot be empty")
	}

	if _, err := os.Stat("/etc/init.d/" + name); err != nil {
		return nil, errors.New(fmt.Sprintf("Service %s doesn't exist\n", name))
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

func notify(notifier *notifier.Notifier, service string, attempts, status int) {
	var logMsg string
	// var title string

	switch status {
	case serviceDown:
		logMsg = fmt.Sprintf("%s Service %s is down", time.Now().String(), service)
	case serviceStart:
		logMsg = fmt.Sprintf("%s Service %s has been sterted after %d attempts", time.Now().String(), service, attempts)
	case serviceCannotStart:
		logMsg = fmt.Sprintf("%s Service %s can't be started after %d attempts", time.Now().String(), service, attempts)
	}

	notifier.Publish(&logMsg)
	log.Println(logMsg)
}

func (watchdog *Watchdog) updateSettings(newSettings *db.Settings) {
	checkInterval, err := time.ParseDuration(newSettings.NumOfSecCheck)
	if err != nil {
		log.Println(err)
	}

	retryInterval, err := time.ParseDuration(newSettings.NumOfSecWait)
	if err != nil {
		log.Println(err)
	}

	attemptVal, err := strconv.Atoi(newSettings.NumOfAttempts)
	if err != nil {
		log.Println(err)
	}

	watchdog.name = newSettings.ListOfServices
	watchdog.numOfSecCheck = checkInterval
	watchdog.numOfSecWait = retryInterval
	watchdog.numOfAttempts = attemptVal

}

func (watchdog *Watchdog) Watch(notifier *notifier.Notifier, stop chan bool) error {
	watching := true
	checkStatus := make(chan time.Time)
	startService := make(chan time.Time)
	checkSettings := make(chan time.Time)

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
			notify(notifier, watchdog.name, 0, serviceDown)
			for i := 0; i <= watchdog.numOfAttempts && watching; i++ {
				if run = watchdog.Start(); run {
					notify(notifier, watchdog.name, i, serviceStart)
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
			notify(notifier, watchdog.name, watchdog.numOfAttempts, serviceCannotStart)
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

		go func() {
			time.Sleep(checkForSettingsTime)
			db := db.New()
			sett, err := db.GetItem("1")
			if err != nil {
				log.Println(err)
			}
			watchdog.updateSettings(sett)
			checkSettings <- time.Now()
		}()
		<-checkSettings
	}
	return nil

}
