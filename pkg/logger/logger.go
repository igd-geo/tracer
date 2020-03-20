package logger

import (
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/broker"
)

type Logger struct {
	rbmq *broker.Session
}

func NewLogger(rbSession *broker.Session) *Logger {
	return &Logger{
		rbmq: rbSession,
	}
}

func (l *Logger) Info(logMsg string) {
	err := l.rbmq.Publish(logMsg, "log.info.tracer")
	if err != nil {
		log.Println(err)
	}
}

func (l *Logger) Error(logMsg string) {
	err := l.rbmq.Publish(logMsg, "log.error.tracer")
	if err != nil {
		log.Println(err)
	}
}

func (l *Logger) Warning(logMsg string) {
	err := l.rbmq.Publish(logMsg, "log.warning.tracer")
	if err != nil {
		log.Println(err)
	}
}
