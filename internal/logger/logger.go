package logger

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

type Logger struct {
	LogFile *os.File
}

func NewLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return &Logger{LogFile: f}, nil
}

func (l *Logger) Write(s string) {
	_, _ = l.LogFile.WriteString(time.Now().Format(time.DateTime) + " " + s + "\n")
}

func Infof(format string, v ...interface{}) {
	fmt.Printf(Blue+"[-] "+format+"\n"+Reset, v...)
}

func Sinfof(format string, v ...interface{}) string {
	return fmt.Sprintf(Blue+"[-] "+format+"\n"+Reset, v...)
}

func Info(v interface{}) {
	fmt.Println(Blue+"[-]", v, Reset)
}

func Fatal(v interface{}) {
	fmt.Println(Purple+"[x]", v, Reset)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	fmt.Printf(Purple+"[x] "+format+"\n"+Reset, v...)
	os.Exit(1)
}

func Notify(v interface{}) {
	fmt.Println(Green+"[+]", v, Reset)
}

func Notifyf(format string, v ...interface{}) {
	fmt.Printf(Green+"[+] "+format+"\n"+Reset, v...)
}

func Warn(v interface{}) {
	fmt.Println(Yellow+"[?]", v, Reset)
}

func Warnf(format string, v ...interface{}) {
	fmt.Printf(Yellow+"[?] "+format+"\n"+Reset, v...)
}

func Error(v interface{}) {
	fmt.Println(Red+"[!]", v, Reset)
}

func Errorf(format string, v ...interface{}) {
	fmt.Printf(Red+"[!] "+format+"\n"+Reset, v...)
}

func Serrorf(format string, v ...interface{}) string {
	return fmt.Sprintf(Red+"[!] "+format+"\n"+Reset, v...)
}

func Silentf(format string, v ...interface{}) {
	logFile := filepath.Join(
		filepath.Dir(viper.ConfigFileUsed()),
		"siphon.log",
	)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.WriteString(time.Now().Format(time.DateTime) + " " + fmt.Sprintf(format, v...) + "\n")
	if err != nil {
		return
	}
}
