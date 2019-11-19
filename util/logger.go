package util

import (
	"io"
	"github.com/sirupsen/logrus"
)

//setUpLogs set the log output ans the log level
func SetUpLogs(out io.Writer, level string) error {

	// Set the the stdout as default output
	logrus.SetOutput(out)

	// Setting the level
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	// Set the global formatter
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// The level is provided from Cobra as global param
	logrus.SetLevel(lvl)
	return nil
}

/**
 * Verify if the system is configured with a given level
 */
func IsLogInDebug() bool {
	return logrus.GetLevel() == logrus.DebugLevel
}