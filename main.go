package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fpiwowarczyk/watchdogGO/db"
	"github.com/fpiwowarczyk/watchdogGO/watchdog"
	"github.com/sevlyar/go-daemon"
)

var (
	serviceId = flag.String("id", "1", "Dynamodb id of row with settings")
)

func main() {
	flag.Parse()
	db := db.New()
	sett, err := db.GetItem(*serviceId)

	if err != nil {
		log.Println(err)
	}

	attemptVal, err := strconv.Atoi(sett.NumOfAttempts)
	if err != nil {
		log.Println(err)
	}

	service, err := watchdog.NewWatchdog(sett.ListOfServices, sett.NumOfSecCheck, sett.NumOfSecWait, attemptVal)
	if err != nil {
		log.Println(err)
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
		stop <- true
	}()

	service.Watch(stop)
}
