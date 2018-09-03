package util

import (
	"io"
	"log"
	"os"
)

// InitLogger initializes the go logger to write to file logFilepath
func InitLogger(logFilepath string) {
	if logFilepath == "" {
		// no path specified, skip
		return
	}

	// initialize logger
	logfile, err := os.OpenFile(logFilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		println("Error opening file for logging: %v", err)
		return
	}
	defer logfile.Close()
	mw := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(mw)
}
