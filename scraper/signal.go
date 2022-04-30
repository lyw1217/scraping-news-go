package scraper

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"scraping-news/util"

	log "github.com/sirupsen/logrus"
)

func WaitSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.Signal(0x01), // SIGHUP
		syscall.Signal(0x02), // SIGINT
		syscall.Signal(0x0f)) // SIGTERM

	log.Error("< SCRAPER > Signal Initialized ....................")

	exitChan := make(chan int)

	go func() {
		for !sysClose {
			s := <-signalChan

			switch s.(syscall.Signal) {
			case 0x01:
				log.Info("< SCRAPER > Receive < SIGHUP > Signal")

			case 0x02:
				log.Info("< SCRAPER > Receive < SIGINT > Signal")
				sysClose = true
				util.SysClose = true
				exitChan <- 2

			case 0x0f:
				log.Info("< SCRAPER > Receive < SIGTREM > Signal")
				sysClose = true
				util.SysClose = true
				exitChan <- 15

			default:
				log.Infof("< SCRAPER > Receive Unknown signal, Sig = '%s(%x)'", s.String(), s.(syscall.Signal))
			}
		}
	}()

	exitCode := <-exitChan
	TerminateHandler(exitCode)
}

func TerminateHandler(c int) {
	log.Infof("< SCRAPER > TERMINATE Handler is called. (reason : %d) ", c)
	time.Sleep(time.Duration(3) * time.Second)

	os.Exit(c)
}
