package watchdog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	serviceDown        = iota
	serviceStart       = iota
	serviceCannotStart = iota
)

type Watchdog struct {
	names         []string
	numOfSecCheck time.Duration
	numOfSecWait  time.Duration
	numOfAttempts int
}

func NewWatchdog(name []string, numOfSecCheck, numOfSecWait string, attempts int) (*Watchdog, error) {
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
	watchdog.names = name
	watchdog.numOfSecCheck = checkInterval
	watchdog.numOfSecWait = retryInterval
	watchdog.numOfAttempts = attempts

	return watchdog, nil
}

func IsRunning(name string) bool {
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

func notify(service string, attempts, status int) {
	var logMsg string

	switch status {
	case serviceDown:
		logMsg = fmt.Sprintf("%s Service %s is down", time.Now().String(), service)
	case serviceStart:
		logMsg = fmt.Sprintf("%s Service %s can't be started after %d attempts", time.Now().String(), service, attempts)
	case serviceCannotStart:
		logMsg = fmt.Sprintf("%s Service %s can't be started after %d attempts", time.Now().String(), service, attempts)
	}

	log.Println(logMsg)
}

func (service *Watchdog) Watch(stop chan bool) error {
	watching := true
	checkStatus := make(chan time.Time)
	startService := make(chan time.Time)

	log.Printf("Watch dog for sevice %s start running", service.name)
	go func() {
		<-stop
		watching = false
		checkStatus <- time.Now()
		startService <- time.Now()
	}()

	for watching {
		for _, name := range service.name {
			run := IsRunning(name)
			if !run {
				notify(service.name, 0, serviceDown)
				for i := 0; i <= service.numOfAttempts && watching; i++ {
					if run = service.Start(); run {
						notify(service.name, i, serviceStart)
						break
					}

					go func() {
						time.Sleep(service.numOfSecWait)
						startService <- time.Now()
					}()
					<-startService
				}
			}
			if !run {
				notify(service.name, service.numOfAttempts, serviceCannotStart)
				return errors.New("Failed to start service")
			}
			if !watching {
				return nil
			}
			go func() {
				time.Sleep(service.numOfSecCheck)
				checkStatus <- time.Now()
			}()
			<-checkStatus
		}
		return nil
	}

}
