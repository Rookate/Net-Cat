package netcat

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

func GenerateLogFileName() string {
	return fmt.Sprintf("ChatLog_%s.log", time.Now().Format(time.Stamp))
}

func OpenLogFile(filename string) error {
	var err error
	logFile, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("error opeing log file: %v", err)
	}
	return nil
}

func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}

func LogMessage(message string) {
	if logFile != nil {
		_, err := logFile.WriteString(message + "\n")
		if err != nil {
			fmt.Println("Error writing to log file:", err)
		}
	}
}
