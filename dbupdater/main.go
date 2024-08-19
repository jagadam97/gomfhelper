package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jagadam97/dbupdater/db"
)

func main() {
	go forever()
	quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
    <-quitChannel
	fmt.Println("Adios!")
}

func forever()  {
	didNotRunforThisweek := true
	fetchedTodaysNav := false
	for {
		currentTime := time.Now().Local()
		if didNotRunforThisweek{
			db.UpdateNavHistoryWatched()
			didNotRunforThisweek = false
		}
		if currentTime.Weekday().String() == "Sunday" && !didNotRunforThisweek {
			didNotRunforThisweek = true
		}
		if (currentTime.Hour() == 23 || currentTime.Hour() == 10) && currentTime.Minute() == 05{
			time.Sleep(time.Minute)
			fetchedTodaysNav = false
		}
		if !fetchedTodaysNav {
			db.UpdateLatestNav()
			fetchedTodaysNav = true
		}
	}

}