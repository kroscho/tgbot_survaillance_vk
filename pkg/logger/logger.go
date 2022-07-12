package logger

import (
	"io"
	"os"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

type Format string

const (
	StackDriver Format = "stackdriver"
)

type Option func(l *logrus.Logger)

func New(opts ...Option) *logrus.Logger {
	l := logrus.New()
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// nolint:gocritic
func WithFormat(format Format) Option {
	return func(l *logrus.Logger) {
		switch format {
		case StackDriver:
			formatter := joonix.NewFormatter()
			formatter.DisableTimestamp = false
			l.SetFormatter(formatter)
		}
	}
}

func WithLevel(level string) Option {
	return func(l *logrus.Logger) {
		lev, err := logrus.ParseLevel(level)
		if err != nil {
			lev = logrus.InfoLevel
		}
		l.SetLevel(lev)
	}
}

func WithFileOutput(file *os.File) Option {
	return func(l *logrus.Logger) {
		l.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
		mw := io.MultiWriter(os.Stdout, file)
		l.SetOutput(mw)
	}
}
