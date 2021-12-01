package scraper

import (
	"os"
	"os/signal"
	"syscall"
	"time"

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
		for !gSysClose {
			s := <-signalChan

			switch s.(syscall.Signal) {
			case 0x01:
				log.Error("< SCRAPER > Receive < SIGHUP > Signal")

			case 0x02:
				log.Error("< SCRAPER > Receive < SIGINT > Signal")
				gSysClose = true
				exitChan <- 2

			case 0x0f:
				log.Error("< SCRAPER > Receive < SIGTREM > Signal")
				gSysClose = true
				exitChan <- 15

			default:
				log.Errorf("< SCRAPER > Receive Unknown signal, Sig = '%s(%x)'", s.String(), s.(syscall.Signal))
			}
		}
	}()

	exitCode := <-exitChan
	TerminateHandler(exitCode)
}

func TerminateHandler(c int) {
	log.Errorf("< SCRAPER > TERMINATE Handler is called. (reason : %d) ", c)
	time.Sleep(time.Duration(3) * time.Second)

	os.Exit(c)
}
