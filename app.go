package main

import (
	"Platform/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	cmd.Log.Info("Starting Platform...")
	cmd.Log.WithFields(logrus.Fields{
		"action": "start",
		"status": "success",
		}).Info("Application started successfully")

	cmd.RunServer()

	cmd.Log.WithFields(logrus.Fields{
        "action": "stop",
        "status": "success",
        }).Info("Application stopped successfully")
	cmd.Log.Info("Platform stopped")
}
