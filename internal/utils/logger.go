package utils

import log "github.com/sirupsen/logrus"

func SetupLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}
