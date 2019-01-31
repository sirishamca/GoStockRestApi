package srslogger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	SRS_LOGGER_LEVEL_FATAL   uint16 = 0
	SRS_LOGGER_LEVEL_ERROR   uint16 = 1
	SRS_LOGGER_LEVEL_WARNING uint16 = 2
	SRS_LOGGER_LEVEL_INFO    uint16 = 3
	SRS_LOGGER_LEVEL_DEBUG   uint16 = 4
	SRS_LOG_FILE_SIZE        int64  = 10 * 1024 * 1024
)

type LogLine struct {
	LogLine string
	Level   uint16
}

type Logger struct {
	LogCh       chan LogLine
	LogDebug    *log.Logger
	LogInfo     *log.Logger
	LogWarning  *log.Logger
	LogError    *log.Logger
	LogFatal    *log.Logger
	logLevel    uint16
	file        *os.File
	fileName    string
	currVersion int
}

func (l *Logger) Init(fileName string) {

	var err error

	l.LogCh = make(chan LogLine)
	l.fileName = fileName
	l.logLevel = SRS_LOGGER_LEVEL_DEBUG

	l.file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", ":", err)
	}

	l.LogDebug = log.New(l.file, "DEBUG: ", log.Ldate|log.Ltime)

	l.LogInfo = log.New(l.file, "INFO: ", log.Ldate|log.Ltime)

	l.LogWarning = log.New(l.file, "WARNING: ", log.Ldate|log.Ltime)

	l.LogError = log.New(l.file, "ERROR: ", log.Ldate|log.Ltime)

	l.LogFatal = log.New(l.file, "FATAL: ", log.Ldate|log.Ltime)

	go func() {
		for {
			line := <-l.LogCh
			switch line.Level {
			case SRS_LOGGER_LEVEL_DEBUG:
				l.LogDebug.Println(line.LogLine)
			case SRS_LOGGER_LEVEL_INFO:
				l.LogInfo.Println(line.LogLine)
			case SRS_LOGGER_LEVEL_WARNING:
				l.LogWarning.Println(line.LogLine)
			case SRS_LOGGER_LEVEL_ERROR:
				l.LogError.Println(line.LogLine)
			case SRS_LOGGER_LEVEL_FATAL:
				l.LogFatal.Println(line.LogLine)
			default:
				l.LogDebug.Println(line.LogLine)
			}
			l.CheckSize(fileName)
		}
	}()
}

func (l *Logger) SetLoglevel(level string) {
	var levelInt uint16
	switch level {
	case "DEBUG":
		levelInt = SRS_LOGGER_LEVEL_DEBUG
	case "INFO":
		levelInt = SRS_LOGGER_LEVEL_INFO
	case "ERROR":
		levelInt = SRS_LOGGER_LEVEL_ERROR
	case "WARNING":
		levelInt = SRS_LOGGER_LEVEL_WARNING
	case "FATAL":
		levelInt = SRS_LOGGER_LEVEL_FATAL
	default:
		levelInt = SRS_LOGGER_LEVEL_FATAL
	}
	l.logLevel = levelInt
}

func (l *Logger) GetCaller() (string, int) {

	_, f, n, _ := runtime.Caller(2)

	x := strings.Split(f, "/")
	return x[len(x)-1], n
}

func (l *Logger) Fatal(args ...interface{}) {

	fileName, lineno := l.GetCaller()
	line := fmt.Sprintln(fileName, lineno, args)
	l.LogCh <- LogLine{line, SRS_LOGGER_LEVEL_FATAL}
}

func (l *Logger) Fatalf(format string, args ...interface{}) {

	fileName, lineno := l.GetCaller()
	line := fmt.Sprintf("%s:%d"+format, fileName, lineno, args)
	l.LogCh <- LogLine{line, SRS_LOGGER_LEVEL_FATAL}

}

func (l *Logger) Error(errStr ...interface{}) {

	if l.logLevel < SRS_LOGGER_LEVEL_ERROR {
		return
	}

	fileName, lineno := l.GetCaller()
	line := fmt.Sprintln(fileName, lineno, errStr)
	l.LogCh <- LogLine{line, SRS_LOGGER_LEVEL_ERROR}

}

func (l *Logger) Warning(warningStr ...interface{}) {
	if l.logLevel < SRS_LOGGER_LEVEL_WARNING {
		return
	}

	fileName, lineno := l.GetCaller()
	line := fmt.Sprintln(fileName, lineno, warningStr)
	l.LogCh <- LogLine{line, SRS_LOGGER_LEVEL_WARNING}

}

func (l *Logger) Info(infoStr ...interface{}) {
	if l.logLevel < SRS_LOGGER_LEVEL_INFO {
		return
	}

	fileName, lineno := l.GetCaller()
	line := fmt.Sprintln(fileName, lineno, infoStr)
	l.LogCh <- LogLine{line, SRS_LOGGER_LEVEL_INFO}

}

func (l *Logger) Debug(debugStr ...interface{}) {
	if l.logLevel < SRS_LOGGER_LEVEL_DEBUG {
		return
	}

	fileName, lineno := l.GetCaller()
	line := fmt.Sprintln(fileName, lineno, debugStr)
	l.LogCh <- LogLine{line, SRS_LOGGER_LEVEL_DEBUG}

}

func (l *Logger) CheckSize(fileName string) {
	fileInfo, err := l.file.Stat()
	if err != nil {
		fmt.Println("Error checking file size", err)
		return
	}
	if fileInfo.Size() > SRS_LOG_FILE_SIZE {
		fmt.Println("Log file running out of space, flushing the logs.")
		os.Remove(l.fileName)
		l.file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", ":", err)
		}

		l.LogDebug = log.New(l.file, "DEBUG: ", log.Ldate|log.Ltime)

		l.LogInfo = log.New(l.file, "INFO: ", log.Ldate|log.Ltime)

		l.LogWarning = log.New(l.file, "WARNING: ", log.Ldate|log.Ltime)

		l.LogError = log.New(l.file, "ERROR: ", log.Ldate|log.Ltime)

		l.LogFatal = log.New(l.file, "FATAL: ", log.Ldate|log.Ltime)
	}
}
