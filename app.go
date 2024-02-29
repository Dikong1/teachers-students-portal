package main

import (
	"Platform/cmd"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	log.Info("Starting Platform...")
	log.WithFields(logrus.Fields{
		"action": "start",
		"status": "success",
	}).Info("Application started successfully")

	cmd.RunServer()

	log.WithFields(logrus.Fields{
		"action": "stop",
		"status": "success",
	}).Info("Application stopped successfully")
	log.Info("Platform stopped")
}
