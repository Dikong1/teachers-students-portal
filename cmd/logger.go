package cmd

import (
	"os"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	file, err := os.OpenFile("actions.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.WithError(err).Fatal("Failed to open log file")
	}

	Log.SetOutput(file)
	Log.SetFormatter(&logrus.JSONFormatter{})
}