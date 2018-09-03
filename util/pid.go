package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// WritePid writes the current process id into pidfile
func WritePid(pidfile string) {
	if pidfile == "" {
		// no path specified, skip
		return
	}

	// write a pid file too avoid searching the processes:
	// 'kill -SIGUSR1 $(cat analysis.pid)'
	log.Printf("Writing PID %d to file '%s'.", os.Getpid(), pidfile)
	err := ioutil.WriteFile(pidfile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		log.Fatal("Could not write PID File.")
	}
}
