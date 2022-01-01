package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/fpiwowarczyk/watchdogGO/db"
	"github.com/fpiwowarczyk/watchdogGO/notifier"
	"github.com/fpiwowarczyk/watchdogGO/watchdog"
	"github.com/sevlyar/go-daemon"
)

var (
	settingsID = flag.String("id", "1", "Dynamodb id of settings")
)

func updateSettings(notifier *notifier.Notifier, wg *sync.WaitGroup) {
	db := db.New()
	working := true
	stop := make(chan bool)
	sigc := make(chan os.Signal, 1)

	signal.Notify(sigc,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for working {
		var services []*watchdog.Watchdog
		sett, err := db.GetItem(*settingsID)
		if err != nil {
			log.Println(err)
		}

		attemptVal, err := strconv.Atoi(sett.NumOfAttempts)
		if err != nil {
			log.Println(err)
		}

		for _, serv := range sett.ListOfServices {
			service, err := watchdog.NewWatchdog(serv, sett.NumOfSecCheck, sett.NumOfSecWait, attemptVal)
			if err != nil {
				log.Println(err)
			}
			services = append(services, service)
		}

		go func() {
			_ = <-sigc
			for range services {
				stop <- true
			}
			working = false
			wg.Done()
		}()

		for _, service := range services {
			go service.Watch(notifier, stop)
		}

		// Update all settings every 15 mins
		time.Sleep(time.Minute * 15)
		for range services {
			stop <- true
		}
	}

}

func main() {
	var wg sync.WaitGroup
	flag.Parse()

	notifier := notifier.New()

	context := daemon.Context{
		LogFileName: "watchdog.log",
		LogFilePerm: 0644,
	}
	child, err := context.Reborn()
	if err != nil {
		log.Fatal(err)
	}
	if child != nil {
		return
	}
	wg.Add(1)
	go updateSettings(notifier, &wg)
	wg.Wait()

	defer context.Release()

}
