package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Entry

func InitLogger() {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	l.SetOutput(os.Stdout)

	Log = logrus.NewEntry(l)
}
