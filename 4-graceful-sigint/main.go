package main

import (
	"os"
	"os/signal"
)

func main() {
	proc := MockProcess{}
	signals := make(chan os.Signal, 1)
	procStopped := make(chan struct{})
	procDone := make(chan struct{})
	go func() {
		proc.Run()
		procDone <- struct{}{}
	}()
	signal.Notify(signals, os.Interrupt)
	select {
	case <-signals:
		go func() {
			proc.Stop()
			procStopped <- struct{}{}
		}()
		select {
		case <-signals:
			os.Exit(1)
		case <-procStopped:
			return
		}
	case <-procDone:
		return
	}
}
