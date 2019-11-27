package logger

import (
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rbmq"
	"log"
)

type Logger struct {
	rbmq *rbmq.Session
}

func NewLogger(rbSession *rbmq.Session) *Logger {
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
