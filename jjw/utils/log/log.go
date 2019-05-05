package log

import (
	"os"
	"strings"
	"errors"
	"fmt"
	baseLog "log"

	"github.com/findthefirst/quant-rest/jjw/utils/config"
)

const logSize int64 = 1024 * 1024
var (
	logPath string
)

func init()  {
	conf, err := config.GetConfig()
	if err != nil {
		baseLog.Println(" init log error, can not find config")
	}
	logPath = conf.LogPath
}

func Debug(logStr string, v ...interface{}) {
	log(logStr, "DEBUG", v)
}

func Info(logStr string, v ...interface{}) {
	log(logStr, "INFO", v)
}

func WarningAndWrap(logStr string, v ...interface{}) (err error) {
	return log(logStr, "WARNING", v)
}

func ErrorAndWrap(logStr string, v ...interface{} ) (err error) {
	return log(logStr, "ERROR", v)
}

func log(logStr string, level string, v ...interface{}) (rtnErr error) {
	levelLogPath := strings.Join([]string{logPath, "-", level, ".log"}, "")
	file, err := os.OpenFile(levelLogPath, os.O_APPEND|os.O_CREATE, 0777)
	defer file.Close()
	/*if os.IsNotExist(err) {
		file, err = os.Create(logPath)
		if err != nil {
			baseLog.Println(" create logFile  error", err)
		}
	} else*/
	rtnErr = errors.New(fmt.Sprintf(logStr, v))
	if err != nil {
		baseLog.Println("open logFile Error() error", err)
		return rtnErr
	}
	if checkFileSize(file) {
		logger := baseLog.New(file, "[ "+level+" ]", 1)
		logger.Println(logStr)
	//	colorPrintln(level, "[ "+level+" ] " + logStr)
		baseLog.Println("[ "+level+" ] " + logStr)
		return rtnErr
	} else {
		baseLog.Println("logFile Size too big")
		return rtnErr
	}
}

func checkFileSize(file *os.File) (canLog bool) {
	info, err := file.Stat()
	if err != nil {
		//Error(" check fileSize error")
	}
	return info.Size() < logSize
	//fmt.Printf(" file size %d \n", info.Size())
}

func colorPrintln(level string, v ...interface{}) {
	switch level {
	case "INFO":
		fmt.Println(green, fmt.Sprint(v), reset)
	default:
		baseLog.Println(fmt.Sprint(v))
	}
}

var (
	greenBg      = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	whiteBg      = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellowBg     = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	redBg        = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blueBg       = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magentaBg    = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyanBg       = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	green        = string([]byte{27, 91, 51, 50, 109})
	white        = string([]byte{27, 91, 51, 55, 109})
	yellow       = string([]byte{27, 91, 51, 51, 109})
	red          = string([]byte{27, 91, 51, 49, 109})
	blue         = string([]byte{27, 91, 51, 52, 109})
	magenta      = string([]byte{27, 91, 51, 53, 109})
	cyan         = string([]byte{27, 91, 51, 54, 109})
	reset        = string([]byte{27, 91, 48, 109})
	disableColor = false
)