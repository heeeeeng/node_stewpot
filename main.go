package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

var recorder Recorder

func main() {
	readConfig()
	initLogger()

	rand.Seed(time.Now().UnixNano())

	stewpot := NewStewpot()
	stewpot.InitNetwork(100)
	stewpot.PrintOutNodes()
	stewpot.Start()

	//stewpot.SendNewMsg()

	//time.Sleep(time.Second * 1000)

	// server
	http.HandleFunc("/", stewpot.MainPage)
	http.HandleFunc("/static/", stewpot.StaticController)
	http.HandleFunc("/graph", stewpot.GetNetworkGraph)
	http.HandleFunc("/send_msg", stewpot.SendMsg)
	http.HandleFunc("/time_unit", stewpot.GetTimeUnit)

	port := "10000"
	fmt.Println("---------Server Start!---------")
	fmt.Println("Port: ", port)
	logrus.Fatal(http.ListenAndServe(":"+port, nil))
}

func readConfig() {
	configPath := "config.toml"

	if !fileExists(configPath) {
		panicIfError(fmt.Errorf("cannot find the config: %s", configPath), "")
	}

	absPath, err := filepath.Abs(configPath)
	panicIfError(err, fmt.Sprintf("Error on parsing config file path: %s", absPath))

	file, err := os.Open(absPath)
	panicIfError(err, fmt.Sprintf("Error on opening config file: %s", absPath))
	defer file.Close()

	viper.SetConfigType("toml")
	err = viper.MergeConfig(file)
	panicIfError(err, fmt.Sprintf("Error on reading config file: %s", absPath))
	return
}

func initLogger() {
	logdir := viper.GetString("log.log_dir")

	var writer io.Writer

	if logdir != "" {
		folderPath, err := filepath.Abs(logdir)
		panicIfError(err, fmt.Sprintf("Error on parsing log path: %s", logdir))

		_, err = filepath.Abs(path.Join(logdir, "run"))
		panicIfError(err, fmt.Sprintf("Error on parsing log file path: %s", logdir))

		err = os.MkdirAll(folderPath, os.ModePerm)
		panicIfError(err, fmt.Sprintf("Error on creating log dir: %s", folderPath))

	} else {
		// stdout only
		fmt.Println("Will be logged to stdout")
		writer = os.Stdout
	}

	logrus.SetOutput(writer)

	// Only log the warning severity or above.
	switch viper.GetString("log.level") {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		fmt.Println("Unknown level", viper.GetString("log.level"), "Set to INFO")
		logrus.SetLevel(logrus.InfoLevel)
	}

	Formatter := new(logrus.TextFormatter)
	Formatter.ForceColors = logdir == ""
	//Formatter.DisableColors = true
	Formatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	Formatter.FullTimestamp = true

	logrus.SetFormatter(Formatter)

	lineNum := viper.GetBool("log_line_number")
	if lineNum {
		logrus.SetReportCaller(true)
	}

	//logger := logrus.StandardLogger()
	//logrus.Debug("Logger initialized.")

}

func panicIfError(err error, message string) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
