package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/fpiwowarczyk/watchdogGO/db"
	"github.com/fpiwowarczyk/watchdogGO/notifier"
	"github.com/fpiwowarczyk/watchdogGO/watchdog"
	"github.com/sevlyar/go-daemon"
)

var (
	serviceId = flag.String("id", "1", "Dynamodb id of row with settings")
)

func main() {
	var wg sync.WaitGroup
	var services []*watchdog.Watchdog
	flag.Parse()
	db := db.New()
	sett, err := db.GetItem(*serviceId)
	if err != nil {
		log.Println(err)
	}

	notifier := notifier.New()

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

	defer context.Release()

	stop := make(chan bool)
	sigc := make(chan os.Signal, 1)

	signal.Notify(sigc,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		_ = <-sigc
		for range services {
			stop <- true
		}
	}()

	for _, service := range services {
		wg.Add(1)
		go service.Watch(notifier, stop, &wg)
	}
	wg.Wait()

}
